package service

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/utils"
	"math"
	"strconv"
	"sync"
	"time"
)

var (
	TheContinue3DaysReportMaster = NewContinue3DaysReportMaster()
	TheBigSwingsReportMaster     = NewBigSwingsReportMaster()
)

type Continue3DaysReportMaster struct {
	lock   *sync.RWMutex
	report string
}

func NewContinue3DaysReportMaster() *Continue3DaysReportMaster {
	return &Continue3DaysReportMaster{
		lock: &sync.RWMutex{},
	}
}

func (m *Continue3DaysReportMaster) SetReport(report string) {
	m.lock.Lock()
	m.report = report
	m.lock.Unlock()
}

func (m *Continue3DaysReportMaster) GetReport() string {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.report
}

type BigSwingsReportMaster struct {
	lock   *sync.RWMutex
	report string
}

func NewBigSwingsReportMaster() *BigSwingsReportMaster {
	return &BigSwingsReportMaster{
		lock: &sync.RWMutex{},
	}
}

func (m *BigSwingsReportMaster) SetReport(report string) {
	m.lock.Lock()
	m.report = report
	m.lock.Unlock()
}

func (m *BigSwingsReportMaster) GetReport() string {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.report
}

type SpecialTradingsReport struct {
	Date             string
	Percent          float64
	tradings         *ARKTradings
	previousTradings []*ARKTradings
}

type StockTradingTxtReport struct {
	Records map[string]map[string]*StockTrading
}

func NewStockTradingTxtReport() *StockTradingTxtReport {
	return &StockTradingTxtReport{
		Records: make(map[string]map[string]*StockTrading),
	}
}

func (h *StockTradingTxtReport) Add(stock string, trading *StockTrading) {
	if h.Records[stock] == nil {
		h.Records[stock] = make(map[string]*StockTrading)
	}
	h.Records[stock][trading.Fund] = trading
}

/*
ZY：ARKG建仓，买入xxx股；ARKK建仓，买入xxx股；ARKK增持5.6%，持股数从xxx股增到xxx股；
ARKK减持5.6%，持股数从xxx股减到xxx股；ARKG清仓，卖出xxx股；ARK总持股数从4014903股增到5241248股。
*/
func (h *StockTradingTxtReport) Report() []byte {
	var (
		report string
	)
	for ticker, record := range h.Records {
		var (
			stockReport                        = ticker + "："
			totalHolding, previousTotalHolding float64
		)

		for _, trading := range record {
			stockReport += NewStockTradingTxtFromTrading(trading)
			totalHolding += trading.Holding
			previousTotalHolding += trading.PreviousHolding
		}
		if len(record) > 1 {
			stockReport += NewHoldingChangeTxt(totalHolding, previousTotalHolding)
		}
		report += stockReport + "\n"
	}

	return []byte(report)
}

func NewSpecialTradingsReport(date time.Time, percent float64) *SpecialTradingsReport {
	var (
		r = &SpecialTradingsReport{
			Date:             date.Format(TheDateFormat),
			Percent:          percent,
			tradings:         TheLibrary.GetTradings(date),
			previousTradings: TheLibrary.GetPreviousTradings(date, 2),
		}
	)

	utils.CheckFolder(r.ReportFolder())

	return r
}

func (r *SpecialTradingsReport) Report() error {
	var (
		err      error
		fileName = r.ExcelPath()
	)

	if r.tradings == nil {
		glog.Warningf("Empty report")
		return nil
	}

	err = r.InitExcelFromTemplate()
	if err != nil {
		glog.Errorf("failed to init excel from template, err: %v", err)
		return err
	}

	f, err := excelize.OpenFile(fileName)
	if err != nil {
		glog.Errorf("failed to open excel %s, err: %v", fileName, err)
		return err
	}

	var (
		sheet                = defaultSheet
		idx                  = 2
		specialTradingReport = NewStockTradingTxtReport()
		//higherThan10Report        = NewStockTxtReport()
		continuousDirectionReport = NewStockTxtReport()
	)

	getTradingPercent := func(trading *ARKTradings, fund, ticker string) float64 {
		if trading == nil {
			return 0
		}
		theTradings := trading.GetFundStockTradings(fund)
		if theTradings == nil {
			return 0
		}
		stockTrading := theTradings.Tradings[ticker]
		if stockTrading == nil {
			return 0
		}
		return stockTrading.Percent
	}

	getTradingFixedDirection := func(trading *ARKTradings, fund, ticker string) TradeDirection {
		if trading == nil {
			return TradeDoNothing
		}
		theTradings := trading.GetFundStockTradings(fund)
		if theTradings == nil {
			return TradeDoNothing
		}
		stockTrading := theTradings.Tradings[ticker]
		if stockTrading == nil {
			return TradeDoNothing
		}
		return stockTrading.FixedDirection
	}

	getTradingHolding := func(trading *ARKTradings, fund, ticker string) float64 {
		if trading == nil {
			return 0
		}
		theTradings := trading.GetFundStockTradings(fund)
		if theTradings == nil {
			return 0
		}
		stockTrading := theTradings.Tradings[ticker]
		if stockTrading == nil {
			return 0
		}
		return stockTrading.Holding
	}

	for _, fund := range allARKTypes {
		var (
			toReportTradings []*StockTrading
		)
		tradings := r.tradings.GetFundStockTradings(fund)
		if tradings == nil {
			return errEmptyReport
		}

		for _, trading := range tradings.Tradings {
			if toSkipTicker(trading.Ticker) {
				continue
			}
			if !r.IsSpecialTradings(trading) {
				continue
			}
			toReportTradings = append(toReportTradings, trading)
		}

		for _, trading := range toReportTradings {
			// special trading
			if math.Abs(trading.Percent) > 10 {
				//higherThan10Report.Add(trading.Ticker, NewSpecialTradingTxtFromTrading(trading))
				specialTradingReport.Add(trading.Ticker, trading)
			}

			// trading trend
			previousFixedDirection1 := getTradingFixedDirection(r.previousTradings[1], fund, trading.Ticker)
			previousFixedDirection2 := getTradingFixedDirection(r.previousTradings[0], fund, trading.Ticker)
			if previousFixedDirection1 == trading.FixedDirection && previousFixedDirection2 == trading.FixedDirection {
				continuousDirectionReport.Add(trading.Ticker, NewContinuousDirectionSpecialTradingTxtFromTrading(trading,
					getTradingPercent(r.previousTradings[1], fund, trading.Ticker),
					getTradingPercent(r.previousTradings[0], fund, trading.Ticker),
					getTradingHolding(r.previousTradings[1], fund, trading.Ticker),
					getTradingHolding(r.previousTradings[0], fund, trading.Ticker),
				))
			}

			line := strconv.Itoa(idx)
			f.SetCellValue(sheet, "A"+line, trading.Fund)
			f.SetCellValue(sheet, "B"+line, trading.Ticker)
			f.SetCellValue(sheet, "C"+line, floatToStringIntOnlyWithSign(trading.Shards))
			f.SetCellValue(sheet, "D"+line, floatToPercentStringWithSign(trading.Percent))
			f.SetCellValue(sheet, "E"+line, floatToPercentStringWithSign(getTradingPercent(r.previousTradings[1], fund, trading.Ticker)))
			f.SetCellValue(sheet, "F"+line, floatToPercentStringWithSign(getTradingPercent(r.previousTradings[0], fund, trading.Ticker)))

			idx++

		}
	}

	err = f.Save()
	if err != nil {
		glog.Errorf("failed to save excel %s, err: %v", fileName, err)
		return err
	}

	higherThan10TxtContent := specialTradingReport.Report()
	if len(higherThan10TxtContent) > 0 {
		higherThan10TxtContent = append([]byte(fmt.Sprintf("基于ARK基金公开的截止%s（不含）的持仓数据，以下个股在上个交易日变动超过%.0f%%：\n", r.Date, r.Percent)), higherThan10TxtContent...)
		TheBigSwingsReportMaster.SetReport(string(higherThan10TxtContent))
		err = utils.NewFileStoreSvc(r.HigherThan10TxtPath()).Save(higherThan10TxtContent)
		if err != nil {
			glog.Errorf("failed to save txt %s, err: %v", r.HigherThan10TxtPath(), err)
			return err
		}
		utils.SendAlertV2("异动股票"+r.Date, string(higherThan10TxtContent))
	} else {
		glog.V(4).Infof("No higher than 10 tradings")
	}

	continuousDirectionTxtContent := continuousDirectionReport.Report()
	if len(continuousDirectionTxtContent) > 0 {
		continuousDirectionTxtContent = append([]byte(fmt.Sprintf("基于ARK基金公开的截止%s（不含）的持仓数据，以下个股在连续三个交易日发生同向变动：\n", r.Date)), continuousDirectionTxtContent...)
		TheContinue3DaysReportMaster.SetReport(string(continuousDirectionTxtContent))
		err = utils.NewFileStoreSvc(r.ContinuousDirectionTxtPath()).Save(continuousDirectionTxtContent)
		if err != nil {
			glog.Errorf("failed to save txt %s, err: %v", r.ContinuousDirectionTxtPath(), err)
			return err
		}
		utils.SendAlertV2("连续变动股票"+r.Date, string(continuousDirectionTxtContent))
	} else {
		glog.V(4).Infof("No Continuous Direction tradings")
	}

	glog.V(4).Infof("SpecialTradingsReport %s is provided", fileName)

	return nil
}

func (r *SpecialTradingsReport) ReportFolder() string {
	return reportPath + "/" + r.Date
}

func (r *SpecialTradingsReport) ExcelPath() string {
	return r.ReportFolder() + "/" + r.ExcelName()
}

func (r *SpecialTradingsReport) ExcelName() string {
	return prefixSpecialTradings + r.Date + ".xlsx"
}

func (r *SpecialTradingsReport) HigherThan10TxtPath() string {
	return r.ReportFolder() + "/" + r.HigherThan10TxtName()
}

func (r *SpecialTradingsReport) HigherThan10TxtName() string {
	return prefixSpecialTradingsHigherThan10 + r.Date + ".txt"
}

func (r *SpecialTradingsReport) ContinuousDirectionTxtPath() string {
	return r.ReportFolder() + "/" + r.ContinuousDirectionTxtName()
}

func (r *SpecialTradingsReport) ContinuousDirectionTxtName() string {
	return prefixSpecialTradingsContinuousDirection + r.Date + ".txt"
}

func (r *SpecialTradingsReport) InitExcelFromTemplate() error {
	var fileName = r.ExcelPath()
	if utils.CheckFileExist(fileName) {
		utils.DeleteFile(fileName)
	}
	utils.CopyFile(specialTradingsExcelTemplate, fileName)
	glog.V(4).Infof("Init fileName: %s", fileName)
	return nil
}

func (r *SpecialTradingsReport) InitTxt() error {
	var fileName = r.HigherThan10TxtPath()
	if utils.CheckFileExist(fileName) {
		utils.DeleteFile(fileName)
	}
	glog.V(4).Infof("Init fileName: %s", fileName)
	return nil
}

func (r *SpecialTradingsReport) IsSpecialTradings(trading *StockTrading) bool {
	return math.Abs(trading.Percent) >= r.Percent && trading.FixedDirection != TradeKeep
}

func NewSpecialTradingTxtFromTrading(trading *StockTrading) []byte {
	var (
		result string
	)
	switch trading.Direction {
	case TradeBuy:
		if trading.Percent == 100 {
			result = fmt.Sprintf("%s在%s中建仓，买入%.0f股。\n", trading.Ticker, trading.Fund, trading.Holding)
		} else {
			result = fmt.Sprintf("%s在%s中获得%.2f%%的增持，持股数从%.0f股增到%.0f股。\n", trading.Ticker, trading.Fund,
				math.Abs(trading.Percent), trading.PreviousHolding, trading.Holding)
		}
	case TradeSell:
		if trading.Percent == -100 {
			result = fmt.Sprintf("%s在%s中清仓，卖出%.0f股。\n", trading.Ticker, trading.Fund, trading.Shards)
		} else {
			result = fmt.Sprintf("%s在%s中被减持了%.2f%%，持股数从%.0f股减少到%.0f股。\n", trading.Ticker, trading.Fund,
				math.Abs(trading.Percent), trading.PreviousHolding, trading.Holding)
		}
	}
	return []byte(result)
}

func NewStockTradingTxtFromTrading(trading *StockTrading) string {
	var (
		result string
	)
	switch trading.Direction {
	case TradeBuy:
		if trading.Percent == 100 {
			result = fmt.Sprintf("%s建仓，买入%.0f股。", trading.Fund, trading.Holding)
			go utils.SendAlertV2(fmt.Sprintf("%s建仓%s，买入%.0f股。", trading.Fund, trading.Ticker, trading.Holding), "")
		} else {
			result = fmt.Sprintf("%s增持%.2f%%，买入%.0f股，持股数从%.0f股增到%.0f股。", trading.Fund,
				trading.Percent, trading.Shards, trading.PreviousHolding, trading.Holding)
			if trading.Percent > 50 {
				go utils.SendAlertV2(fmt.Sprintf("%s增持%s%.2f%%，买入%.0f股，持股数从%.0f股增到%.0f股。", trading.Fund,
					trading.Ticker, trading.Percent, trading.Shards, trading.PreviousHolding, trading.Holding), "")
			}
		}
	case TradeSell:
		if trading.Percent == -100 {
			result = fmt.Sprintf("%s清仓，卖出%.0f股。", trading.Fund, trading.PreviousHolding)
			go utils.SendAlertV2(fmt.Sprintf("%s清仓%s，卖出%.0f股。", trading.Fund, trading.Ticker, trading.Holding), "")
		} else {
			result = fmt.Sprintf("%s减持%.2f%%，卖出%.0f股，持股数从%.0f股减少到%.0f股。", trading.Fund,
				math.Abs(trading.Percent), math.Abs(trading.Shards), trading.PreviousHolding, trading.Holding)
			if trading.Percent > 50 {
				go utils.SendAlertV2(fmt.Sprintf("%s减持%s%.2f%%，卖出%.0f股，持股数从%.0f股减少到%.0f股。", trading.Fund, trading.Ticker,
					math.Abs(trading.Percent), math.Abs(trading.Shards), trading.PreviousHolding, trading.Holding), "")
			} else if trading.Percent > 20 && TheChinaStockManager.IsChinaStock(trading.Ticker) {
				go utils.SendAlertV2(fmt.Sprintf("%s减持%s%.2f%%，卖出%.0f股，持股数从%.0f股减少到%.0f股。", trading.Fund, trading.Ticker,
					math.Abs(trading.Percent), math.Abs(trading.Shards), trading.PreviousHolding, trading.Holding), "")
			}
		}
	}
	return result
}

// 今日， 昨日， 前日
func NewContinuousDirectionSpecialTradingTxtFromTrading(trading *StockTrading, previousPercent1, previousPercent2,
	previousHolding1, previousHolding2 float64) []byte {
	var (
		result string
	)
	switch trading.Direction {
	case TradeSell:
		result = fmt.Sprintf("%s最近三日在%s中均被减持，分别是%.2f%%、%.2f%%以及%.2f%%，持股数分别为%.0f股、%.0f股和%.0f股。\n", trading.Ticker, trading.Fund,
			math.Abs(trading.Percent), math.Abs(previousPercent1), math.Abs(previousPercent2),
			trading.Holding, previousHolding1, previousHolding2)
	case TradeBuy:
		result = fmt.Sprintf("%s最近三日在%s中都获得增持，分别是%.2f%%、%.2f%%以及%.2f%%，持股数分别为%.0f股、%.0f股和%.0f股。\n", trading.Ticker, trading.Fund,
			math.Abs(trading.Percent), math.Abs(previousPercent1), math.Abs(previousPercent2),
			trading.Holding, previousHolding1, previousHolding2)
	}
	return []byte(result)
}

func NewHoldingChangeTxt(holding, previousHolding float64) string {
	if holding > previousHolding {
		return fmt.Sprintf("ARK本日总计买入%.0f股。", holding-previousHolding)
	} else if holding < previousHolding {
		return fmt.Sprintf("ARK本日总计卖出%.0f股。", previousHolding-holding)
	} else {
		return "ARK本日总计持股数未发生变化。"
	}
}

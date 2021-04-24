package service

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/utils"
	"math"
	"strconv"
	"time"
)

type SpecialTradingsReport struct {
	Date             string
	Percent          float64
	tradings         *ARKTradings
	previousTradings []*ARKTradings
}

type StockTxtReport struct {
	Records map[string][]byte
}

func NewStockTxtReport() *StockTxtReport {
	return &StockTxtReport{
		Records: make(map[string][]byte),
	}
}

func (h *StockTxtReport) Add(stock string, record []byte) {
	h.Records[stock] = append(h.Records[stock], record...)
}

func (h *StockTxtReport) Report() []byte {
	var (
		report []byte
	)
	for _, record := range h.Records {
		report = append(report, record...)
	}

	return report
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
		sheet                     = defaultSheet
		idx                       = 2
		higherThan10Report        = NewStockTxtReport()
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
			previousFixedDirection1 := getTradingFixedDirection(r.previousTradings[1], fund, trading.Ticker)
			previousFixedDirection2 := getTradingFixedDirection(r.previousTradings[0], fund, trading.Ticker)
			if previousFixedDirection1 == trading.FixedDirection && previousFixedDirection2 == trading.FixedDirection {
				continuousDirectionReport.Add(trading.Ticker, NewContinuousDirectionSpecialTradingTxtFromTrading(trading,
					getTradingPercent(r.previousTradings[1], fund, trading.Ticker),
					getTradingPercent(r.previousTradings[0], fund, trading.Ticker),
					getTradingHolding(r.previousTradings[1], fund, trading.Ticker),
					getTradingHolding(r.previousTradings[0], fund, trading.Ticker),
				))
			} else {
				if math.Abs(trading.Percent) > 10 {
					higherThan10Report.Add(trading.Ticker, NewSpecialTradingTxtFromTrading(trading))
				}
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

	higherThan10TxtContent := higherThan10Report.Report()
	if len(higherThan10TxtContent) > 0 {
		err = utils.NewFileStoreSvc(r.HigherThan10TxtPath()).Save(higherThan10TxtContent)
		if err != nil {
			glog.Errorf("failed to save txt %s, err: %v", r.HigherThan10TxtPath(), err)
			return err
		}
	} else {
		glog.V(4).Infof("No higher than 10 tradings")
	}

	continuousDirectionTxtContent := continuousDirectionReport.Report()
	if len(continuousDirectionTxtContent) > 0 {
		err = utils.NewFileStoreSvc(r.ContinuousDirectionTxtPath()).Save(continuousDirectionTxtContent)
		if err != nil {
			glog.Errorf("failed to save txt %s, err: %v", r.ContinuousDirectionTxtPath(), err)
			return err
		}
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
			result = fmt.Sprintf("%s在%s中清仓，卖出%.0f股。\n", trading.Ticker, trading.Fund, trading.PreviousHolding)
		} else {
			result = fmt.Sprintf("%s在%s中被减持了%.2f%%，持股数从%.0f股减少到%.0f股。\n", trading.Ticker, trading.Fund,
				math.Abs(trading.Percent), trading.PreviousHolding, trading.Holding)
		}
	}
	return []byte(result)
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

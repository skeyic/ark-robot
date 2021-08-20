package service

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/config"
	"github.com/skeyic/ark-robot/utils"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type StockDateRangeReport struct {
	Ticker     string
	FromDate   time.Time
	EndDate    time.Time
	ReportTime time.Time
	TotalDays  int64

	Funds []string

	// generate by load
	Details *stockDateRangeDetails
}

func NewStockDateRangeReport(ticker string, fromDate, endDate time.Time, funds string) *StockDateRangeReport {
	r := &StockDateRangeReport{
		Ticker:     ticker,
		FromDate:   fromDate,
		EndDate:    endDate,
		ReportTime: time.Now(),
	}
	if funds == "" {
		r.Funds = allARKTypes
	} else {
		r.Funds = strings.Split(funds, ",")
	}
	utils.CheckFolder(r.ReportFolder())
	return r
}

func NewStockDateRangeReportFromDays(ticker string, days int64, funds string) *StockDateRangeReport {
	r := &StockDateRangeReport{
		Ticker:     ticker,
		EndDate:    time.Now(),
		ReportTime: time.Now(),
		TotalDays:  days,
	}
	if funds == "" {
		r.Funds = allARKTypes
	} else {
		r.Funds = strings.Split(funds, ",")
	}
	utils.CheckFolder(r.ReportFolder())
	return r
}

type stockDateRangeDetails struct {
	dateList    []time.Time
	fundList    []string
	dailyDetail map[time.Time]*stockDailyDetail

	tradingSummary *stockDataRangeTradingAnalysis
}

func (r *stockDateRangeDetails) GenerateTradingAnalysis() {
	var (
		tradingAnalysis = &stockDataRangeTradingAnalysis{}
	)

	for _, theDate := range r.dateList {
		tradingAnalysis.AddTrading(theDate, r.dailyDetail[theDate].totalTradingShards)
	}

	r.tradingSummary = tradingAnalysis
}

func (r *stockDateRangeDetails) TxtReport() string {
	var (
		holdingReport = "持仓信息：\n"
		tradingReport string
		keyReport     string
	)

	for dateIdx, theDate := range r.dateList {
		var (
			holdings                = r.dailyDetail[theDate].holdings
			tradings                = r.dailyDetail[theDate].tradings
			totalHoldingShards      = r.dailyDetail[theDate].totalHoldingShards
			totalHoldingMarketValue = r.dailyDetail[theDate].totalHoldingMarketValue
			totalTradingShards      = r.dailyDetail[theDate].totalTradingShards
			today                   = fmt.Sprintf("%d月%d日，", theDate.Month(), theDate.Day())
			txtDailyHoldingReport   = today
			txtDailyHoldingTemp     string
			txtDailyTradingTemp     string
		)

		for idx, fund := range r.fundList {
			var (
				holding *StockHolding
				trading = tradings.GetFundTrading(fund)
			)

			if holdings != nil {
				holding = holdings.GetFundHolding(fund)
			}

			if holding != nil {
				txtDailyHoldingTemp = txtDailyHoldingTemp + fmt.Sprintf("%s持有%s股(比重%.2f%%)，", holding.Fund,
					utils.ThousandFormatFloat64(holding.Shards), holding.Weight)
			}

			if trading != nil {
				if trading.IsBuy() {
					txtDailyTradingTemp += fund + fmt.Sprintf("增持%s股，", utils.ThousandFormatFloat64(trading.Shards))
					if trading.Percent > 50 {
						if trading.Percent == 100 {
							keyReport += fmt.Sprintf("    %s%s建仓\n", theDate.Format(TheDateFormat), fund)
						}
						keyReport += fmt.Sprintf("    %s%s大幅增持%f%%\n", theDate.Format(TheDateFormat), fund, trading.Percent)
					}
				} else if trading.IsSell() {
					if holding == nil {
						txtDailyTradingTemp += fund + fmt.Sprintf("清仓全部的%s股，", utils.ThousandFormatFloat64(-1*trading.Shards))
						keyReport += fmt.Sprintf("    %s%s清仓\n", theDate.Format(TheDateFormat), fund)
					} else {
						txtDailyTradingTemp += fund + fmt.Sprintf("减持%s股，", utils.ThousandFormatFloat64(-1*trading.Shards))
						if trading.Percent < -50 {
							keyReport += fmt.Sprintf("    %s%s大幅减持%f%%\n", theDate.Format(TheDateFormat), fund, trading.Percent*-1)
						}
					}
				}
			}

			// Set the total
			if idx == len(r.fundList)-1 {
				if r.dailyDetail[theDate].totalHoldingShards != 0 {
					txtDailyHoldingTemp = fmt.Sprintf("ARK（%s）共持有%s股，市值为%s美元，其中",
						r.fundList, utils.ThousandFormatFloat64(totalHoldingShards), utils.ThousandFormatFloat64(totalHoldingMarketValue)) + txtDailyHoldingTemp
				} else {
					txtDailyHoldingTemp = "ARK未持有"
				}

				txtDailyTradingTemp += fmt.Sprintf("ARK（%s）总持有股数", r.fundList)
				if totalTradingShards > 0 {
					txtDailyTradingTemp += fmt.Sprintf("增加%s股。\n", utils.ThousandFormatFloat64(totalTradingShards))
				} else if totalTradingShards < 0 {
					txtDailyTradingTemp += fmt.Sprintf("减少%s股。\n", utils.ThousandFormatFloat64(-1*totalTradingShards))
				} else {
					txtDailyTradingTemp += "没有变化。\n"
				}
			}

		}
		if dateIdx == 0 {
			txtDailyHoldingReport = "\t期初" + txtDailyHoldingReport
			txtDailyHoldingReport += strings.TrimSuffix(txtDailyHoldingTemp, "，") + "。"
			holdingReport += txtDailyHoldingReport + "\n"
		} else if dateIdx == len(r.dateList)-1 {
			txtDailyHoldingReport = " \t期末" + txtDailyHoldingReport
			txtDailyHoldingReport += strings.TrimSuffix(txtDailyHoldingTemp, "，") + "。"
			holdingReport += txtDailyHoldingReport + "\n"
		}
		tradingReport += "\t" + today + txtDailyTradingTemp
	}

	if keyReport == "" && r.tradingSummary.TxtKeyReport() == "" {
		keyReport = "重点：无。\n"
	} else {
		keyReport = "重点：\n" + keyReport + r.tradingSummary.TxtKeyReport()
	}

	return keyReport + "\n" + holdingReport + "\n" + r.tradingSummary.TxtReport() + tradingReport
}

func (r *stockDateRangeDetails) TxtReportOld() string {
	var (
		holdingReport string
		tradingReport string
		keyReport     string
	)

	for dateIdx, theDate := range r.dateList {
		var (
			holdings                = r.dailyDetail[theDate].holdings
			tradings                = r.dailyDetail[theDate].tradings
			totalHoldingShards      = r.dailyDetail[theDate].totalHoldingShards
			totalHoldingMarketValue = r.dailyDetail[theDate].totalHoldingMarketValue
			totalTradingShards      = r.dailyDetail[theDate].totalTradingShards
			today                   = fmt.Sprintf("%d月%d日，", theDate.Month(), theDate.Day())
			txtDailyHoldingReport   = today
			txtDailyHoldingTemp     string
			txtDailyTradingTemp     string
		)

		for idx, fund := range r.fundList {
			var (
				holding = holdings.GetFundHolding(fund)
				trading = tradings.GetFundTrading(fund)
			)
			if holding != nil {
				txtDailyHoldingTemp = txtDailyHoldingTemp + fmt.Sprintf("%s持有%s股(比重%.2f%%)，", holding.Fund,
					utils.ThousandFormatFloat64(holding.Shards), holding.Weight)
			}

			if trading.IsBuy() {
				txtDailyTradingTemp += fund + fmt.Sprintf("增持%s股，", utils.ThousandFormatFloat64(trading.Shards))
				if trading.Percent > 50 {
					if trading.Percent == 100 {
						keyReport += fmt.Sprintf("%s%s建仓\n", theDate.Format(TheDateFormat), fund)
					}
					keyReport += fmt.Sprintf("%s%s增持%f%%\n", theDate.Format(TheDateFormat), fund, trading.Percent)
				}
			} else if trading.IsSell() {
				txtDailyTradingTemp += fund + fmt.Sprintf("减持%s股，", utils.ThousandFormatFloat64(trading.Shards*-1))
				if trading.Percent < -50 {
					if trading.Percent == -100 {
						keyReport += fmt.Sprintf("%s%s清仓\n", theDate.Format(TheDateFormat), fund)
					} else {
						keyReport += fmt.Sprintf("%s%s减持%f%%\n", theDate.Format(TheDateFormat), fund, trading.Percent*-1)
					}
				}
			}

			// Set the total
			if idx == len(r.fundList)-1 {
				if r.dailyDetail[theDate].totalHoldingShards != 0 {
					txtDailyHoldingTemp = fmt.Sprintf("ARK共持有%s股，市值为%s美元，其中",
						utils.ThousandFormatFloat64(totalHoldingShards), utils.ThousandFormatFloat64(totalHoldingMarketValue)) + txtDailyHoldingTemp
				} else {
					txtDailyHoldingTemp = "ARK未持有"
				}

				txtDailyTradingTemp += "ARK总持有股数"
				if totalTradingShards > 0 {
					txtDailyTradingTemp += fmt.Sprintf("增加%s股。\n", utils.ThousandFormatFloat64(totalTradingShards))
				} else if totalTradingShards < 0 {
					txtDailyTradingTemp += fmt.Sprintf("减少%s股。\n", utils.ThousandFormatFloat64(-1*totalTradingShards))
				} else {
					txtDailyTradingTemp += "没有变化。\n"
				}
			}

		}
		if dateIdx == 0 {
			txtDailyHoldingReport = "期初" + txtDailyHoldingReport
			txtDailyHoldingReport += strings.TrimSuffix(txtDailyHoldingTemp, "，") + "。"
			holdingReport += txtDailyHoldingReport + "\n"
		} else if dateIdx == len(r.dateList)-1 {
			txtDailyHoldingReport = "期末" + txtDailyHoldingReport
			txtDailyHoldingReport += strings.TrimSuffix(txtDailyHoldingTemp, "，") + "。"
			holdingReport += txtDailyHoldingReport + "\n"
		}
		tradingReport += "  " + today + txtDailyTradingTemp
	}

	if keyReport == "" && r.tradingSummary.TxtKeyReport() == "" {
		keyReport = "重点：无。\n"
	} else {
		keyReport = "重点：\n" + keyReport + r.tradingSummary.TxtKeyReport()
	}

	return keyReport + holdingReport + r.tradingSummary.TxtReport() + tradingReport
}

type stockDailyDetail struct {
	funds []string

	theDate  time.Time
	holdings *StockARKHoldings
	tradings *StockARKTradings

	totalHoldingShards      float64
	totalHoldingMarketValue float64
	totalTradingShards      float64
}

func newStockDailyDetail(funds []string, theDate time.Time, holdings *StockARKHoldings, tradings *StockARKTradings) *stockDailyDetail {
	detail := &stockDailyDetail{
		funds:    funds,
		theDate:  theDate,
		holdings: holdings,
		tradings: tradings,
	}

	detail.Sum()
	return detail
}

func (d *stockDailyDetail) Sum() {
	var (
		totalHoldingShards, totalHoldingMarketValue, totalTradingShards float64
	)
	for _, fund := range d.funds {
		if d.holdings != nil {
			holding := d.holdings.GetFundHolding(fund)
			if holding != nil {
				totalHoldingShards += holding.Shards
				totalHoldingMarketValue += holding.MarketValue
			}
		}
		trading := d.tradings.GetFundTrading(fund)
		if trading != nil {
			totalTradingShards += trading.Shards
		}
	}
	d.totalHoldingShards = totalHoldingShards
	d.totalHoldingMarketValue = totalHoldingMarketValue
	d.totalTradingShards = totalTradingShards
}

func (r *StockDateRangeReport) Load() error {
	var (
		stockDetails = &stockDateRangeDetails{
			dailyDetail: make(map[time.Time]*stockDailyDetail),
		}
	)

	stock := TheStockLibraryMaster.StockLibraries[r.Ticker]
	if stock == nil {
		return errStockNotFound
	}

	var (
		dateList timeList
		byDays   = r.TotalDays != 0
	)

	for theDate := range stock.HistoryStockTradings {
		if byDays {
			if theDate.Equal(r.EndDate) || theDate.Before(r.EndDate) {
				dateList = append(dateList, theDate)
			}
		} else {
			if theDate.Equal(r.FromDate) || theDate.Equal(r.EndDate) || (theDate.After(r.FromDate) && theDate.Before(r.EndDate)) {
				dateList = append(dateList, theDate)
				r.TotalDays++
			}
		}
	}
	if len(dateList) == 0 {
		return errNoDataInDateRange
	}

	sort.Sort(dateList)

	if byDays && len(dateList) > int(r.TotalDays) {
		dateList = dateList[len(dateList)-int(r.TotalDays):]
	}

	r.FromDate = dateList[0]

	var (
		fundList []string
	)

	for _, fund := range r.Funds {
		for i := 0; i < len(dateList); i++ {
			holdings := stock.HistoryStockHoldings[dateList[i]]
			if holdings != nil {
				holding := holdings.GetFundHolding(fund)
				if holding != nil {
					fundList = append(fundList, fund)
					break
				}
			}
		}
	}

	glog.V(4).Infof("%s was holding in %v from %s to %s", r.Ticker, fundList, r.FromDate.Format(TheDateFormat), r.EndDate.Format(TheDateFormat))
	stockDetails.fundList = fundList
	stockDetails.dateList = dateList

	for i := 0; i < len(dateList); i++ {
		var (
			theDate = dateList[i]
		)
		stockDetails.dailyDetail[theDate] = newStockDailyDetail(r.Funds, theDate, stock.HistoryStockHoldings[theDate], stock.HistoryStockTradings[theDate])
	}

	r.Details = stockDetails
	r.Details.GenerateTradingAnalysis()

	return nil
}

type stockDataRangeTradingAnalysis struct {
	maxBuyShards float64
	maxBuyDate   time.Time

	maxSellShards float64
	maxSellDate   time.Time

	buyDays  int
	sellDays int
	keepDays int

	maxContinueBuyDays  int
	maxContinueSellDays int
	continueBuyDays     int
	continueSellDays    int
	lastDirection       TradeDirection

	totalShards float64
}

func (t *stockDataRangeTradingAnalysis) AddTrading(date time.Time, shards float64) {
	glog.V(4).Infof("DATE: %s, SHARDS: %f", date.Format(TheDateFormat), shards)
	if shards > 0 {
		t.buyDays++
		if t.lastDirection == TradeBuy || t.lastDirection == TradeEmpty {
			t.continueBuyDays++
			if t.continueBuyDays > t.maxContinueBuyDays {
				t.maxContinueBuyDays = t.continueBuyDays
			}
		} else {
			if t.lastDirection == TradeSell {
				t.continueSellDays = 0
			}
		}
		t.lastDirection = TradeBuy

		if t.maxBuyShards < shards {
			t.maxBuyShards = shards
			t.maxBuyDate = date
		}
	} else if shards < 0 {
		t.sellDays++
		if t.lastDirection == TradeSell || t.lastDirection == TradeEmpty {
			t.continueSellDays++
			if t.continueSellDays > t.maxContinueSellDays {
				t.maxContinueSellDays = t.continueSellDays
			}
		} else {
			if t.lastDirection == TradeBuy {
				t.continueBuyDays = 0
			}
		}
		t.lastDirection = TradeSell
		if math.Abs(t.maxSellShards) < math.Abs(shards) {
			t.maxSellShards = shards
			t.maxSellDate = date
		}
	} else {
		t.keepDays++
	}

	t.totalShards += shards
}

func (t *stockDataRangeTradingAnalysis) TxtKeyReport() string {
	var (
		msg = ""
	)

	if t.continueBuyDays > 2 {
		msg += fmt.Sprintf("    连续%d日获增持。\n", t.continueBuyDays)
	}

	if t.continueSellDays > 2 {
		msg += fmt.Sprintf("    连续%d日被减持。\n", t.continueSellDays)
	}

	return msg
}

func (t *stockDataRangeTradingAnalysis) TxtReport() string {
	var (
		msg = "交易信息：\n    本期总计"
	)

	if t.totalShards >= 0 {
		msg += fmt.Sprintf("增持%s股，", utils.ThousandFormatFloat64(t.totalShards))
	} else if t.totalShards < 0 {
		msg += fmt.Sprintf("减持%s股，", utils.ThousandFormatFloat64(-1*t.totalShards))
	}

	msg += fmt.Sprintf("共计增持%d日，减持%d日，没有变动%d日。", t.buyDays, t.sellDays, t.keepDays)

	if t.buyDays > 0 {
		msg += fmt.Sprintf("最大增持发生在%d月%d日，增持%s股。", t.maxBuyDate.Month(), t.maxBuyDate.Day(),
			utils.ThousandFormatFloat64(t.maxBuyShards))
	}

	if t.sellDays > 0 {
		msg += fmt.Sprintf("最大减持发生在%d月%d日，减持%s股。", t.maxSellDate.Month(), t.maxSellDate.Day(),
			utils.ThousandFormatFloat64(-1*t.maxSellShards))
	}

	msg += "\n    具体如下：\n"

	return msg
}

func (r *StockDateRangeReport) ReportExcel() error {
	var (
		err      error
		fileName = r.ExcelPath()
	)

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
		sheet = defaultSheet

		dateIdxList = []string{"B", "C", "D", "E", "F", "G",
			"H", "I", "J", "K", "L", "M", "N", "O", "P",
			"Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
			"AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI",
			"AJ", "AK", "AL", "AM", "AN", "AO", "AP", "AQ", "AR",
			"AS", "AT", "AU", "AV", "AW", "AX", "AY", "AZ",
		}
		maxDate = len(dateIdxList)
	)

	for dateIdx, theDate := range r.Details.dateList {
		if dateIdx > maxDate {
			glog.Warningf("MAX DATE RANGE REACHED, we have: %d, max: %d", len(r.Details.dateList), maxDate)
			break
		}

		var (
			totalShards        float64
			totalMarketValue   float64
			totalTradingShards float64
			fundIdx            = 25
			holdings           = r.Details.dailyDetail[theDate].holdings
			tradings           = r.Details.dailyDetail[theDate].tradings
		)

		for idx, fund := range r.Details.fundList {
			var (
				holding *StockHolding
			)
			if holdings != nil {
				holding = holdings.GetFundHolding(fund)
			}
			var currentShards float64

			line := strconv.Itoa(fundIdx)
			// Set the date column
			if idx == 0 {
				f.SetCellValue(sheet, dateIdxList[dateIdx]+line, theDate.Format(TheDateFormat))
				fundIdx++
			}

			line = strconv.Itoa(fundIdx)
			if dateIdx == 0 {
				f.SetCellValue(sheet, "A"+line, fund)
			}
			if holding != nil {
				totalShards += holding.Shards
				totalMarketValue += holding.MarketValue
				currentShards = holding.Shards
			} else {
				//txtDailyHoldingTemp = txtDailyHoldingTemp + fmt.Sprintf("%s持有%.0f股(比重%.2f%%)，", holding.Fund, holding.Shards, holding.Weight)
			}
			f.SetCellValue(sheet, dateIdxList[dateIdx]+line, fmt.Sprintf("%.0f", currentShards))
			fundIdx++

			trading := tradings.GetFundTrading(fund)
			if trading != nil {
				//glog.V(4).Infof("TRADING: %v", trading)
				totalTradingShards += trading.Shards
			}

			// Set the total
			if idx == len(r.Details.fundList)-1 {
				line = strconv.Itoa(fundIdx)
				if dateIdx == 0 {
					f.SetCellValue(sheet, "A"+line, "TOTAL")
				}
				f.SetCellValue(sheet, dateIdxList[dateIdx]+line, fmt.Sprintf("%.0f", totalShards))
			}
		}

	}

	err = f.Save()
	if err != nil {
		glog.Errorf("failed to save excel %s, err: %v", fileName, err)
		return err
	}

	//glog.V(4).Infof("%s", txtReport)
	//glog.V(4).Infof("%s", txtTradingReport)

	return nil
}

func (r *StockDateRangeReport) ReportImage() error {
	var (
		dates               []string
		currentHoldingsData []float64
	)

	for _, theDate := range r.Details.dateList {
		dates = append(dates, theDate.Format(TheDateFormat))
		currentHoldingsData = append(currentHoldingsData, r.Details.dailyDetail[theDate].totalHoldingShards)
	}

	// create a bar and line
	var (
		bar  = charts.NewBar()
		line = charts.NewLine()
	)

	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: `ARK持仓中` + r.Ticker + fmt.Sprintf("持股数（%d月%d日至%d月%d日）",
				r.FromDate.Month(), r.FromDate.Day(), r.EndDate.Month(), r.EndDate.Day()),
			//Top: "5%",
			//Bottom: "20%",
			Left: "center",
			//Right: "20%",
		}), charts.WithLegendOpts(opts.Legend{
			Show: true,
			Top:  "7%",
		}))
	bar.SetXAxis(dates).AddSeries("当前持股数", utils.ToBarData("Current", currentHoldingsData))
	line.SetXAxis(dates).AddSeries("", utils.ToPercentLineData("Current", currentHoldingsData, 2.0),
		charts.WithLineStyleOpts(opts.LineStyle{
			//Color: "white",
			Width: 2,
		}))
	bar.Overlap(line)

	var (
		htmlPath  = r.HtmlPath()
		imagePath = r.ImagePath()
	)
	f, err := os.Create(htmlPath)
	if err != nil {
		glog.Errorf("failed to create html file %s", htmlPath)
		return err
	}
	err = bar.Render(f)
	if err != nil {
		glog.Errorf("failed to render html file %s", htmlPath)
		return err
	}

	// TODO do not use chrome to generate image, will add another micro service install
	err = utils.TheChartPainter.GenerateImage(htmlPath, imagePath)
	if err != nil {
		glog.Errorf("failed to save image file %s", imagePath)
		return err
	}

	return nil
}

func (r *StockDateRangeReport) Report() error {
	var (
		err error
	)

	err = r.Load()
	if err != nil {
		return err
	}

	var (
		txtReport = `对ARK持仓中` + r.Ticker + fmt.Sprintf("（基金：%s, 日期：%d月%d日至%d月%d日）的分析: \n",
			r.Funds, r.FromDate.Month(), r.FromDate.Day(), r.EndDate.Month(), r.EndDate.Day())
	)

	if config.Config.Report.WithExcel {
		err = r.ReportExcel()
		if err != nil {
			glog.Errorf("failed to report excel, err: %v", err)
			return err
		}
	}

	err = r.ReportImage()
	if err != nil {
		glog.Errorf("failed to report image, err: %v", err)
		return err
	}

	err = utils.NewFileStoreSvc(r.TxtPath()).Save([]byte(txtReport + r.Details.TxtReport()))
	if err != nil {
		glog.Errorf("failed to save txt %s, err: %v", r.TxtPath(), err)
		return err
	}

	return nil
}

//func (r *StockDateRangeReport) OldReport() error {
//	var (
//		err       error
//		fileName  = r.ExcelPath()
//		txtReport = `对ARK持仓中` + r.Ticker + fmt.Sprintf("（%d月%d日至%d月%d日）的分析: \n",
//			r.FromDate.Month(), r.FromDate.Day(), r.EndDate.Month(), r.EndDate.Day())
//		theTradingReport = &stockDataRangeTradingAnalysis{}
//		txtTradingReport string
//	)
//
//	err = r.Load()
//	if err != nil {
//		return err
//	}
//
//	err = r.InitExcelFromTemplate()
//	if err != nil {
//		glog.Errorf("failed to init excel from template, err: %v", err)
//		return err
//	}
//
//	f, err := excelize.OpenFile(fileName)
//	if err != nil {
//		glog.Errorf("failed to open excel %s, err: %v", fileName, err)
//		return err
//	}
//
//	/*
//				    A           B           C           D           E           F
//				25		        2021-03-29	2021-03-30	2021-03-31	2021-04-01	2021-04-02
//				26	ARKW持仓变动	10013448 	102200000 	10013408 	10002000 	10013448
//				27	ARKF持仓变动	100230000 	10001245 	10013448 	100014000 	100014000
//				28	ARK总持仓变动	110243448 	112201245 	20026856 	110016000 	110027448
//
//			txtReport
//			关于HUYA，从4月5日至4月9日:
//
//			txtDailyHoldingReport
//		    期初4月5日，ARK共持有10023448股，市值为1234567美元，其中ARKK持有500000股（比重1%），ARKW持有500000股（比重1%），期末4月9日，ARK共持有10002485股，市值为1234567美元，其中ARKK400000股（比重1%），ARKW400000股（比重1%）。
//
//			txtDailyTradingReport
//			本期总计减持/增持20963股，共计减持2日，增持2日，没有变动1日。最大减持发生在4月8日，减持30000股，最大增持发生在4月6日，增持30000股。
//			具体情况如下：
//
//			txtDailyTradingDetail
//			  4月5日，ARKK增持10000股，ARKW减持20000股，ARK总持有股数减少10000股；
//			  4月6日，ARKK增持10000股，ARKW增持20000股，ARK总持有股数增加30000股；
//			  4月7日，ARKK增持20000股，ARKW没有变动，ARK总持有股数增加20000股；
//			  4月8日，ARKK减持10000股，ARKW减持20000股，ARK总持股数有减少30000股；
//			  4月9日，ARK总持股数没有变动。
//	*/
//
//	var (
//		sheet = defaultSheet
//
//		dateIdxList = []string{"B", "C", "D", "E", "F", "G",
//			"H", "I", "J", "K", "L", "M", "N", "O", "P",
//			"Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
//		maxDate = len(dateIdxList)
//	)
//
//	for dateIdx, theDate := range r.Details.dateList {
//		if dateIdx > maxDate {
//			glog.Warningf("MAX DATE RANGE REACHED, we have: %d, max: %d", len(r.Details.dateList), maxDate)
//			break
//		}
//
//		var (
//			totalShards           float64
//			totalMarketValue      float64
//			totalTradingShards    float64
//			fundIdx               = 25
//			holdings              = r.Details.dailyDetail[theDate].holdings
//			tradings              = r.Details.dailyDetail[theDate].tradings
//			today                 = fmt.Sprintf("%d月%d日，", theDate.Month(), theDate.Day())
//			txtDailyHoldingReport = today
//			txtDailyHoldingTemp   string
//			txtDailyTradingTemp   string
//		)
//
//		for idx, fund := range r.Details.fundList {
//			holding := holdings[fund]
//			var currentShards float64
//
//			line := strconv.Itoa(fundIdx)
//			// Set the date column
//			if idx == 0 {
//				f.SetCellValue(sheet, dateIdxList[dateIdx]+line, theDate.Format(TheDateFormat))
//				fundIdx++
//			}
//
//			line = strconv.Itoa(fundIdx)
//			if dateIdx == 0 {
//				f.SetCellValue(sheet, "A"+line, fund)
//			}
//			if holding != nil {
//				totalShards += holding.Shards
//				totalMarketValue += holding.MarketValue
//				currentShards = holding.Shards
//				txtDailyHoldingTemp = txtDailyHoldingTemp + fmt.Sprintf("%s持有%s股(比重%.2f%%)，", holding.Fund,
//					utils.ThousandFormatFloat64(holding.Shards), holding.Weight)
//			} else {
//				//txtDailyHoldingTemp = txtDailyHoldingTemp + fmt.Sprintf("%s持有%.0f股(比重%.2f%%)，", holding.Fund, holding.Shards, holding.Weight)
//			}
//			f.SetCellValue(sheet, dateIdxList[dateIdx]+line, fmt.Sprintf("%.0f", currentShards))
//			fundIdx++
//
//			trading := tradings[fund]
//			//glog.V(4).Infof("TRADING: %v", trading)
//			totalTradingShards += trading.Shards
//
//			if trading.IsBuy() {
//				txtDailyTradingTemp += fund + fmt.Sprintf("增持%s股，", utils.ThousandFormatFloat64(trading.Shards))
//			} else if trading.IsSell() {
//				txtDailyTradingTemp += fund + fmt.Sprintf("减持%s股，", utils.ThousandFormatFloat64(-1*trading.Shards))
//			}
//
//			// Set the total
//			if idx == len(r.Details.fundList)-1 {
//				line = strconv.Itoa(fundIdx)
//				if dateIdx == 0 {
//					f.SetCellValue(sheet, "A"+line, "TOTAL")
//				}
//				f.SetCellValue(sheet, dateIdxList[dateIdx]+line, fmt.Sprintf("%.0f", totalShards))
//				if totalShards != 0 {
//					txtDailyHoldingTemp = fmt.Sprintf("ARK共持有%s股，市值为%s美元，其中",
//						utils.ThousandFormatFloat64(totalShards), utils.ThousandFormatFloat64(totalMarketValue)) + txtDailyHoldingTemp
//				} else {
//					txtDailyHoldingTemp = "ARK未持有"
//				}
//
//				theTradingReport.AddTrading(theDate, totalTradingShards)
//				txtDailyTradingTemp += "ARK总持有股数"
//				if totalTradingShards > 0 {
//					txtDailyTradingTemp += fmt.Sprintf("增加%s股。\n", utils.ThousandFormatFloat64(totalTradingShards))
//				} else if totalTradingShards < 0 {
//					txtDailyTradingTemp += fmt.Sprintf("减少%s股。\n", utils.ThousandFormatFloat64(-1*totalTradingShards))
//				} else {
//					txtDailyTradingTemp += "没有变化。\n"
//				}
//			}
//		}
//
//		if dateIdx == 0 {
//			txtDailyHoldingReport = "期初" + txtDailyHoldingReport
//			txtDailyHoldingReport += txtDailyHoldingTemp
//			txtReport += txtDailyHoldingReport + "\n"
//		} else if dateIdx == len(r.Details.dateList)-1 {
//			txtDailyHoldingReport = "期末" + txtDailyHoldingReport
//			txtDailyHoldingReport += txtDailyHoldingTemp
//			txtReport += txtDailyHoldingReport + "\n"
//		}
//
//		txtTradingReport += "  " + today + txtDailyTradingTemp
//	}
//
//	txtTradingReport = theTradingReport.TxtReport() + txtTradingReport
//
//	err = f.Save()
//	if err != nil {
//		glog.Errorf("failed to save excel %s, err: %v", fileName, err)
//		return err
//	}
//
//	err = utils.NewFileStoreSvc(r.TxtPath()).Save([]byte(txtReport + txtTradingReport))
//	if err != nil {
//		glog.Errorf("failed to save txt %s, err: %v", r.TxtPath(), err)
//		return err
//	}
//
//	//glog.V(4).Infof("%s", txtReport)
//	//glog.V(4).Infof("%s", txtTradingReport)
//
//	return nil
//}

func (r *StockDateRangeReport) ReportFolder() string {
	return stockReportPath + "/" + r.Ticker + "_" + r.ReportTime.Format("2006-01-02-15-04-05")
}

func (r *StockDateRangeReport) ExcelPath() string {
	return r.ReportFolder() + "/" + r.ExcelName()
}

func (r *StockDateRangeReport) ExcelName() string {
	return r.FileName() + ".xlsx"
}

func (r *StockDateRangeReport) TxtPath() string {
	return r.ReportFolder() + "/" + r.TxtName()
}

func (r *StockDateRangeReport) TxtName() string {
	return r.FileName() + ".txt"
}

func (r *StockDateRangeReport) HtmlPath() string {
	return r.ReportFolder() + "/" + r.HtmlName()
}

func (r *StockDateRangeReport) HtmlName() string {
	return r.FileName() + ".html"
}

func (r *StockDateRangeReport) ImagePath() string {
	return r.ReportFolder() + "/" + r.ImageName()
}

func (r *StockDateRangeReport) ImageName() string {
	return r.FileName() + ".png"
}

func (r *StockDateRangeReport) FileName() string {
	return fmt.Sprintf("%s%s_from_%s_to_%s", prefixDateRangeStockReport, r.Ticker,
		r.FromDate.Format(TheDateFormat), r.EndDate.Format(TheDateFormat))
}

func (r *StockDateRangeReport) InitExcelFromTemplate() error {
	var fileName = r.ExcelPath()
	if utils.CheckFileExist(fileName) {
		utils.DeleteFile(fileName)
	}
	utils.CopyFile(stockReportExcelTemplate, fileName)
	glog.V(4).Infof("Init fileName: %s", fileName)
	return nil
}

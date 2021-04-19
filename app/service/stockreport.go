package service

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/utils"
	"math"
	"sort"
	"strconv"
	"time"
)

type StockReport struct {
	Ticker   string
	FromDate time.Time
	EndDate  time.Time
	Details  *stockDetails
}

func NewStockReport(ticker string, fromDate, endDate time.Time) *StockReport {
	return &StockReport{
		Ticker:   ticker,
		FromDate: fromDate,
		EndDate:  endDate,
	}
}

type stockDetails struct {
	dateList    []time.Time
	fundList    []string
	dailyDetail map[time.Time]*stockDailyDetail
}

type stockDailyDetail struct {
	date     time.Time
	holdings map[string]*StockHolding
	tradings map[string]*StockTrading
}

func (r *StockReport) Load() error {
	var (
		stockDetails = &stockDetails{
			dailyDetail: make(map[time.Time]*stockDailyDetail),
		}
	)

	stock := TheStockLibraryMaster.StockLibraries[r.Ticker]
	if stock == nil {
		return errStockNotFound
	}

	var (
		dateList timeList
	)

	for theDate := range stock.HistoryStockHoldings {
		if theDate.Equal(r.FromDate) || theDate.Equal(r.EndDate) || (theDate.After(r.FromDate) && theDate.Before(r.EndDate)) {
			dateList = append(dateList, theDate)
		}
	}
	if len(dateList) == 0 {
		return errNoDataInDateRange
	}

	sort.Sort(dateList)

	var (
		fundList []string
	)

	for _, fund := range allARKTypes {
		for i := 0; i < len(dateList); i++ {
			holdings := stock.HistoryStockHoldings[dateList[i]]

			holding := holdings[fund]
			if holding != nil {
				fundList = append(fundList, fund)
				break
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
		stockDetails.dailyDetail[theDate] = &stockDailyDetail{
			date:     theDate,
			holdings: stock.HistoryStockHoldings[theDate],
			tradings: stock.HistoryStockTradings[theDate],
		}
	}

	r.Details = stockDetails
	return nil
}

type tradingReport struct {
	maxBuyShards float64
	maxBuyDate   time.Time

	maxSellShards float64
	maxSellDate   time.Time

	buyDays  int
	sellDays int
	keepDays int

	totalShards float64
}

func (t *tradingReport) AddTrading(date time.Time, shards float64) {
	if shards > 0 {
		t.buyDays++
		if t.maxBuyShards < shards {
			t.maxBuyShards = shards
			t.maxBuyDate = date
		}
	} else if shards < 0 {
		t.sellDays++
		if math.Abs(t.maxSellShards) < math.Abs(shards) {
			t.maxSellShards = shards
			t.maxSellDate = date
		}
	} else {
		t.keepDays++
	}

	t.totalShards += shards
}

func (t *tradingReport) TxtReport() string {
	var (
		msg = "本期总计"
	)

	if t.totalShards >= 0 {
		msg += fmt.Sprintf("增持%.0f股，", t.totalShards)
	} else if t.totalShards < 0 {
		msg += fmt.Sprintf("减持%.0f股，", -1*t.totalShards)
	}

	msg += fmt.Sprintf("共计增持%d日，减持%d日，没有变动%d日。", t.buyDays, t.sellDays, t.keepDays)

	if t.buyDays > 0 {
		msg += fmt.Sprintf("最大增持发生在%d月%d日，增持了%.0f股。", t.maxBuyDate.Month(), t.maxBuyDate.Day(), t.maxBuyShards)
	}

	if t.sellDays > 0 {
		msg += fmt.Sprintf("最大减持发生在%d月%d日，减持了%.0f股。", t.maxSellDate.Month(), t.maxSellDate.Day(), -1*t.maxSellShards)
	}

	msg += "\n\n具体如下：\n"

	return msg
}

func (r *StockReport) ToExcel() error {
	var (
		err       error
		fileName  = r.ExcelPath()
		txtReport = `对ARK持仓中` + r.Ticker + fmt.Sprintf("（%d月%d日至%d月%d日）的分析: \n",
			r.FromDate.Month(), r.FromDate.Day(), r.EndDate.Month(), r.EndDate.Day())
		theTradingReport = &tradingReport{}
		txtTradingReport string
	)

	err = r.Load()
	if err != nil {
		return err
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

	/*
				    A           B           C           D           E           F
				25		        2021-03-29	2021-03-30	2021-03-31	2021-04-01	2021-04-02
				26	ARKW持仓变动	10013448 	102200000 	10013408 	10002000 	10013448
				27	ARKF持仓变动	100230000 	10001245 	10013448 	100014000 	100014000
				28	ARK总持仓变动	110243448 	112201245 	20026856 	110016000 	110027448

			txtReport
			关于HUYA，从4月5日至4月9日:

			txtDailyHoldingReport
		    期初4月5日，ARK共持有10023448股，其中ARKK持有500000股（比重1%），ARKW持有500000股（比重1%），期末4月9日，ARK共持有10002485股，其中ARKK400000股（比重1%），ARKW400000股（比重1%）。

			txtDailyTradingReport
			本期总计减持/增持了20963股，共计减持2日，增持2日，没有变动1日。最大减持发生在4月8日，减持了30000股，最大增持发生在4月6日，增持了30000股。
			具体情况如下：

			txtDailyTradingDetail
			  4月5日，ARKK增持10000股，ARKW减持20000股，ARK总持有股数减少了10000股；
			  4月6日，ARKK增持10000股，ARKW增持20000股，ARK总持有股数增加了30000股；
			  4月7日，ARKK增持20000股，ARKW没有变动，ARK总持有股数增加了20000股；
			  4月8日，ARKK减持10000股，ARKW减持20000股，ARK总持股数有减少了30000股；
			  4月9日，ARK总持股数没有变动。
	*/

	var (
		sheet = defaultSheet

		dateIdxList = []string{"B", "C", "D", "E", "F", "G",
			"H", "I", "J", "K", "L", "M", "N", "O", "P",
			"Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
		maxDate = len(dateIdxList)
	)

	for dateIdx, theDate := range r.Details.dateList {
		if dateIdx > maxDate {
			glog.Warningf("MAX DATE RANGE REACHED, we have: %d, max: %d", len(r.Details.dateList), maxDate)
			break
		}

		var (
			totalShards           float64
			totalTradingShards    float64
			fundIdx               = 25
			holdings              = r.Details.dailyDetail[theDate].holdings
			tradings              = r.Details.dailyDetail[theDate].tradings
			today                 = fmt.Sprintf("%d月%d日，", theDate.Month(), theDate.Day())
			txtDailyHoldingReport = today
			txtDailyHoldingTemp   string
			txtDailyTradingTemp   string
		)

		for idx, fund := range r.Details.fundList {
			holding := holdings[fund]
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
				currentShards = holding.Shards
				txtDailyHoldingTemp = txtDailyHoldingTemp + fmt.Sprintf("%s持有%.0f股(比重%.2f%%)，", holding.Fund, holding.Shards, holding.Weight)
			} else {
				//txtDailyHoldingTemp = txtDailyHoldingTemp + fmt.Sprintf("%s持有%.0f股(比重%.2f%%)，", holding.Fund, holding.Shards, holding.Weight)
			}
			f.SetCellValue(sheet, dateIdxList[dateIdx]+line, fmt.Sprintf("%.0f", currentShards))
			fundIdx++

			trading := tradings[fund]
			//glog.V(4).Infof("TRADING: %v", trading)
			totalTradingShards += trading.Shards

			if trading.IsBuy() {
				txtDailyTradingTemp += fund + fmt.Sprintf("增持了%.0f股，", trading.Shards)
			} else if trading.IsSell() {
				txtDailyTradingTemp += fund + fmt.Sprintf("减持了%.0f股，", -1*trading.Shards)
			}

			// Set the total
			if idx == len(r.Details.fundList)-1 {
				line = strconv.Itoa(fundIdx)
				if dateIdx == 0 {
					f.SetCellValue(sheet, "A"+line, "TOTAL")
				}
				f.SetCellValue(sheet, dateIdxList[dateIdx]+line, fmt.Sprintf("%.0f", totalShards))
				if totalShards != 0 {
					txtDailyHoldingTemp = fmt.Sprintf("ARK共持有%.0f股，其中", totalShards) + txtDailyHoldingTemp
				} else {
					txtDailyHoldingTemp = "ARK未持有"
				}

				theTradingReport.AddTrading(theDate, totalTradingShards)
				txtDailyTradingTemp += "ARK总持有股数"
				if totalTradingShards > 0 {
					txtDailyTradingTemp += fmt.Sprintf("增加了%.0f股。\n", totalTradingShards)
				} else if totalTradingShards < 0 {
					txtDailyTradingTemp += fmt.Sprintf("减少了%.0f股。\n", -1*totalTradingShards)
				} else {
					txtDailyTradingTemp += "没有变化。\n"
				}
			}
		}

		if dateIdx == 0 {
			txtDailyHoldingReport = "期初" + txtDailyHoldingReport
			txtDailyHoldingReport += txtDailyHoldingTemp
			txtReport += txtDailyHoldingReport + "\n"
		} else if dateIdx == len(r.Details.dateList)-1 {
			txtDailyHoldingReport = "期末" + txtDailyHoldingReport
			txtDailyHoldingReport += txtDailyHoldingTemp
			txtReport += txtDailyHoldingReport + "\n"
		}

		txtTradingReport += "  " + today + txtDailyTradingTemp
	}

	txtTradingReport = theTradingReport.TxtReport() + txtTradingReport

	err = f.Save()
	if err != nil {
		glog.Errorf("failed to save excel %s, err: %v", fileName, err)
		return err
	}

	err = utils.NewFileStoreSvc(r.TxtPath()).Save([]byte(txtReport + txtTradingReport))
	if err != nil {
		glog.Errorf("failed to save txt %s, err: %v", r.TxtPath(), err)
		return err
	}

	//glog.V(4).Infof("%s", txtReport)
	//glog.V(4).Infof("%s", txtTradingReport)

	return nil
}

func (r *StockReport) ReportFolder() string {
	return stockReportPath
}

func (r *StockReport) ExcelPath() string {
	return r.ReportFolder() + "/" + r.ExcelName()
}

func (r *StockReport) ExcelName() string {
	return r.FileName() + ".xlsx"
}

func (r *StockReport) TxtPath() string {
	return r.ReportFolder() + "/" + r.TxtName()
}

func (r *StockReport) TxtName() string {
	return r.FileName() + ".txt"
}

func (r *StockReport) FileName() string {
	return fmt.Sprintf("%s_%s%s_from_%s_to_%s", time.Now().Format("20062102150405"), prefixStockReport, r.Ticker,
		r.FromDate.Format(TheDateFormat), r.EndDate.Format(TheDateFormat))
}

func (r *StockReport) InitExcelFromTemplate() error {
	var fileName = r.ExcelPath()
	if utils.CheckFileExist(fileName) {
		utils.DeleteFile(fileName)
	}
	utils.CopyFile(stockReportExcelTemplate, fileName)
	glog.V(4).Infof("Init fileName: %s", fileName)
	return nil
}

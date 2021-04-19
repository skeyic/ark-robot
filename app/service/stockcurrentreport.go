package service

import (
	"time"
)

type StockCurrentReport struct {
	Ticker  string
	EndDate time.Time
	Details *stockDateRangeDetails
}

func NewStockCurrentReport(ticker string, fromDate, endDate time.Time) *StockCurrentReport {
	return &StockCurrentReport{
		Ticker:  ticker,
		EndDate: time.Now(),
	}
}

//func (r *StockCurrentReport) Load() error {
//	var (
//		stockDetails = &stockDateRangeDetails{
//			dailyDetail: make(map[time.Time]*stockDailyDetail),
//		}
//	)
//
//	stock := TheStockLibraryMaster.StockLibraries[r.Ticker]
//	if stock == nil {
//		return errStockNotFound
//	}
//
//	var (
//		dateList timeList
//	)
//
//	for theDate := range stock.HistoryStockHoldings {
//		if theDate.Equal(r.FromDate) || theDate.Equal(r.EndDate) || (theDate.After(r.FromDate) && theDate.Before(r.EndDate)) {
//			dateList = append(dateList, theDate)
//		}
//	}
//	if len(dateList) == 0 {
//		return errNoDataInDateRange
//	}
//
//	sort.Sort(dateList)
//
//	var (
//		fundList []string
//	)
//
//	for _, fund := range allARKTypes {
//		for i := 0; i < len(dateList); i++ {
//			holdings := stock.HistoryStockHoldings[dateList[i]]
//
//			holding := holdings[fund]
//			if holding != nil {
//				fundList = append(fundList, fund)
//				break
//			}
//		}
//	}
//
//	glog.V(4).Infof("%s was holding in %v from %s to %s", r.Ticker, fundList, r.FromDate.Format(TheDateFormat), r.EndDate.Format(TheDateFormat))
//	stockDetails.fundList = fundList
//	stockDetails.dateList = dateList
//
//	for i := 0; i < len(dateList); i++ {
//		var (
//			theDate = dateList[i]
//		)
//		stockDetails.dailyDetail[theDate] = &stockDailyDetail{
//			date:     theDate,
//			holdings: stock.HistoryStockHoldings[theDate],
//			tradings: stock.HistoryStockTradings[theDate],
//		}
//	}
//
//	r.Details = stockDetails
//	return nil
//}
//
//func (r *StockCurrentReport) Report() error {
//	var (
//		err       error
//		fileName  = r.ExcelPath()
//		txtReport = `对ARK持仓中` + r.Ticker + fmt.Sprintf("（%d月%d日至%d月%d日）的分析: \n",
//			r.FromDate.Month(), r.FromDate.Day(), r.EndDate.Month(), r.EndDate.Day())
//		theTradingReport = &tradingReport{}
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
//		截止4月18日，ARK持有HUYA共计10023448股，市值200000美元，其中ARKK持有500000股（占比1%，市值10000美元），ARKW持有500000股（占比1%，市值10000美元）。
//		最近一次操作发生在4月17日，ARK总持股从90000股减少到80000股，其中ARKK增持10000股（占比1%），持股数从40000股增加到50000股；ARKW减持20000股（占比1%），持股数从50000股减少到30000股。
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
//
//func (r *StockCurrentReport) ReportFolder() string {
//	return stockReportPath
//}
//
//func (r *StockCurrentReport) ExcelPath() string {
//	return r.ReportFolder() + "/" + r.ExcelName()
//}
//
//func (r *StockCurrentReport) ExcelName() string {
//	return r.FileName() + ".xlsx"
//}
//
//func (r *StockCurrentReport) TxtPath() string {
//	return r.ReportFolder() + "/" + r.TxtName()
//}
//
//func (r *StockCurrentReport) TxtName() string {
//	return r.FileName() + ".txt"
//}
//
//func (r *StockCurrentReport) FileName() string {
//	return fmt.Sprintf("%s_%s%s_from_%s_to_%s", time.Now().Format("20062102150405"), prefixStockReport, r.Ticker,
//		r.FromDate.Format(TheDateFormat), r.EndDate.Format(TheDateFormat))
//}
//
//func (r *StockCurrentReport) InitExcelFromTemplate() error {
//	var fileName = r.ExcelPath()
//	if utils.CheckFileExist(fileName) {
//		utils.DeleteFile(fileName)
//	}
//	utils.CopyFile(stockReportExcelTemplate, fileName)
//	glog.V(4).Infof("Init fileName: %s", fileName)
//	return nil
//}

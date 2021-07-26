package service

import (
	"time"
)

// ChinaStockHoldingReport ...
// Current holding - phase 1
// Latest trading - TODO
// Data range report of last 5 days - phase 1
type ChinaStockHoldingReport struct {
	Ticker          string
	ReportDate      time.Time
	CurrentHolding  *ARKHoldings
	PreviousHolding *ARKHoldings
}

func NewChinaStockHoldingReport(ticker string) *ChinaStockHoldingReport {
	return &ChinaStockHoldingReport{
		Ticker:     ticker,
		ReportDate: time.Now(),
	}
}

func (r *ChinaStockHoldingReport) Load() error {
	//var (
	//	err error
	//)
	//
	//latestDate := TheLibrary.GetLatestHoldingDate()
	//if latestDate.IsZero() {
	//	return errNoLatestDate
	//}
	//
	//latestStockHolding := TheStockLibraryMaster.GetStockLatestHolding(r.Ticker)
	//if latestStockHolding.Date != latestDate {
	//	return errStockNotHold
	//}
	//r.CurrentHolding = latestStockHolding
	//
	//r.DataRangeReport = NewStockDateRangeReportFromDays(r.Ticker, HistoryDays)
	//err = r.DataRangeReport.Load()
	//if err != nil {
	//	glog.Errorf("failed to load the data range report, ticker: %s, from: %s, end: %s, days: %d, err: %v",
	//		r.Ticker, r.DataRangeReport.FromDate, r.DataRangeReport.EndDate, r.DataRangeReport.TotalDays, err)
	//	return err
	//}
	//
	//allDates := r.DataRangeReport.Details.dateList
	//if allDates[len(allDates)-1] != latestDate {
	//	glog.Warningf("The stock %s was not in the ARK holding of %s", r.Ticker, latestDate)
	//	return errStockNotHold
	//}

	return nil
}

//func (r *ChinaStockHoldingReport) TxtReport() string {
//	var (
//		report = r.CurrentHolding.TxtReport()
//	)
//
//	report += "\n分析最近五个交易日的数据：\n" + r.DataRangeReport.Details.TxtReport()
//	return report
//}

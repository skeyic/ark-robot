package service

import (
	"time"
)

// ChinaStockHoldingReport ...
// Current holding - phase 1
// Latest trading - TODO
// Data range report of last 5 days - phase 1
type ChinaStockHoldingReport struct {
	ReportDate      time.Time
	CurrentHolding  *ARKHoldings
	PreviousHolding *ARKHoldings
}

func NewChinaStockHoldingReport(ticker string) *ChinaStockHoldingReport {
	return &ChinaStockHoldingReport{
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

/*
ARK持有的中概股总市值为xx美元，相比上个交易日增加/减少了xxx美元，持有市值最多的是xx，共计xxx美元。
其中：
  建仓了xxx，市值xxx美元；
  清仓了xxx；
  增持最多的是xxx，增加了xxx美元；
  减持最多的是xxx，减少了xxx美元。
具体如下：
  TCEHY：持有xxx股，市值xxx美元，相比上个交易日增加了xxx股；
  JD：持有xxx股，市值xxx美元，相比上个交易日减少了xxx股；
  BIDU：持有xxx股，市值xxx美元，相比上个交易日没有变化；
*/
//func (r *ChinaStockHoldingReport) TxtReport() string {
//	var (
//		report string
//		totalMarketValue, previousMarketValue float64
//		maxMarketValue float64
//		maxMarketValueTicker string
//		maxBuyTicker, maxSoldTicker string
//		firstBuyReport, soldOutReport, detailReport string
//	)
//
//	for _, fund := range allARKTypes {
//		holdings := r.CurrentHolding.GetFundStockHoldings(fund)
//		for ticker, holding := range holdings.Holdings {
//			if TheChinaStockManager.IsChinaStock(ticker) {
//				totalMarketValue += holding.MarketValue
//				if holding.MarketValue > maxMarketValue {
//					maxMarketValue = holding.MarketValue
//					maxMarketValueTicker = ticker
//				}
//
//				detailReport += fmt.Sprintf("  %s：持有%s股，市值%s美元，相比上个交易日", ticker,
//					utils.ThousandFormatFloat64(holding.Shards), utils.ThousandFormatFloat64(holding.MarketValue))
//
//				previousHolding := r.PreviousHolding.GetFundStockHoldings(fund).GetStockHolding(ticker)
//				if previousHolding.Shards == 0 {
//					firstBuyReport += fmt.Sprintf("建仓了%s，市值%s美元；", ticker, utils.ThousandFormatFloat64(holding.MarketValue))
//				}
//
//				if previousHolding.Shards > holding.Shards {
//
//				}
//
//			}
//		}
//	}
//
//
//
//	return report
//}

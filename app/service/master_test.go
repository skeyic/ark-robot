package service

import (
	"flag"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/utils"
	"testing"
)

func Test_MasterStart(t *testing.T) {
	flag.Set("logtostderr", "true")
	flag.Set("v", "10")
	flag.Parse()
	TheMaster.StartDownload()
}

func Test_MasterFreshInit(t *testing.T) {
	var (
		err error
	)

	utils.EnableGlogForTesting()
	err = TheMaster.FreshInit()
	if err != nil {
		glog.Errorf("failed to fresh init the master, err: %v", err)
		return
	}
}

func Test_MasterReport(t *testing.T) {
	var (
		err    error
		report = &Report{}
	)

	utils.EnableGlogForTesting()

	err = TheLibrary.LoadFromFileStore()
	if err != nil {
		glog.Errorf("failed to load library, err: %v", err)
		return
	}

	err = TheStockLibraryMaster.LoadAllStocks()
	if err != nil {
		glog.Errorf("failed to load stock library master, err: %v", err)
		return
	}

	//stockCurrentHoldings := TheStockLibraryMaster.GetStockCurrentHolding("MORGAN STANLEY GOVT INSTL 8035", "ARKF")
	//glog.V(4).Infof("HOLDINGS: %+v", stockCurrentHoldings)

	latestTradings := TheLibrary.LatestStockTradings
	for _, fund := range allARKTypes {
		tradings := latestTradings.GetFundStockTradings(fund)
		for _, trading := range tradings.SortedTradingList() {
			stockCurrentHoldings := TheStockLibraryMaster.GetStockCurrentHolding(trading.Ticker, trading.Fund)
			if report.Date == "" {
				report.Date = trading.Date.Format("2006-01-02")
			}
			report.StockReports = append(report.StockReports, &StockReport{
				Date:                  trading.Date.Format("2006-01-02"),
				StockTicker:           trading.Ticker,
				Company:               trading.Company,
				Cusip:                 trading.Cusip,
				Fund:                  trading.Fund,
				CurrentHoldingShards:  stockCurrentHoldings.Shards,
				CurrentDirection:      trading.Direction,
				FixDirection:          trading.FixedDirection,
				CurrentTradingShards:  trading.Shards,
				CurrentTradingPercent: trading.Percent,
				FundDirection:         tradings.Direction,
				FundTradingPercent:    tradings.Percent,
			},
			)
		}
	}

	//for _, r := range report.StockReports {
	//	if r.FixDirection == TradeKeep && r.StockTicker != "RPTX" {
	//		continue
	//	}
	//	glog.V(4).Infof("%s STOCK: %s, FUND: %s, CurrentHoldingShards: %f, DIRECTION: %s, FixDIRECTION: %s, SHARDS: %f, PERCENT: %f, FundDirection: %s, FundPERCENT: %f", r.Date, r.StockTicker, r.Fund, r.CurrentHoldingShards,
	//		r.CurrentDirection, r.FixDirection, r.CurrentTradingShards, r.CurrentTradingPercent, r.FundDirection, r.FundTradingPercent)
	//}

	err = report.ToExcel(false)
	if err != nil {
		glog.Errorf("Report to excel failed, err: %v", err)
		return
	}

}

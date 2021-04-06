package service

import (
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/utils"
	"sort"
	"testing"
)

func TestStockLibrary(t *testing.T) {
	utils.EnableGlogForTesting()
	err := TheStockLibraryMaster.LoadAllStocks()
	if err != nil {
		glog.Errorf("failed to load all stocks, err: %v", err)
		return
	}

	//for ticker, stock := range TheStockLibraryMaster.StockLibraries {
	//	glog.V(4).Infof("TICKER: %s, STOCK: %+v", ticker, stock)
	//}
	var (
		ticker = "SYRS"
	)

	stock := TheStockLibraryMaster.StockLibraries[ticker]
	if stock == nil {
		panic("ticker not found")
	}

	var (
		dateList timeList
	)

	for theDate := range stock.HistoryStockTradings {
		dateList = append(dateList, theDate)
	}

	sort.Sort(dateList)
	//for idx, date := range dateList {
	//	glog.V(4).Infof("IDX: %d, DATE: %s", idx, date)
	//}

	if dateList == nil {
		return
	}

	for i := 0; i < len(dateList); i++ {
		fundTradings := stock.HistoryStockTradings[dateList[i]]
		for _, fund := range allARKTypes {
			if fund != "ARKG" {
				continue
			}
			trading := fundTradings[fund]
			if trading != nil {
				glog.V(4).Infof("DATE: %s, FUND: %s, FD: %s, SHARDS: %f, PERCENT: %f", dateList[i], fund, trading.FixedDirection, trading.Shards, trading.Percent)
			}
		}
	}

}

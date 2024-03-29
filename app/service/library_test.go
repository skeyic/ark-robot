package service

import (
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/utils"
	"testing"
	"time"
)

func Test_LibraryInit(t *testing.T) {
	utils.EnableGlogForTesting()

	//glog.V(4).Infof("The LIBRARY: %+v", TheLibrary)
	err := TheLibrary.LoadFromDirectory()
	if err != nil {
		glog.Errorf("failed to load from directory, err: %v", err)
		return
	}

	pDate, _ := time.Parse("2006-01-02", "2021-02-25")
	cDate, _ := time.Parse("2006-01-02", "2021-02-26")
	//tDate, _ := time.Parse("2006-01-02", "2021-02-10")

	//for key, value := range TheLibrary.HistoryStockHoldings[date]["ARKW"].Holdings {
	//	glog.V(4).Infof("KEY: %s, VALUE: %+v", key, value)
	//}

	//for key, value := range TheLibrary.LatestStockHoldings["ARKW"].Holdings {
	//	glog.V(4).Infof("KEY: %s, VALUE: %+v", key, value)
	//}

	// ARKW	02/16/2021	Sell	SE	81141R100	SEA LTD	31,931	0.0964
	//glog.V(4).Infof("SE P: %+v", TheLibrary.HistoryStockHoldings[pDate]["ARKW"].Holdings["SE"])
	//glog.V(4).Infof("SE C: %+v", TheLibrary.HistoryStockHoldings[cDate]["ARKW"].Holdings["SE"])
	//glog.V(4).Infof("SE C: %+v", TheLibrary.HistoryStockHoldings[pDate]["ARKW"].Holdings["SE"])

	tradings := TheLibrary.HistoryStockHoldings[cDate].GetFundStockHoldings("ARKG").GenerateTrading(TheLibrary.HistoryStockHoldings[pDate].GetFundStockHoldings("ARKG"))

	//glog.V(4).Infof("TRADINGS: %+v", tradings.Tradings)
	//glog.V(4).Infof("TRADINGS SE: %+v", tradings.Tradings["SE"])

	//28	ARKW	02/16/2021	Buy	PLTR	69608A108	PALANTIR TECHNOLOGIES INC	1,560,200	0.4961
	//29	ARKW	02/16/2021	Sell	PSTG	74624M102	PURE STORAGE INC	265,044	0.0814
	//30	ARKW	02/16/2021	Sell	PD	69553P100	PAGERDUTY INC	31,012	0.0188
	//31	ARKW	02/16/2021	Sell	Z	98954M200	ZILLOW GROUP INC	67,864	0.1481
	//32	ARKW	02/16/2021	Sell	API	00851L103	AGORA INC	25,621	0.0287
	//33	ARKW	02/16/2021	Sell	SE	81141R100	SEA LTD	31,931	0.0964

	tradings.SetFixDirection()
	for idx, trading := range tradings.SortedTradingList() {
		//if trading.Ticker == "PLTR" || trading.Ticker == "PSTG" || trading.Ticker == "PD" || trading.Ticker == "Z" ||
		//	trading.Ticker == "API" || trading.Ticker == "SE" {
		//	glog.V(4).Infof("TRADING: %+v", trading)
		//	glog.V(4).Infof("TICKER: %s, DIRECTION: %s, FDIRECTION: %s, SHARDS: %f, PERCENT: %f", trading.Ticker, trading.Direction, trading.FixedDirection, trading.Shards, trading.Percent)
		//}
		glog.V(4).Infof("IDX: %d, TICKER: %s, DIRECTION: %s, FDIRECTION: %s, SHARDS: %f, PERCENT: %f", idx, trading.Ticker, trading.Direction, trading.FixedDirection, trading.Shards, trading.Percent)
	}
}

func Test_LibraryInit2(t *testing.T) {
	utils.EnableGlogForTesting()

	//glog.V(4).Infof("The LIBRARY: %+v", TheLibrary)
	err := TheLibrary.LoadFromFileStore()
	if err != nil {
		glog.Errorf("failed to load from file store, err: %v", err)
		return
	}

	pDate, _ := time.Parse("2006-01-02", "2021-02-26")
	cDate, _ := time.Parse("2006-01-02", "2021-03-01")
	//tDate, _ := time.Parse("2006-01-02", "2021-02-10")

	//for key, value := range TheLibrary.HistoryStockHoldings[date]["ARKW"].Holdings {
	//	glog.V(4).Infof("KEY: %s, VALUE: %+v", key, value)
	//}

	//for key, value := range TheLibrary.LatestStockHoldings["ARKW"].Holdings {
	//	glog.V(4).Infof("KEY: %s, VALUE: %+v", key, value)
	//}

	// ARKW	02/16/2021	Sell	SE	81141R100	SEA LTD	31,931	0.0964
	//glog.V(4).Infof("SE P: %+v", TheLibrary.HistoryStockHoldings[pDate]["ARKW"].Holdings["SE"])
	//glog.V(4).Infof("SE C: %+v", TheLibrary.HistoryStockHoldings[cDate]["ARKW"].Holdings["SE"])
	//glog.V(4).Infof("SE C: %+v", TheLibrary.HistoryStockHoldings[pDate]["ARKW"].Holdings["SE"])

	tradings := TheLibrary.HistoryStockHoldings[cDate].GetFundStockHoldings("ARKF").GenerateTrading(TheLibrary.HistoryStockHoldings[pDate].GetFundStockHoldings("ARKF"))

	//glog.V(4).Infof("TRADINGS: %+v", tradings.Tradings)
	//glog.V(4).Infof("TRADINGS SE: %+v", tradings.Tradings["SE"])

	//28	ARKW	02/16/2021	Buy	PLTR	69608A108	PALANTIR TECHNOLOGIES INC	1,560,200	0.4961
	//29	ARKW	02/16/2021	Sell	PSTG	74624M102	PURE STORAGE INC	265,044	0.0814
	//30	ARKW	02/16/2021	Sell	PD	69553P100	PAGERDUTY INC	31,012	0.0188
	//31	ARKW	02/16/2021	Sell	Z	98954M200	ZILLOW GROUP INC	67,864	0.1481
	//32	ARKW	02/16/2021	Sell	API	00851L103	AGORA INC	25,621	0.0287
	//33	ARKW	02/16/2021	Sell	SE	81141R100	SEA LTD	31,931	0.0964

	tradings.SetFixDirection()
	for idx, trading := range tradings.SortedTradingList() {
		//if trading.Ticker == "PLTR" || trading.Ticker == "PSTG" || trading.Ticker == "PD" || trading.Ticker == "Z" ||
		//	trading.Ticker == "API" || trading.Ticker == "SE" {
		//	glog.V(4).Infof("TRADING: %+v", trading)
		//	glog.V(4).Infof("TICKER: %s, DIRECTION: %s, FDIRECTION: %s, SHARDS: %f, PERCENT: %f", trading.Ticker, trading.Direction, trading.FixedDirection, trading.Shards, trading.Percent)
		//}
		glog.V(4).Infof("IDX: %d, TICKER: %s, DIRECTION: %s, FDIRECTION: %s, SHARDS: %f, PERCENT: %f", idx, trading.Ticker, trading.Direction, trading.FixedDirection, trading.Shards, trading.Percent)
	}
}

func Test_GenerateTradings(t *testing.T) {
	utils.EnableGlogForTesting()
	err := TheLibrary.LoadFromDirectory()
	if err != nil {
		glog.Errorf("failed to load from directory, err: %v", err)
		return
	}
	TheLibrary.GenerateTradings()
}

func Test_LibraryLoad(t *testing.T) {
	var (
		err error
	)

	utils.EnableGlogForTesting()
	err = TheLibrary.LoadFromFileStore()
	if err != nil {
		glog.Errorf("failed to load library from the file store, err: %v", err)
		return
	}

	//cDate, _ := time.Parse("2006-01-02", "2021-03-04")

	//for _, fund := range allARKTypes {
	//	for ticker, holding := range TheLibrary.HistoryStockHoldings[cDate].GetFundStockHoldings(fund).Holdings {
	//		glog.V(4).Infof("TICKER: %s, HOLDING: %+v", ticker, holding)
	//	}
	//}

	for date, _ := range TheLibrary.HistoryStockHoldings {
		glog.V(4).Infof("DATE: %s", date)
	}
}

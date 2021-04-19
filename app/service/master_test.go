package service

import (
	"flag"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/utils"
	"strings"
	"testing"
	"time"
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

	for date := range TheLibrary.HistoryStockHoldings {
		glog.V(4).Infof("H DATE: %s", date)
	}

	for date := range TheLibrary.HistoryStockTradings {
		glog.V(4).Infof("T DATE: %s", date)
	}
}

func Test_MasterGetPreviousHoldings(t *testing.T) {
	var (
		err error
	)

	utils.EnableGlogForTesting()
	err = TheMaster.FreshInit()
	if err != nil {
		glog.Errorf("failed to fresh init the master, err: %v", err)
		return
	}

	for date := range TheLibrary.HistoryStockHoldings {
		glog.V(4).Infof("H DATE: %s", date)
	}

	for _, holding := range TheLibrary.GetPreviousHoldings(TheLibrary.GetLatestHoldingDate(), 3) {
		glog.V(4).Infof("P DATE: %s", holding.Date)
	}
}

func Test_MasterFreshInitWithDownload(t *testing.T) {
	var (
		err error
	)

	utils.EnableGlogForTesting()
	err = TheMaster.FreshInit()
	if err != nil {
		glog.Errorf("failed to fresh init the master, err: %v", err)
		return
	}

	glog.V(4).Infof("Latest stock holding date: %s", TheLibrary.LatestStockHoldings.Date)
	glog.V(4).Infof("Latest stock trading date: %s", TheLibrary.LatestStockTradings.Date)

	err = TheDownloader.DownloadAllARKCSVs()
	if err != nil {
		glog.Errorf("failed to download csv, err: %v", err)
		return
	}

	glog.V(4).Infof("Latest stock holding date: %s", TheLibrary.LatestStockHoldings.Date)
	glog.V(4).Infof("Latest stock trading date: %s", TheLibrary.LatestStockTradings.Date)
}

func Test_MasterStaleInitWithDownload(t *testing.T) {
	var (
		err error
	)

	utils.EnableGlogForTesting()
	err = TheMaster.StaleInit()
	if err != nil {
		glog.Errorf("failed to fresh init the master, err: %v", err)
		return
	}

	glog.V(4).Infof("Latest stock holding date: %s", TheLibrary.LatestStockHoldings.Date)
	glog.V(4).Infof("Latest stock trading date: %s", TheLibrary.LatestStockTradings.Date)

	err = TheDownloader.DownloadAllARKCSVs()
	if err != nil {
		glog.Errorf("failed to download csv, err: %v", err)
		return
	}

	glog.V(4).Infof("Latest stock holding date: %s", TheLibrary.LatestStockHoldings.Date)
	glog.V(4).Infof("Latest stock trading date: %s", TheLibrary.LatestStockTradings.Date)
}

func Test_MasterReportLatest(t *testing.T) {
	var (
		err error
	)

	utils.EnableGlogForTesting()

	err = TheMaster.StaleInit()
	if err != nil {
		glog.Errorf("failed to fresh init the master, err: %v", err)
		return
	}

	err = TheMaster.ReportLatestTrading(true)
	if err != nil {
		glog.Errorf("failed to report latest trading, err: %v", err)
		return
	}
}

func Test_MasterCheckHoldings(t *testing.T) {
	var (
		err error
	)

	utils.EnableGlogForTesting()
	err = TheMaster.StaleInit()
	if err != nil {
		glog.Errorf("failed to fresh init the master, err: %v", err)
		return
	}

	//for _, fund := range allARKTypes {
	//	tradings := TheLibrary.LatestStockTradings.GetFundStockTradings(fund)
	//	glog.V(4).Infof("FUND: %s, TRADING NUM: %d", fund, len(tradings.Tradings))
	//}
	//
	for theDate := range TheLibrary.HistoryStockHoldings {
		glog.V(4).Infof("H DATE: %s", theDate)
	}

	for theDate := range TheLibrary.HistoryStockTradings {
		glog.V(4).Infof("T DATE: %s", theDate)
	}

	glog.V(4).Infof("LATEST H: %s, LATEST T: %s", TheLibrary.LatestStockHoldings.Date,
		TheLibrary.LatestStockHoldings.Date)

	//if err != nil {
	//	glog.Errorf("failed to report latest trading, err: %v", err)
	//	return
	//}
}

func Test_MasterIndexToES(t *testing.T) {
	var (
		err error
	)

	utils.EnableGlogForTesting()
	err = TheMaster.FreshInit()
	if err != nil {
		glog.Errorf("failed to fresh init the master, err: %v", err)
		return
	}

	err = TheMaster.IndexToES()
	if err != nil {
		glog.Errorf("failed to index to es, err: %v", err)
		return
	}
}

func Test_MasterCheckChinaStock(t *testing.T) {
	var (
		err error
	)

	utils.EnableGlogForTesting()

	err = TheChinaStockManager.FreshInit()
	if err != nil {
		glog.Errorf("failed to fresh init the china stock manager, err: %v", err)
		return
	}

	glog.V(4).Infof("NUM: %d", len(TheChinaStockManager.stocks))
	for _, stock := range TheChinaStockManager.stocks {
		glog.V(4).Infof("STOCK: %+v", stock)
	}
}

func Test_MasterReportStocks(t *testing.T) {
	var (
		err    error
		stocks = []string{"JD", "HUYA", "BIDU", "PDD", "BABA"}
		//stocks      = []string{"JD"}
		fromDate, _ = time.Parse(TheDateFormat, "2021-04-05")
		endDate, _  = time.Parse(TheDateFormat, "2021-04-13")
	)

	utils.EnableGlogForTesting()
	err = TheMaster.FreshInit()
	if err != nil {
		glog.Errorf("failed to fresh init the master, err: %v", err)
		return
	}

	for _, stock := range stocks {
		err = TheMaster.ReportStock(stock, fromDate, endDate)
		if err != nil {
			glog.Errorf("failed to report stock %s from %s to %s, err: %v", stock, fromDate, endDate, err)
			return
		}
	}

	glog.V(4).Info("REPORTED")
}

func Test_MasterCheckStocks(t *testing.T) {
	var (
		err error
	)

	utils.EnableGlogForTesting()
	err = TheMaster.StaleInit()
	if err != nil {
		glog.Errorf("failed to fresh init the master, err: %v", err)
		return
	}

	TheStockLibraryMaster.lock.RLock()
	for theStock := range TheStockLibraryMaster.StockLibraries {
		if strings.ContainsAny(theStock, "_") {
			glog.V(4).Infof("TICKER: %s", theStock)
		}
	}
	TheStockLibraryMaster.lock.RUnlock()
}

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

func Test_MasterReportLatest(t *testing.T) {
	var (
		err error
	)

	utils.EnableGlogForTesting()
	TheChinaStockManager.Init()

	err = TheMaster.FreshInit()
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

func Test_MasterCheckTradings(t *testing.T) {
	var (
		err error
	)

	utils.EnableGlogForTesting()
	err = TheMaster.FreshInit()
	if err != nil {
		glog.Errorf("failed to fresh init the master, err: %v", err)
		return
	}

	for _, fund := range allARKTypes {
		tradings := TheLibrary.LatestStockTradings.GetFundStockTradings(fund)
		glog.V(4).Infof("FUND: %s, TRADING NUM: %d", fund, len(tradings.Tradings))
	}
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
	//var (
	//	err error
	//)

	utils.EnableGlogForTesting()

	glog.V(4).Infof("NUM: %d", len(TheChinaStockManager.stocks))
}

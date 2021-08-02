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

	glog.V(4).Infof("LATEST DATE: %s", TheLibrary.GetLatestHoldingDate())
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

	//TheDate, _ := time.Parse(TheDateFormat, "2021-04-20")
	//err = TheMaster.Report(TheDate, true, config.Config.Report.SpecialTradingPercent)
	//if err != nil {
	//	glog.Errorf("failed to report latest trading, err: %v", err)
	//	return
	//}
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
		stocks = []string{"TSLA", "COIN", "TCEHY", "3690", "BEKE", "JD", "HUYA", "BIDU", "PDD", "BABA", "BYDDY", "NIU"}
		//stocks      = []string{"JD"}
		fromDate, _ = time.Parse(TheDateFormat, "2021-04-26")
		endDate, _  = time.Parse(TheDateFormat, "2021-04-30")
	)

	utils.EnableGlogForTesting()
	err = TheMaster.StaleInit()
	if err != nil {
		glog.Errorf("failed to stale init the master, err: %v", err)
		return
	}

	for _, stock := range stocks {
		err = TheMaster.ReportStock(stock, fromDate, endDate, "")
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

func Test_MasterReportStocks2(t *testing.T) {
	var (
		err error
		//stocks = []string{"TSLA"}
		stocks      = []string{"JD"}
		fromDate, _ = time.Parse(TheDateFormat, "2021-04-26")
		endDate, _  = time.Parse(TheDateFormat, "2021-04-30")
	)

	utils.EnableGlogForTesting()
	err = TheMaster.StaleInit()
	if err != nil {
		glog.Errorf("failed to stale init the master, err: %v", err)
		return
	}

	for _, stock := range stocks {
		//NewStockDateRangeReport(stock, fromDate, endDate).Report()
		err = TheMaster.ReportStock(stock, fromDate, endDate, "")
		if err != nil {
			glog.Errorf("failed to report stock %s from %s to %s, err: %v", stock, fromDate, endDate, err)
			return
		}
	}

	glog.V(4).Info("REPORTED")
}

func Test_MasterReportStocks3(t *testing.T) {
	var (
		err error
		//stocks = []string{"TSLA"}
		stocks = []string{"JD"}
		days   int
		//fromDate, _ = time.Parse(TheDateFormat, "2021-04-26")
		//endDate, _  = time.Parse(TheDateFormat, "2021-04-30")
	)

	utils.EnableGlogForTesting()
	err = TheMaster.StaleInit()
	if err != nil {
		glog.Errorf("failed to stale init the master, err: %v", err)
		return
	}

	for _, stock := range stocks {
		//NewStockDateRangeReport(stock, fromDate, endDate).Report()
		err = TheMaster.ReportStockByDays(stock, 10, "")
		if err != nil {
			glog.Errorf("failed to report stock %s for %d days, err: %v", stock, days, err)
			return
		}
	}

	glog.V(4).Info("REPORTED")
}

func Test_MasterReportStocksCurrent(t *testing.T) {
	var (
		err error
		//stocks = []string{"TSLA"}
		stocks = []string{"TCEHY"}
		//fromDate, _ = time.Parse(TheDateFormat, "2021-04-26")
		//endDate, _  = time.Parse(TheDateFormat, "2021-04-30")
	)

	utils.EnableGlogForTesting()
	err = TheMaster.StaleInit()
	//err = TheMaster.FreshInit()
	if err != nil {
		glog.Errorf("failed to stale init the master, err: %v", err)
		return
	}

	for _, stock := range stocks {
		//NewStockDateRangeReport(stock, fromDate, endDate).Report()
		report, err := TheMaster.ReportStockCurrent(stock)
		if err != nil {
			glog.Errorf("failed to report stock %s, err: %v", stock, err)
			return
		}
		glog.V(4).Infof("REPORT: %s", report)
		h := TheStockLibraryMaster.GetStockLatestHolding(stock)
		glog.V(4).Infof("HHH: %v", h)
		t := TheStockLibraryMaster.GetStockLatestTrading(stock)
		glog.V(4).Infof("TTT: %v", t)
	}

	glog.V(4).Info("REPORTED")
}

func Test_MasterReportTop10(t *testing.T) {
	var (
		fund = "ARKG"
		err  error
	)

	utils.EnableGlogForTesting()
	err = TheMaster.StaleInit()
	if err != nil {
		glog.Errorf("failed to stale init the master, err: %v", err)
		return
	}

	report, err := TheMaster.ReportFundTop10(fund)
	if err != nil {
		glog.Errorf("failed to report fund %s top10, err: %v", fund, err)
		return
	}
	glog.V(4).Infof("REPORT: \n%s", report)

	glog.V(4).Info("REPORTED")
}

func Test_MasterNewStockLibraryMaster(t *testing.T) {
	var (
		err error
		//stocks = []string{"TSLA"}
		stocks = []string{"HUYA"}
		//fromDate, _ = time.Parse(TheDateFormat, "2021-04-26")
		//endDate, _  = time.Parse(TheDateFormat, "2021-04-30")
	)

	utils.EnableGlogForTesting()
	err = TheMaster.StaleInit()
	if err != nil {
		glog.Errorf("failed to stale init the master, err: %v", err)
		return
	}

	for _, stock := range stocks {
		//NewStockDateRangeReport(stock, fromDate, endDate).Report()
		currentHolding := TheStockLibraryMaster.GetStockLatestHolding(stock)
		if currentHolding != nil {
			for _, fund := range allARKTypes {
				holding := currentHolding.GetFundHolding(fund)
				if holding != nil {
					glog.V(4).Infof("FUND: %s, Holding %v", fund, holding)
				} else {
					glog.V(4).Infof("FUND: %s, no holding", fund)
				}
			}
		} else {
			glog.V(4).Infof("currentHolding is empty")
		}

		currentTrading := TheStockLibraryMaster.GetStockLatestTrading(stock)
		if currentTrading != nil {
			for _, fund := range allARKTypes {
				holding := currentTrading.GetFundTrading(fund)
				if holding != nil {
					glog.V(4).Infof("FUND: %s, Trading %v", fund, holding)
				} else {
					glog.V(4).Infof("FUND: %s, no trading", fund)
				}
			}
		} else {
			glog.V(4).Infof("currentTrading is empty")
		}

		//theLibrary := TheStockLibraryMaster.GetStockLibrary(stock)
		//glog.V(4).Infof("H: %v", theLibrary.LatestStockHolding)
		//glog.V(4).Infof("T: %v", theLibrary.LatestStockTrading)
		//for _, holding := range theLibrary.HistoryStockHoldings {
		//	glog.V(4).Infof("HH: %v", holding)
		//}
	}
}

func Test_MasterNewChinaStockReport(t *testing.T) {
	var (
		err error
		//stocks = []string{"TSLA"}
		//fromDate, _ = time.Parse(TheDateFormat, "2021-04-26")
		//endDate, _  = time.Parse(TheDateFormat, "2021-04-30")
	)

	utils.EnableGlogForTesting()
	err = TheMaster.StaleInit()
	if err != nil {
		glog.Errorf("failed to stale init the master, err: %v", err)
		return
	}

	report := NewChinaStockHoldingReport()
	err = report.Load()
	if err != nil {
		glog.Errorf("Load failed, err: %v", err)
		return
	}

	glog.V(4).Infof("REPORTS: %s", report.TxtReport())
}

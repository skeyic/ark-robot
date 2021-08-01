package service

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/config"
	"sync"
	"time"
)

var (
	masterDownloadTimeInterval = 2 * time.Hour
)

var (
	TheMaster = NewMaster()
)

type Master struct {
	lock *sync.RWMutex
}

func NewMaster() *Master {
	m := &Master{
		lock: &sync.RWMutex{},
	}
	m.MustInit()
	return m
}

func (m *Master) MustInit() {
	var (
		err error
	)

	err = TheDownloader.Init()
	if err != nil {
		glog.Errorf("init the downloader failed")
		panic(err)
	}

	err = ThePorter.Init()
	if err != nil {
		glog.Errorf("init the porter failed")
		panic(err)
	}

	err = TheChinaStockManager.Init()
	if err != nil {
		glog.Errorf("init the china stock manager failed")
		panic(err)
	}

	err = TheLibrary.Init()
	if err != nil {
		glog.Errorf("init the library failed")
		panic(err)
	}

	err = TheStockLibraryMaster.Init()
	if err != nil {
		glog.Errorf("init the stock library master failed")
		panic(err)
	}
}

// StaleInit Init TheLibrary, TheStockLibraryMaster from the stored holdings in directory
func (m *Master) StaleInit() error {
	var (
		err error
	)

	err = TheLibrary.StaleInit()
	if err != nil {
		glog.Errorf("failed to stale init the library, err: %v", err)
		return err
	}

	err = TheStockLibraryMaster.StaleInit()
	if err != nil {
		glog.Errorf("failed to stale init the stock library master, err: %v", err)
		return err
	}

	TheTop10HoldingsReportMaster.Refresh()
	return nil
}

// FreshInit Init the holding, trading and stock library from the stored holdings in directory
func (m *Master) FreshInit() error {
	var (
		err error
	)

	err = TheChinaStockManager.FreshInit()
	if err != nil {
		glog.Errorf("failed to init china stock, err: %v", err)
		return err
	}

	err = ThePorter.LoadFromDirectory()
	if err != nil {
		glog.Errorf("failed to load holdings from directory, err: %v", err)
		return err
	}

	TheLibrary.GenerateTradings()
	TheTop10HoldingsReportMaster.Refresh()

	m.MustSave()
	return nil
}

func (m *Master) MustSave() {
	TheLibrary.MustSave()
	TheStockLibraryMaster.MustSave()
}

func (m *Master) ReportLatestTrading(full bool) error {
	latestDate := TheLibrary.GetLatestHoldingDate()
	if latestDate.IsZero() {
		return errNoLatestDate
	}

	return m.Report(latestDate, full, config.Config.Report.SpecialTradingPercent)
}

func (m *Master) Report(date time.Time, full bool, SpecialTradingsPercent float64) error {
	var (
		err error
	)

	tradingsReport := NewTradingsReport(date)
	err = tradingsReport.Report(full)
	if err != nil {
		glog.Errorf("tradingsReport to excel failed, err: %v", err)
		return err
	}

	top10HoldingsReport := NewTop10HoldingsReport(date, allARKTypes)
	err = top10HoldingsReport.Report()
	if err != nil {
		glog.Errorf("top10HoldingsReport to excel failed, err: %v", err)
		return err
	}

	specialTradingsReport := NewSpecialTradingsReport(date, SpecialTradingsPercent)
	err = specialTradingsReport.Report()
	if err != nil {
		glog.Errorf("specialTradingsReport to excel failed, err: %v", err)
		return err
	}

	chinaStockTradingsReport := NewChinaStockTradingsReport(date)
	err = chinaStockTradingsReport.Report()
	if err != nil {
		glog.Errorf("chinaStockTradingsReport to excel failed, err: %v", err)
		return err
	}

	return nil
}

func (m *Master) ReportStock(ticker string, fromDate, endDate time.Time, funds string) error {
	var (
		err error
	)

	err = NewStockDateRangeReport(ticker, fromDate, endDate, funds).Report()
	if err != nil {
		glog.Errorf("report stock %s from %s to %s failed, err: %v", ticker, fromDate.Format(TheDateFormat),
			endDate.Format(TheDateFormat), err)
		return err
	}

	return nil
}

func (m *Master) ReportStockByDays(ticker string, days int64, funds string) error {
	var (
		err error
	)

	err = NewStockDateRangeReportFromDays(ticker, days, funds).Report()
	if err != nil {
		glog.Errorf("report stock %s for %d days failed, err: %v", ticker, days, err)
		return err
	}

	return nil
}

func (m *Master) ReportStockCurrent(ticker string) (report string, err error) {
	r := NewStockCurrentReport(ticker)

	err = r.Load()
	if err != nil {
		glog.Warningf("TICKER: %s, ERR: %v", ticker, err)
		if err == errStockNotHold {
			var (
				report = fmt.Sprintf("ARK当前未持有%s", ticker)
			)
			h := TheStockLibraryMaster.GetStockLatestTrading(ticker)
			if h != nil {
				report += fmt.Sprintf("，ARK曾经持有，但于%s清完全部的持仓。", h.Date.Format(TheDateFormat))
			}
			return report, nil
		}
		glog.Errorf("report stock %s failed, err: %v", ticker, err)
		return
	}

	report = r.TxtReport()

	return
}

func (m *Master) ReportFundTop10(fund string) (report string, err error) {
	report = TheTop10HoldingsReportMaster.GetFundTop10(fund)
	return
}

func (m *Master) ReportContinue3Days() (report string, err error) {
	report = TheContinue3DaysReportMaster.GetReport()
	return
}

func (m *Master) ReportBigSwings() (report string, err error) {
	report = TheBigSwingsReportMaster.GetReport()
	return
}

func (m *Master) ReportChinaStock() (report string, err error) {
	report = TheChinaStockReportMaster.GetReport()
	return
}

func (m *Master) IsTicker(ticker string) bool {
	return TheStockLibraryMaster.IsTicker(ticker)
}

func (m *Master) GetAllTickers() []string {
	return TheStockLibraryMaster.GetAllTickers()
}

func (m *Master) IndexLatestToES() (err error) {
	latestDate := TheLibrary.GetLatestHoldingDate()
	if latestDate.IsZero() {
		return errNoLatestDate
	}

	return m.IndexDateToES(latestDate)
}

func (m *Master) IndexDateToES(date time.Time) (err error) {
	latestHoldings := TheLibrary.GetHoldings(date)
	for _, fund := range allARKTypes {
		err = TheESConnector.IndexStockHoldings(latestHoldings.GetFundStockHoldings(fund))
		if err != nil {
			glog.Errorf("Index date %s holdings to ES failed, err: %v", date, err)
			return
		}
	}

	latestTradings := TheLibrary.GetTradings(date)
	for _, fund := range allARKTypes {
		err = TheESConnector.IndexStockTradings(latestTradings.GetFundStockTradings(fund))
		if err != nil {
			glog.Errorf("Index date %s tradings to ES failed, err: %v", date, err)
			return
		}
	}
	return nil
}

func (m *Master) IndexToES() error {
	var (
		err error
	)

	for _, arkHoldings := range TheLibrary.HistoryStockHoldings {
		for _, fund := range allARKTypes {
			err = TheESConnector.IndexStockHoldings(arkHoldings.GetFundStockHoldings(fund))
			if err != nil {
				glog.Errorf("failed to index holdings %s %s, err: %v", arkHoldings.Date, fund, err)
				return err
			}
		}
	}

	for _, arkTradings := range TheLibrary.HistoryStockTradings {
		for _, fund := range allARKTypes {
			err = TheESConnector.IndexStockTradings(arkTradings.GetFundStockTradings(fund))
			if err != nil {
				glog.Errorf("failed to index holdings %s %s, err: %v", arkTradings.Date, fund, err)
				return err
			}
		}
	}

	return nil
}

func (m *Master) StartDownload() {
	var (
		ticker   = time.NewTicker(masterDownloadTimeInterval)
		initChan = make(chan time.Time, 1)
	)

	go func() {
		time.Sleep(30 * time.Second)
		glog.V(4).Infof("kick off master at the beginning")
		initChan <- time.Now().UTC()
	}()

	processFunc := func(a time.Time) (bool, error) {
		var (
			err error
		)

		err = TheDownloader.DownloadAllARKCSVs()
		if err != nil {
			glog.Errorf("download All ARK CSVs failed, wait and retry, current time: %s, err: %v", a, err)
			return false, err
		}

		return true, nil
	}

	for {
		select {
		case a := <-ticker.C:
			a = a.UTC()
			result, err := processFunc(a)
			if err != nil {
				glog.Errorf("failed to process, current time: %s, err: %v", a, err)
			}
			if result {
				glog.V(4).Infof("process done, current time: %s", a)
			}
		case a := <-initChan:
			result, err := processFunc(a)
			if err != nil {
				glog.Errorf("failed to process, current time: %s, err: %v", a, err)
			}
			if result {
				glog.V(4).Infof("process done, current time: %s", a)
			}
		}
	}
}

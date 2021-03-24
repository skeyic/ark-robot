package service

import (
	"github.com/golang/glog"
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
}

// Init the holding, trading and stock library from the stored holdings in directory
func (m *Master) FreshInit() error {
	var (
		err error
	)

	err = ThePorter.LoadFromDirectory()
	if err != nil {
		glog.Errorf("failed to load holdings from directory, err: %v", err)
		return err
	}

	TheLibrary.GenerateTradings()

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

	return m.Report(latestDate, full)
}

func (m *Master) Report(date time.Time, full bool) error {
	var (
		err error
	)

	tradingsReport := NewTradingsReport(date)
	err = tradingsReport.ToExcel(full)
	if err != nil {
		glog.Errorf("tradingsReport to excel failed, err: %v", err)
		return err
	}

	top10HoldingsReport := NewTop10HoldingsReport(date)
	err = top10HoldingsReport.ToExcel()
	if err != nil {
		glog.Errorf("top10HoldingsReport to excel failed, err: %v", err)
		return err
	}

	specialTradingsReport := NewSpecialTradingsReport(date)
	err = specialTradingsReport.ToExcel()
	if err != nil {
		glog.Errorf("specialTradingsReport to excel failed, err: %v", err)
		return err
	}

	chinaStockTradingsReport := NewChinaStockTradingsReport(date)
	err = chinaStockTradingsReport.ToExcel()
	if err != nil {
		glog.Errorf("chinaStockTradingsReport to excel failed, err: %v", err)
		return err
	}

	return nil
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

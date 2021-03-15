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
	return &Master{
		lock: &sync.RWMutex{},
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
	var (
		err error
	)

	latestDate := TheLibrary.GetLatestHoldingDate()
	if latestDate.IsZero() {
		return errNoLatestDate
	}

	tradingsReport := NewTradingsReport(latestDate)
	err = tradingsReport.ToExcel(full)
	if err != nil {
		glog.Errorf("tradingsReport to excel failed, err: %v", err)
		return err
	}

	top10HoldingsReport := NewTop10HoldingsReport(latestDate)
	err = top10HoldingsReport.ToExcel()
	if err != nil {
		glog.Errorf("top10HoldingsReport to excel failed, err: %v", err)
		return err
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

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
	lock       *sync.RWMutex
	latestDate time.Time
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

	return nil
}

func (m *Master) ReportLatestTrading() *Report {
	var (
		report = &Report{}
	)

	return report
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

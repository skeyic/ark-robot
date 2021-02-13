package service

import (
	"github.com/golang/glog"
	"sync"
	"time"
)

var (
	masterProcessTimeInterval = 5 * time.Minute
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

func (m *Master) SetLatestDate(a time.Time) {
	m.lock.Lock()
	m.latestDate = a.UTC()
	m.lock.Unlock()
}

func (m *Master) GetLatestDate() time.Time {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.latestDate
}

func (m *Master) Start() {
	var (
		ticker   = time.NewTicker(masterProcessTimeInterval)
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

		latestDay := m.GetLatestDate().Day()
		if a.Day() == latestDay {
			glog.V(10).Infof("we have already processed today, current time: %s, current day: %d", a, latestDay)
			return true, nil
		}

		if a.Hour() != downloaderUTCStartHour {
			glog.V(10).Infof("not the correct time to process, skip, current time: %s, current hour: %d, expect hour: %d", a, a.Hour(), downloaderUTCStartHour)
			return false, nil
		}

		previousDate := TheLibrary.GetLatestHoldingDate()

		err = TheDownloader.DownloadAllARKCSVs()
		if err != nil {
			glog.Errorf("download All ARK CSVs failed, wait and retry, current time: %s, err: %v", a, err)
			return false, err
		}

		latestDate := TheLibrary.GetLatestHoldingDate()
		if previousDate == latestDate {
			glog.V(4).Infof("no need to update this round, current time: %s, previous: %s", a, latestDate)
			return true, nil
		}

		glog.V(4).Infof("update latest holding in library from %s to %s, current time: %s", previousDate, latestDate, a)
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
				m.SetLatestDate(a)
			}
		case a := <-initChan:
			result, err := processFunc(a)
			if err != nil {
				glog.Errorf("failed to process, current time: %s, err: %v", a, err)
			}
			if result {
				glog.V(4).Infof("process done, current time: %s", a)
				m.SetLatestDate(a)
			}
		}
	}
}

package service

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/config"
	"github.com/skeyic/ark-robot/utils"
	"os"
	"sync"
	"time"
)

var (
	libraryFolder       = config.Config.DataFolder + "/library/"
	theLibraryFileStore = utils.NewFileStoreSvc(libraryFolder + "TheLibrary")
)

var (
	TheLibrary = NewLibrary()
)

type Library struct {
	lock                 *sync.RWMutex
	CurrentStockHoldings map[string]*StockHoldings
	LatestStockTradings  map[string]*StockTradings
	HistoryStockHoldings map[time.Time]map[string]*StockHoldings
	HistoryStockTradings map[time.Time]map[string]*StockTradings
}

func NewLibrary() *Library {
	r := &Library{
		lock:                 &sync.RWMutex{},
		CurrentStockHoldings: make(map[string]*StockHoldings),
		LatestStockTradings:  make(map[string]*StockTradings),
		HistoryStockHoldings: make(map[time.Time]map[string]*StockHoldings),
		HistoryStockTradings: make(map[time.Time]map[string]*StockTradings),
	}
	r.init()
	return r
}

func (r *Library) GetLatestHoldingDate() time.Time {
	var (
		latestTime = time.Time{}
	)
	r.lock.RLock()
	defer r.lock.RUnlock()

	if r.CurrentStockHoldings != nil {
		for fund, holdings := range r.CurrentStockHoldings {
			glog.V(10).Infof("Fund %s, latest date: %s", fund, holdings.Date)
			if latestTime.IsZero() || latestTime.After(holdings.Date) {
				latestTime = holdings.Date
			}
		}
	}

	return latestTime
}

func (r *Library) init() {
	utils.CheckFolder(libraryFolder)
	err := r.LoadFromFileStore()
	if err != nil {
		panic(fmt.Sprintf("failed to load library from the saved file, err: %v", err))
	}
	glog.V(4).Infof("library init completed")
}

func (r *Library) LoadFromFileStore() error {
	theBytes, err := theLibraryFileStore.Read()
	if err != nil {
		if os.IsNotExist(err) {
			glog.V(4).Info("No saved file for library")
			return nil
		}
		glog.Errorf("failed to load library from the saved file")
		return err
	}

	err = json.Unmarshal(theBytes, &r)
	if err != nil {
		glog.Errorf("failed to unmarshal the saved file to library")
		return err
	}

	glog.V(10).Infof("library after load: %+v", r)
	return nil
}

func (r *Library) LoadFromDirectory() error {

	return nil
}

func (r *Library) Save() error {
	uByte, _ := json.Marshal(r)
	err := theLibraryFileStore.Save(uByte)
	if err != nil {
		glog.Errorf("failed to save library, err: %v", err)
		return err
	}
	return nil
}

func (r *Library) MustSave() {
	err := r.Save()
	if err != nil {
		panic(err)
	}
}

func (r *Library) AddStockHoldings(s *StockHoldings) {
	r.lock.Lock()
	if r.HistoryStockHoldings[s.Date] == nil {
		r.HistoryStockHoldings[s.Date] = make(map[string]*StockHoldings)
	}
	r.HistoryStockHoldings[s.Date][s.Fund] = s
	if r.CurrentStockHoldings[s.Fund] == nil || r.CurrentStockHoldings[s.Fund].Date.Before(s.Date) {
		r.CurrentStockHoldings[s.Fund] = s
	}
	r.lock.Unlock()
	r.MustSave()
}

func (r *Library) AddStockTradings(s *StockTradings) {
	r.lock.Lock()
	if r.HistoryStockTradings[s.Date] == nil {
		r.HistoryStockTradings[s.Date] = make(map[string]*StockTradings)
	}
	r.HistoryStockTradings[s.Date][s.Fund] = s
	if r.LatestStockTradings[s.Fund] == nil || r.LatestStockTradings[s.Fund].Date.Before(s.Date) {
		r.LatestStockTradings[s.Fund] = s
	}
	r.lock.Unlock()
	r.MustSave()
}

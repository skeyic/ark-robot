package service

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/config"
	"github.com/skeyic/ark-robot/utils"
	"os"
	"sort"
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
	LatestStockHoldings  map[string]*StockHoldings
	LatestStockTradings  map[string]*StockTradings
	HistoryStockHoldings map[time.Time]map[string]*StockHoldings
	HistoryStockTradings map[time.Time]map[string]*StockTradings
}

func NewLibrary() *Library {
	r := &Library{
		lock:                 &sync.RWMutex{},
		LatestStockHoldings:  make(map[string]*StockHoldings),
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

	if r.LatestStockHoldings != nil {
		for fund, holdings := range r.LatestStockHoldings {
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
	//err := r.LoadFromDirectory()
	//if err != nil {
	//	panic(fmt.Sprintf("failed to load library from the saved csv file, err: %v", err))
	//}
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

func (r *Library) LoadFromDirectory() (err error) {
	// Load all holdings
	files, err := ThePorter.ListAllCSVs()
	if err != nil {
		glog.Errorf("failed to list all csv files, err: %v", err)
		return
	}

	for _, theFile := range files {
		glog.V(10).Infof("File: %s", theFile)
		err = ThePorter.ReadCSV(theFile)
		if err != nil {
			glog.Errorf("failed to read csv file %s, err: %v", theFile, err)
			return
		}
	}

	// Generate all tradings

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
	if r.LatestStockHoldings[s.Fund] == nil || r.LatestStockHoldings[s.Fund].Date.Before(s.Date) {
		r.LatestStockHoldings[s.Fund] = s
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

func (r *Library) AddStockTradingsWithoutLock(s *StockTradings) {
	if r.HistoryStockTradings[s.Date] == nil {
		r.HistoryStockTradings[s.Date] = make(map[string]*StockTradings)
	}
	r.HistoryStockTradings[s.Date][s.Fund] = s
	if r.LatestStockTradings[s.Fund] == nil || r.LatestStockTradings[s.Fund].Date.Before(s.Date) {
		r.LatestStockTradings[s.Fund] = s
	}
	r.MustSave()
}

type timeList []time.Time

func (s timeList) Len() int {
	return len(s)
}

func (s timeList) Less(i, j int) bool {
	return s[i].Before(s[j])
}

func (s timeList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (r *Library) GenerateTradings() {
	r.lock.Lock()
	defer r.lock.Unlock()

	var (
		dateList timeList
	)

	for theDate := range r.HistoryStockHoldings {
		dateList = append(dateList, theDate)
	}

	sort.Sort(dateList)
	//for idx, date := range dateList {
	//	glog.V(4).Infof("IDX: %d, DATE: %s", idx, date)
	//}

	if dateList == nil {
		return
	}

	for i := 1; i < len(dateList); i++ {
		for _, theFund := range allARKTypes {
			tradings := TheLibrary.HistoryStockHoldings[dateList[i]][theFund].GenerateTrading(TheLibrary.HistoryStockHoldings[dateList[i-1]][theFund])
			tradings.SetFixDirection()
			r.AddStockTradingsWithoutLock(tradings)
		}
	}
}

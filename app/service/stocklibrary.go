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

// data
//   - stock
//     - HUYA
//       - 20210218HUYA

var (
	stockLibraryFolder       = config.Config.DataFolder + "/stock/"
	theStockLibraryFileStore = utils.NewFileStoreSvc(libraryFolder + "TheStockLibrary")
)

var (
	TheStockLibraryMaster = NewStockLibrary()
)

type StockLibraryMaster struct {
	lock           *sync.RWMutex
	StockLibraries map[string]*StockLibrary
}

func NewStockLibraryMaster() *StockLibraryMaster {
	r := &StockLibraryMaster{
		lock:           &sync.RWMutex{},
		StockLibraries: make(map[string]*StockLibrary),
	}
	r.init()
	return r
}

func (r *StockLibraryMaster) init() {
	utils.CheckFolder(stockLibraryFolder)

	glog.V(4).Infof("StockLibraryMaster init completed")
}

type StockLibrary struct {
	lock                 *sync.RWMutex
	Ticker               string
	LatestStockHolding   map[string]*StockHolding
	LatestStockTrading   map[string]*StockTrading
	HistoryStockHoldings map[time.Time]map[string]*StockHolding
	HistoryStockTradings map[time.Time]map[string]*StockTrading
}

func NewStockLibrary() *StockLibrary {
	r := &StockLibrary{
		lock:                 &sync.RWMutex{},
		LatestStockHolding:   make(map[string]*StockHolding),
		LatestStockTrading:   make(map[string]*StockTrading),
		HistoryStockHoldings: make(map[time.Time]map[string]*StockHolding),
		HistoryStockTradings: make(map[time.Time]map[string]*StockTrading),
	}
	r.init()
	return r
}

func (r *StockLibrary) init() {
	utils.CheckFolder(stockLibraryFolder)
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

func (r *StockLibrary) LoadFromFileStore() error {
	theBytes, err := theStockLibraryFileStore.Read()
	if err != nil {
		if os.IsNotExist(err) {
			glog.V(4).Info("No saved file for stock library")
			return nil
		}
		glog.Errorf("failed to load stock library from the saved file")
		return err
	}

	err = json.Unmarshal(theBytes, &r)
	if err != nil {
		glog.Errorf("failed to unmarshal the saved file to stock library")
		return err
	}

	glog.V(10).Infof("stock library after load: %+v", r)
	return nil
}

func (r *StockLibrary) Save() error {
	uByte, _ := json.Marshal(r)
	err := theStockLibraryFileStore.Save(uByte)
	if err != nil {
		glog.Errorf("failed to save stock library, err: %v", err)
		return err
	}
	return nil
}

func (r *StockLibrary) MustSave() {
	err := r.Save()
	if err != nil {
		panic(err)
	}
}

func (r *StockLibrary) AddStockHolding(s *StockHolding) {
	r.lock.Lock()
	if r.HistoryStockHoldings[s.Date] == nil {
		r.HistoryStockHoldings[s.Date] = make(map[string]*StockHolding)
	}
	r.HistoryStockHoldings[s.Date][s.Fund] = s
	if r.LatestStockHolding[s.Fund] == nil || r.LatestStockHolding[s.Fund].Date.Before(s.Date) {
		r.LatestStockHolding[s.Fund] = s
	}
	r.lock.Unlock()
	r.MustSave()
}

func (r *StockLibrary) AddStockTrading(s *StockTrading) {
	r.lock.Lock()
	if r.HistoryStockTradings[s.Date] == nil {
		r.HistoryStockTradings[s.Date] = make(map[string]*StockTrading)
	}
	r.HistoryStockTradings[s.Date][s.Fund] = s
	if r.LatestStockTrading[s.Fund] == nil || r.LatestStockTrading[s.Fund].Date.Before(s.Date) {
		r.LatestStockTrading[s.Fund] = s
	}
	r.lock.Unlock()
	r.MustSave()
}

func (r *StockLibrary) AddStockTradingWithoutLock(s *StockTrading) {
	if r.HistoryStockTradings[s.Date] == nil {
		r.HistoryStockTradings[s.Date] = make(map[string]*StockTrading)
	}
	r.HistoryStockTradings[s.Date][s.Fund] = s
	if r.LatestStockTrading[s.Fund] == nil || r.LatestStockTrading[s.Fund].Date.Before(s.Date) {
		r.LatestStockTrading[s.Fund] = s
	}
	r.MustSave()
}

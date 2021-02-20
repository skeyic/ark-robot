package service

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/config"
	"github.com/skeyic/ark-robot/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// data
//   - stock
//     - HUYA
//     - PLTR

var (
	stockLibraryFolder = config.Config.DataFolder + "/stock/"
)

var (
	TheStockLibraryMaster = NewStockLibraryMaster(stockLibraryFolder)
)

type StockLibraryMaster struct {
	lock           *sync.RWMutex
	fileStore      *utils.MultiFileStoreSvc
	fileStorePath  string
	StockLibraries map[string]*StockLibrary
}

func NewStockLibraryMaster(filePath string) *StockLibraryMaster {
	r := &StockLibraryMaster{
		lock:           &sync.RWMutex{},
		fileStore:      utils.NewMultiFileStoreSvc(filePath, ""),
		fileStorePath:  filePath,
		StockLibraries: make(map[string]*StockLibrary),
	}
	r.init()
	return r
}

func (r *StockLibraryMaster) init() {
	utils.CheckFolder(stockLibraryFolder)

	glog.V(4).Infof("StockLibraryMaster init completed")
}

func (r *StockLibraryMaster) ListAllStocks() (files []string, err error) {
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	}

	err = filepath.Walk(r.fileStorePath, walkFunc)
	if err != nil {
		return
	}

	return
}

func (r *StockLibraryMaster) LoadAllStocks() error {
	paths, err := r.ListAllStocks()
	if err != nil {
		return err
	}

	for _, path := range paths {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		stock := NewStockLibraryFromBytes(content)
		r.StockLibraries[stock.Ticker] = stock
	}

	return nil
}

func (r *StockLibraryMaster) AddStockHoldings(holdings *StockHoldings) {
	r.lock.Lock()
	defer r.lock.Unlock()
	for ticker, holding := range holdings.Holdings {
		stockLibrary := r.StockLibraries[ticker]
		if stockLibrary == nil {
			stockLibrary = NewStockLibrary(ticker)
			r.StockLibraries[ticker] = stockLibrary
		}
		stockLibrary.AddStockHolding(holding)
	}
}

func (r *StockLibraryMaster) AddStockTradings(tradings *StockTradings) {
	r.lock.Lock()
	defer r.lock.Unlock()
	for ticker, trading := range tradings.Tradings {
		stockLibrary := r.StockLibraries[ticker]
		if stockLibrary == nil {
			stockLibrary = NewStockLibrary(ticker)
			r.StockLibraries[ticker] = stockLibrary
		}
		stockLibrary.AddStockTrading(trading)
	}
}

type StockLibrary struct {
	lock                 *sync.RWMutex
	Ticker               string
	fileStore            *utils.FileStoreSvc
	LatestStockHolding   map[string]*StockHolding
	LatestStockTrading   map[string]*StockTrading
	HistoryStockHoldings map[time.Time]map[string]*StockHolding
	HistoryStockTradings map[time.Time]map[string]*StockTrading
}

func NewStockLibrary(ticker string) *StockLibrary {
	r := &StockLibrary{
		Ticker:               ticker,
		fileStore:            utils.NewFileStoreSvc(stockLibraryFolder + strings.TrimSpace(ticker)),
		lock:                 &sync.RWMutex{},
		LatestStockHolding:   make(map[string]*StockHolding),
		LatestStockTrading:   make(map[string]*StockTrading),
		HistoryStockHoldings: make(map[time.Time]map[string]*StockHolding),
		HistoryStockTradings: make(map[time.Time]map[string]*StockTrading),
	}
	r.init()
	return r
}

func NewStockLibraryFromBytes(theBytes []byte) *StockLibrary {
	r := &StockLibrary{}

	err := json.Unmarshal(theBytes, &r)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal the saved file to stock library, err: %v", err))
	}

	r.lock = &sync.RWMutex{}
	return r
}

func (r *StockLibrary) init() {
	//utils.CheckFolder(stockLibraryFolder + r.Ticker)
	//err := r.LoadFromFileStore()
	//if err != nil {
	//	panic(fmt.Sprintf("failed to load library from the saved file, err: %v", err))
	//}
	//err := r.LoadFromDirectory()
	//if err != nil {
	//	panic(fmt.Sprintf("failed to load library from the saved csv file, err: %v", err))
	//}
	glog.V(10).Infof("stock library %s init completed", r.Ticker)
}

func (r *StockLibrary) LoadFromFileStore() error {
	theBytes, err := r.fileStore.Read()
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
	err := r.fileStore.Save(uByte)
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

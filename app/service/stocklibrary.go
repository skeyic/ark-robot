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
	"sort"
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

func (r *StockLibraryMaster) MustSave() {
	r.lock.RLock()
	defer r.lock.RUnlock()
	for _, library := range r.StockLibraries {
		library.MustSave()
	}
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

func (r *StockLibraryMaster) AddStockHoldings(arkHoldings *ARKHoldings) {
	r.lock.Lock()
	defer r.lock.Unlock()
	for _, fund := range allARKTypes {
		for ticker, holding := range arkHoldings.GetFundStockHoldings(fund).Holdings {
			stockLibrary := r.StockLibraries[ticker]
			if stockLibrary == nil {
				stockLibrary = NewStockLibrary(ticker)
				r.StockLibraries[ticker] = stockLibrary
			}
			stockLibrary.AddStockHolding(holding)
		}
	}
}

func (r *StockLibraryMaster) AddStockTradings(arkTradings *ARKTradings) {
	r.lock.Lock()
	defer r.lock.Unlock()
	for _, fund := range allARKTypes {
		for ticker, trading := range arkTradings.GetFundStockTradings(fund).Tradings {
			stockLibrary := r.StockLibraries[ticker]
			if stockLibrary == nil {
				stockLibrary = NewStockLibrary(ticker)
				r.StockLibraries[ticker] = stockLibrary
			}
			stockLibrary.AddStockTrading(trading)
		}
	}
}

func (r *StockLibraryMaster) GetStockCurrentHolding(ticker, fund string) *StockHolding {
	r.lock.RLock()
	defer r.lock.RUnlock()

	stockLibrary := r.StockLibraries[ticker]
	if stockLibrary != nil {
		if stockLibrary.LatestStockHolding != nil {
			return stockLibrary.LatestStockHolding[fund]
		}
	}

	return nil
}

func (r *StockLibraryMaster) GetStockPreviousHoldings(ticker, fund string, days int) []*StockHolding {
	r.lock.RLock()
	defer r.lock.RUnlock()

	var (
		latestDate  time.Time
		holdingList []*StockHolding
	)

	emptyList := func(holdings []*StockHolding) []*StockHolding {
		for i := 0; i < days; i++ {
			holdings = append(holdings, nil)
		}
		return holdings
	}
	stockLibrary := r.StockLibraries[ticker]
	if stockLibrary != nil {
		if stockLibrary.LatestStockHolding != nil {
			latestHolding := stockLibrary.LatestStockHolding[fund]
			if latestHolding != nil {
				latestDate = latestHolding.Date
			}
		}
	} else {
		return emptyList(holdingList)
	}

	var (
		dateList timeList
	)

	for theDate := range stockLibrary.HistoryStockHoldings {
		if stockLibrary.HistoryStockHoldings[theDate][fund] != nil {
			dateList = append(dateList, theDate)
		}
	}

	if dateList == nil {
		return emptyList(holdingList)
	}

	sort.Sort(dateList)
	for _, date := range dateList {
		if latestDate.IsZero() || date.Before(latestDate) {
			holdingList = append(holdingList, stockLibrary.HistoryStockHoldings[date][fund])
			days--
		}
		if days == 0 {
			break
		}
	}

	for days > 0 {
		holdingList = append([]*StockHolding{nil}, holdingList...)
		days--
	}

	return holdingList
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

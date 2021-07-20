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
	return r
}

func (r *StockLibraryMaster) Init() error {
	utils.CheckFolder(stockLibraryFolder)
	glog.V(4).Infof("StockLibraryMaster init completed")
	return nil
}

func (r *StockLibraryMaster) StaleInit() error {
	return r.LoadAllStocks()
}

func (r *StockLibraryMaster) MustSave() {
	r.lock.RLock()
	defer r.lock.RUnlock()
	for _, library := range r.StockLibraries {
		library.MustSave()
	}
}

func (r *StockLibraryMaster) IsTicker(ticker string) bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	_, ok := r.StockLibraries[ticker]
	return ok
}

func (r *StockLibraryMaster) GetAllTickers() (p []string) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	for ticker := range r.StockLibraries {
		if !toSkipTicker(ticker) {
			p = append(p, ticker)
		}
	}
	return
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
		theHoldings := arkHoldings.GetFundStockHoldings(fund)
		if theHoldings != nil {
			for ticker, holding := range theHoldings.Holdings {
				stockLibrary := r.StockLibraries[ticker]
				if stockLibrary == nil {
					stockLibrary = NewStockLibrary(ticker)
					r.StockLibraries[ticker] = stockLibrary
				}
				stockLibrary.AddStockHolding(holding)
			}
		}
	}
}

func (r *StockLibraryMaster) AddStockTradings(arkTradings *ARKTradings) {
	r.lock.Lock()
	defer r.lock.Unlock()
	for _, fund := range allARKTypes {
		theTradings := arkTradings.GetFundStockTradings(fund)
		if theTradings != nil {
			for ticker, trading := range theTradings.Tradings {
				stockLibrary := r.StockLibraries[ticker]
				if stockLibrary == nil {
					stockLibrary = NewStockLibrary(ticker)
					r.StockLibraries[ticker] = stockLibrary
				}
				stockLibrary.AddStockTrading(trading)
			}
		}
	}
}

func (r *StockLibraryMaster) GetStockLibrary(ticker string) *StockLibrary {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.StockLibraries[ticker]
}

func (r *StockLibraryMaster) GetStockLatestHolding(ticker string) *StockARKHoldings {
	r.lock.RLock()
	defer r.lock.RUnlock()

	stockLibrary := r.StockLibraries[ticker]
	if stockLibrary != nil {
		if stockLibrary.LatestStockHolding != nil {
			return stockLibrary.LatestStockHolding
		}
	}

	return nil
}

func (r *StockLibraryMaster) GetStockCurrentFundHolding(ticker, fund string) *StockHolding {
	r.lock.RLock()
	defer r.lock.RUnlock()

	stockLibrary := r.StockLibraries[ticker]
	if stockLibrary != nil {
		if stockLibrary.LatestStockHolding != nil {
			return stockLibrary.LatestStockHolding.GetFundHolding(fund)
		}
	}

	return nil
}

func (r *StockLibraryMaster) GetStockLatestTrading(ticker string) *StockARKTradings {
	r.lock.RLock()
	defer r.lock.RUnlock()

	stockLibrary := r.StockLibraries[ticker]
	if stockLibrary != nil {
		if stockLibrary.LatestStockTrading != nil {
			return stockLibrary.LatestStockTrading
		}
	}

	return nil
}

func (r *StockLibraryMaster) GetStockCurrentFundTrading(ticker, fund string) *StockTrading {
	r.lock.RLock()
	defer r.lock.RUnlock()

	stockLibrary := r.StockLibraries[ticker]
	if stockLibrary != nil {
		if stockLibrary.LatestStockTrading != nil {
			return stockLibrary.LatestStockTrading.GetFundTrading(fund)
		}
	}

	return nil
}

type StockLibrary struct {
	lock                 *sync.RWMutex
	Ticker               string
	fileStore            *utils.FileStoreSvc
	LatestStockHolding   *StockARKHoldings
	LatestStockTrading   *StockARKTradings
	HistoryStockHoldings map[time.Time]*StockARKHoldings
	HistoryStockTradings map[time.Time]*StockARKTradings
}

func NewStockLibrary(ticker string) *StockLibrary {
	r := &StockLibrary{
		Ticker:               ticker,
		fileStore:            utils.NewFileStoreSvc(stockLibraryFolder + strings.TrimSpace(ticker)),
		lock:                 &sync.RWMutex{},
		HistoryStockHoldings: make(map[time.Time]*StockARKHoldings),
		HistoryStockTradings: make(map[time.Time]*StockARKTradings),
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
	r.fileStore = utils.NewFileStoreSvc(stockLibraryFolder + strings.TrimSpace(r.Ticker))
	if r.HistoryStockHoldings == nil {
		r.HistoryStockHoldings = make(map[time.Time]*StockARKHoldings)
	}
	if r.HistoryStockTradings == nil {
		r.HistoryStockTradings = make(map[time.Time]*StockARKTradings)
	}
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
		r.HistoryStockHoldings[s.Date] = NewStockARKHoldings(r.Ticker, s.Date)
	}
	r.HistoryStockHoldings[s.Date].Add(s)
	if r.LatestStockHolding == nil || r.LatestStockHolding.Date.Before(s.Date) {
		r.LatestStockHolding = r.HistoryStockHoldings[s.Date]
	}
	r.lock.Unlock()
	r.MustSave()
}

func (r *StockLibrary) AddStockTrading(s *StockTrading) {
	r.lock.Lock()
	if r.HistoryStockTradings[s.Date] == nil {
		r.HistoryStockTradings[s.Date] = NewStockARKTradings(r.Ticker, s.Date)
	}
	r.HistoryStockTradings[s.Date].Add(s)
	if r.LatestStockTrading == nil || r.LatestStockTrading.Date.Before(s.Date) {
		r.LatestStockTrading = r.HistoryStockTradings[s.Date]
	}
	r.lock.Unlock()
	r.MustSave()
}

func (r *StockLibrary) AddStockTradingWithoutLock(s *StockTrading) {
	if r.HistoryStockTradings[s.Date] == nil {
		r.HistoryStockTradings[s.Date] = NewStockARKTradings(r.Ticker, s.Date)
	}
	r.HistoryStockTradings[s.Date].Add(s)
	r.MustSave()
}

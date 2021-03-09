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

var (
	allARKTypes = []string{"ARKF", "ARKG", "ARKK", "ARKQ", "ARKW"}
)

type ARKHoldings struct {
	Date time.Time
	ARKK *StockHoldings
	ARKQ *StockHoldings
	ARKW *StockHoldings
	ARKG *StockHoldings
	ARKF *StockHoldings
}

func NewARKHoldings() *ARKHoldings {
	return &ARKHoldings{}
}

func NewARKHoldingsFromDirectory(dir string) (*ARKHoldings, error) {
	files, err := ThePorter.ListAllCSVs(dir)
	if err != nil {
		return nil, err
	}
	var arkHoldings = NewARKHoldings()
	for _, theFile := range files {
		glog.V(10).Infof("File: %s", theFile)
		holdings, err := ThePorter.ReadCSV(theFile)
		if err != nil {
			glog.Errorf("failed to read csv file %s, err: %v", theFile, err)
			return nil, err
		}
		err = arkHoldings.AddStockHoldings(holdings)
		if err != nil {
			glog.Errorf("failed to add stock holdings, err: %v", err)
			return nil, err
		}
	}
	return arkHoldings, nil
}

func (a *ARKHoldings) Validation() bool {
	return !(a.Date.IsZero() || a.ARKK == nil || a.ARKQ == nil ||
		a.ARKW == nil || a.ARKG == nil || a.ARKF == nil)
}

func (a *ARKHoldings) AddStockHoldings(s *StockHoldings) error {
	if a.Date.IsZero() {
		a.Date = s.Date
	} else {
		if a.Date != s.Date {
			return errDateNotMatch
		}
	}

	switch s.Fund {
	case "ARKK":
		a.ARKK = s
	case "ARKQ":
		a.ARKQ = s
	case "ARKW":
		a.ARKW = s
	case "ARKG":
		a.ARKG = s
	case "ARKF":
		a.ARKF = s
	default:
		return errFundNotMatch
	}
	return nil
}

func (a *ARKHoldings) GetFundStockHoldings(fund string) *StockHoldings {
	switch fund {
	case "ARKK":
		return a.ARKK
	case "ARKQ":
		return a.ARKQ
	case "ARKW":
		return a.ARKW
	case "ARKG":
		return a.ARKG
	case "ARKF":
		return a.ARKF
	default:
		panic(fmt.Sprintf("Incorrect fund type: %s", fund))
	}
}

func (a *ARKHoldings) GenerateTrading(p *ARKHoldings) *ARKTradings {
	var (
		arkTradings = NewARKTradings()
	)
	for _, theFund := range allARKTypes {
		tradings := a.GetFundStockHoldings(theFund).GenerateTrading(p.GetFundStockHoldings(theFund))
		tradings.SetFixDirection()
		_ = arkTradings.AddStockTradings(tradings)
	}
	return arkTradings
}

type ARKTradings struct {
	Date time.Time
	ARKK *StockTradings
	ARKQ *StockTradings
	ARKW *StockTradings
	ARKG *StockTradings
	ARKF *StockTradings
}

func NewARKTradings() *ARKTradings {
	return &ARKTradings{}
}

func (a *ARKTradings) Validation() bool {
	return !(a.Date.IsZero() || a.ARKK == nil || a.ARKQ == nil ||
		a.ARKW == nil || a.ARKG == nil || a.ARKF == nil)
}

func (a *ARKTradings) AddStockTradings(s *StockTradings) error {
	if a.Date.IsZero() {
		a.Date = s.Date
	} else {
		if a.Date != s.Date {
			return errDateNotMatch
		}
	}

	switch s.Fund {
	case "ARKK":
		a.ARKK = s
	case "ARKQ":
		a.ARKQ = s
	case "ARKW":
		a.ARKW = s
	case "ARKG":
		a.ARKG = s
	case "ARKF":
		a.ARKF = s
	default:
		return errFundNotMatch
	}
	return nil
}

func (a *ARKTradings) GetFundStockTradings(fund string) *StockTradings {
	switch fund {
	case "ARKK":
		return a.ARKK
	case "ARKQ":
		return a.ARKQ
	case "ARKW":
		return a.ARKW
	case "ARKG":
		return a.ARKG
	case "ARKF":
		return a.ARKF
	default:
		panic(fmt.Sprintf("Incorrect fund type: %s", fund))
	}
}

type Library struct {
	lock                 *sync.RWMutex
	LatestStockHoldings  *ARKHoldings
	LatestStockTradings  *ARKTradings
	HistoryStockHoldings map[time.Time]*ARKHoldings
	HistoryStockTradings map[time.Time]*ARKTradings
}

func NewLibrary() *Library {
	r := &Library{
		lock:                 &sync.RWMutex{},
		LatestStockHoldings:  NewARKHoldings(),
		LatestStockTradings:  NewARKTradings(),
		HistoryStockHoldings: make(map[time.Time]*ARKHoldings),
		HistoryStockTradings: make(map[time.Time]*ARKTradings),
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
		return r.LatestStockHoldings.Date
	}

	return latestTime
}

func (r *Library) init() {
	utils.CheckFolder(libraryFolder)
	//err := r.LoadFromFileStore()
	//if err != nil {
	//	panic(fmt.Sprintf("failed to load library from the saved file, err: %v", err))
	//}
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

	return nil
}

func (r *Library) LoadFromDirectory() (err error) {
	// Load all holdings
	dates, err := ThePorter.ListAllDates()
	if err != nil {
		glog.Errorf("failed to list all csv files, err: %v", err)
		return
	}

	for _, dateFolder := range dates {
		glog.V(10).Infof("DATE_FOLDER: %s", dateFolder)
		arkHoldings, err := NewARKHoldingsFromDirectory(dateFolder)
		if err != nil {
			return err
		}
		TheLibrary.AddStockHoldings(arkHoldings)
	}

	return nil
}

func (r *Library) Save() error {
	uByte, err := json.Marshal(r)
	if err != nil {
		glog.Errorf("failed to marshal the library, err: %v", err)
		return err
	}
	glog.V(4).Infof("TO SAVE BYTES: %d", len(uByte))
	err = theLibraryFileStore.Save(uByte)
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

func (r *Library) AddStockHoldings(a *ARKHoldings) {
	r.lock.Lock()

	// Never overwrite
	if r.HistoryStockHoldings[a.Date] != nil {
		r.lock.Unlock()
		return
	}

	r.HistoryStockHoldings[a.Date] = a
	if r.LatestStockHoldings == nil || r.LatestStockHoldings.Date.Before(a.Date) {
		r.LatestStockHoldings = a
	}
	r.lock.Unlock()
	r.MustSave()
}

func (r *Library) AddStockTradings(a *ARKTradings) {
	r.lock.Lock()

	// Never overwrite
	if r.HistoryStockTradings[a.Date] != nil {
		r.lock.Unlock()
		return
	}

	r.HistoryStockTradings[a.Date] = a
	if r.LatestStockTradings == nil || r.LatestStockTradings.Date.Before(a.Date) {
		r.LatestStockTradings = a
	}
	r.lock.Unlock()
	r.MustSave()
}

func (r *Library) AddStockTradingsWithoutLock(a *ARKTradings) {
	// Never overwrite
	if r.HistoryStockTradings[a.Date] != nil {
		return
	}

	r.HistoryStockTradings[a.Date] = a
	if r.LatestStockTradings == nil || r.LatestStockTradings.Date.Before(a.Date) {
		r.LatestStockTradings = a
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
		tradings := TheLibrary.HistoryStockHoldings[dateList[i]].GenerateTrading(TheLibrary.HistoryStockHoldings[dateList[i-1]])
		r.AddStockTradingsWithoutLock(tradings)
		TheStockLibraryMaster.AddStockTradings(tradings)
	}
}

func (r *Library) GenerateCurrentTrading(holdings *ARKHoldings) {
	r.lock.Lock()
	defer r.lock.Unlock()

	tradings := holdings.GenerateTrading(TheLibrary.LatestStockHoldings)
	r.AddStockTradingsWithoutLock(tradings)
	TheStockLibraryMaster.AddStockTradings(tradings)
}

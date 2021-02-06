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
	researcherFolder     = config.Config.DataFolder + "/researcher/"
	theResearchFileStore = utils.NewFileStoreSvc(researcherFolder + "TheResearcher")
)

var (
	TheResearcher = NewResearcher()
)

type Researcher struct {
	lock                 *sync.RWMutex
	CurrentStockHoldings *StockHoldings
	LatestStockTradings  *StockTradings
	HistoryStockHoldings map[time.Time]map[string]*StockHoldings
	HistoryStockTradings map[time.Time]map[string]*StockTradings
}

func NewResearcher() *Researcher {
	r := &Researcher{
		lock:                 &sync.RWMutex{},
		CurrentStockHoldings: nil,
		LatestStockTradings:  nil,
		HistoryStockHoldings: make(map[time.Time]map[string]*StockHoldings),
		HistoryStockTradings: make(map[time.Time]map[string]*StockTradings),
	}
	r.init()
	return r
}

func (r *Researcher) init() {
	utils.CheckFolder(researcherFolder)
	err := r.LoadFromFileStore()
	if err != nil {
		panic(fmt.Sprintf("failed to load researcher from the saved file, err: %v", err))
	}
	glog.V(4).Infof("researcher init completed")
}

func (r *Researcher) LoadFromFileStore() error {
	theBytes, err := theResearchFileStore.Read()
	if err != nil {
		if os.IsNotExist(err) {
			glog.V(4).Info("No saved file for researcher")
			return nil
		}
		glog.Errorf("failed to load researcher from the saved file")
		return err
	}

	err = json.Unmarshal(theBytes, &r)
	if err != nil {
		glog.Errorf("failed to unmarshal the saved file to researcher")
		return err
	}

	glog.V(4).Infof("researcher after load: %+v", r)
	return nil
}

func (r *Researcher) Save() error {
	uByte, _ := json.Marshal(r)
	err := theResearchFileStore.Save(uByte)
	if err != nil {
		glog.Errorf("failed to save researcher, err: %v", err)
		return err
	}
	return nil
}

func (r *Researcher) MustSave() {
	err := r.Save()
	if err != nil {
		panic(err)
	}
}

func (r *Researcher) AddStockHoldings(s *StockHoldings) {
	r.lock.Lock()
	if r.HistoryStockHoldings[s.Date] == nil {
		r.HistoryStockHoldings[s.Date] = make(map[string]*StockHoldings)
	}
	r.HistoryStockHoldings[s.Date][s.Fund] = s
	if r.CurrentStockHoldings == nil || r.CurrentStockHoldings.Date.Before(s.Date) {
		r.CurrentStockHoldings = s
	}
	r.lock.Unlock()
	r.MustSave()
}

func (r *Researcher) AddStockTradings(s *StockTradings) {
	r.lock.Lock()
	if r.HistoryStockTradings[s.Date] == nil {
		r.HistoryStockTradings[s.Date] = make(map[string]*StockTradings)
	}
	r.HistoryStockTradings[s.Date][s.Fund] = s
	if r.LatestStockTradings == nil || r.LatestStockTradings.Date.Before(s.Date) {
		r.LatestStockTradings = s
	}
	r.lock.Unlock()
	r.MustSave()
}

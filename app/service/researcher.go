package service

import (
	"sync"
	"time"
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
	return &Researcher{
		lock:                 &sync.RWMutex{},
		CurrentStockHoldings: nil,
		LatestStockTradings:  nil,
		HistoryStockHoldings: make(map[time.Time]map[string]*StockHoldings),
		HistoryStockTradings: make(map[time.Time]map[string]*StockTradings),
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
}

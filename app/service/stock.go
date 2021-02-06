package service

import (
	"strconv"
	"time"
)

// Get from downloaded CSV file
type StockHolding struct {
	Date time.Time

	Fund    string
	Ticker  string
	Cusip   string
	Company string
	Shards  float64

	MarketValue float64
	Weight      float64
}

func NewStockHoldingFromRecord(record []string) *StockHolding {
	date, _ := time.Parse("1/2/2006", record[0])
	shards, _ := strconv.ParseFloat(record[5], 64)
	marketValue, _ := strconv.ParseFloat(record[6], 64)
	weight, _ := strconv.ParseFloat(record[7], 64)

	return &StockHolding{
		Date:        date,
		Fund:        record[1],
		Ticker:      record[3],
		Cusip:       record[4],
		Company:     record[2],
		Shards:      shards,
		MarketValue: marketValue,
		Weight:      weight,
	}
}

type StockHoldings struct {
	Date     time.Time
	Fund     string
	Holdings map[string]*StockHolding
}

func NewStockHoldings(date time.Time, fund string, holdings []*StockHolding) *StockHoldings {
	s := &StockHoldings{
		Date:     date,
		Fund:     fund,
		Holdings: make(map[string]*StockHolding),
	}

	for _, holding := range holdings {
		s.Holdings[holding.Fund] = holding
	}

	return s
}

// Get from E-Mail
type StockTrading struct {
	Date time.Time

	Direction string

	Fund    string
	Ticker  string
	Cusip   string
	Company string
	Shards  float64

	Weight float64
}

type StockTradings struct {
	Date     time.Time
	Fund     string
	Tradings map[string]*StockTrading
}

func NewStockTradings(date time.Time, fund string, tradings []*StockTrading) *StockTradings {
	s := &StockTradings{
		Date:     date,
		Fund:     fund,
		Tradings: make(map[string]*StockTrading),
	}

	for _, trading := range tradings {
		s.Tradings[trading.Fund] = trading
	}

	return s
}

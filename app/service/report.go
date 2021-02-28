package service

import "time"

type Report struct {
	Date         time.Time
	StockReports []*StockReport
}

type StockReport struct {
	Date                 string
	StockTicker          string
	Fund                 string
	CurrentHoldingShards float64

	// Last 3 days
	HistoryShards [3]float64

	CurrentDirection      TradeDirection
	FixDirection          TradeDirection
	CurrentTradingShards  float64
	CurrentTradingPercent float64

	FundDirection      TradeDirection
	FundTradingPercent float64
}

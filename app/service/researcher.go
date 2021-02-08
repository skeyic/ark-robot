package service

import "github.com/golang/glog"

type Researcher struct {
}

// Holding:
//  1.) total value
//  ...
//
// Trading:
//   - each trade
//      1.) start buying
//      2.) keep buying
//      3.) start selling
//      4.) keep selling
//      ...
//   - total trades
//      1.) buy more or sell more
//      ...
//  ...

func (r *Researcher) AnalyseTrading(trade *StockTrading) {

}

// BUY
// - new buy
func CheckNewBuy(trading *StockTrading, latestHoldings *StockHoldings) bool {
	return latestHoldings.Holdings[trading.Ticker] == nil
}

// SOLD
// - sold out

func CalcSoldPercent(trading *StockTrading, latestHoldings *StockHoldings) float64 {
	glog.V(6).Infof("TOTAL: %f, SOLD: %f", latestHoldings.Holdings[trading.Ticker].Shards, trading.Shards)
	return trading.Shards / latestHoldings.Holdings[trading.Ticker].Shards
}

func CheckSoldOut(trading *StockTrading, latestHoldings *StockHoldings) bool {
	return CalcSoldPercent(trading, latestHoldings) == 1
}

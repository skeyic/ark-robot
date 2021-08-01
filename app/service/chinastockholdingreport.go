package service

import (
	"fmt"
	"github.com/skeyic/ark-robot/utils"
	"sync"
	"time"
)

var (
	TheChinaStockReportMaster = NewChinaStockReportMaster()
)

type ChinaStockReportMaster struct {
	lock   *sync.RWMutex
	report string
}

func NewChinaStockReportMaster() *ChinaStockReportMaster {
	return &ChinaStockReportMaster{
		lock: &sync.RWMutex{},
	}
}

func (m *ChinaStockReportMaster) SetReport(report string) {
	m.lock.Lock()
	m.report = report
	m.lock.Unlock()
}

func (m *ChinaStockReportMaster) GetReport() string {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.report
}

// ChinaStockHoldingReport ...
// Current holding - phase 1
// Latest trading - TODO
// Data range report of last 5 days - phase 1
type ChinaStockHoldingReport struct {
	ReportDate      time.Time
	CurrentHolding  *ARKHoldings
	PreviousHolding *ARKHoldings
	details         *detailsReport
}

type stockDetail struct {
	ticker              string
	shards              float64
	previousShards      float64
	shardsDiff          float64
	marketValue         float64
	previousMarketValue float64
	marketValueDiff     float64
}

type detailsReport struct {
	details                                               map[string]*stockDetail
	allMarketValue                                        float64
	allMarketValueDiff                                    float64
	maxMarketValueTicker, maxBuyTicker, maxSellTicker     string
	maxMarketValue, maxBuyMarketValue, maxSellMarketValue float64
}

func newDetailsReport() *detailsReport {
	return &detailsReport{details: make(map[string]*stockDetail)}
}

func (r *detailsReport) Add(holding, previousHolding *StockHolding) {
	var (
		ticker              string
		shards              float64
		previousShards      float64
		shardsDiff          float64
		marketValue         float64
		previousMarketValue float64
		marketValueDiff     float64
	)

	if holding != nil {
		ticker = holding.Ticker
		shards = holding.Shards
		marketValue = holding.MarketValue
	}

	if previousHolding != nil {
		if ticker == "" {
			ticker = previousHolding.Ticker
		}
		previousShards = previousHolding.Shards
		previousMarketValue = previousHolding.MarketValue
	}

	shardsDiff = shards - previousShards
	marketValueDiff = marketValue - previousMarketValue

	if r.details[ticker] == nil {
		r.details[ticker] = &stockDetail{
			ticker:              ticker,
			shards:              shards,
			previousShards:      previousShards,
			shardsDiff:          shardsDiff,
			marketValue:         marketValue,
			previousMarketValue: previousMarketValue,
			marketValueDiff:     marketValueDiff,
		}
	} else {
		detail := r.details[ticker]
		detail.shards += shards
		detail.previousShards += previousShards
		detail.shardsDiff += shardsDiff
		detail.marketValue += marketValue
		detail.previousMarketValue += previousMarketValue
		detail.marketValueDiff += marketValueDiff
	}
}

func (r *detailsReport) report() (maxReport, detailReport string) {
	detailReport = "具体如下：\n"
	for ticker, detail := range r.details {
		if detail.marketValueDiff > r.maxBuyMarketValue {
			r.maxBuyMarketValue = detail.marketValueDiff
			r.maxBuyTicker = ticker
		} else if detail.marketValueDiff < r.maxSellMarketValue {
			r.maxSellMarketValue = detail.marketValueDiff
			r.maxSellTicker = ticker
		}

		if detail.marketValue > r.maxMarketValue {
			r.maxMarketValue = detail.marketValue
			r.maxMarketValueTicker = ticker
		}

		r.allMarketValue += detail.marketValue
		r.allMarketValueDiff += detail.marketValueDiff

		detailReport += fmt.Sprintf("  %s：持有%s股，市值%s美元，相比上个交易日", ticker,
			utils.ThousandFormatFloat64(detail.shards), utils.ThousandFormatFloat64(detail.marketValue))

		if detail.shardsDiff > 0 {
			detailReport += fmt.Sprintf("增加了%s股；\n", utils.ThousandFormatFloat64(detail.shardsDiff))
		} else if detail.shardsDiff < 0 {
			detailReport += fmt.Sprintf("减少了%s股；\n", utils.ThousandFormatFloat64(detail.shardsDiff*-1))
		} else {
			detailReport += "没有变化；\n"
		}
	}

	maxReport = fmt.Sprintf("ARK持有的中概股总市值为%s美元，相比上个交易日", utils.ThousandFormatFloat64(r.allMarketValue))
	if r.allMarketValueDiff > 0 {
		maxReport += fmt.Sprintf("增加了%s美元；\n", utils.ThousandFormatFloat64(r.allMarketValueDiff))
	} else if r.allMarketValueDiff < 0 {
		maxReport += fmt.Sprintf("减少了%s美元；\n", utils.ThousandFormatFloat64(r.allMarketValueDiff*-1))
	} else {
		maxReport += "没有变化；\n"
	}

	maxReport += fmt.Sprintf("总计持有市值最多的是%s，共计%s美元；\n", r.maxMarketValueTicker,
		utils.ThousandFormatFloat64(r.maxMarketValue))
	maxReport += fmt.Sprintf("总计持有市值增长最多的是%s，增加了%s美元；\n", r.maxBuyTicker,
		utils.ThousandFormatFloat64(r.maxBuyMarketValue))
	maxReport += fmt.Sprintf("总计持有市值减少最多的是%s，减少了%s美元；\n", r.maxSellTicker,
		utils.ThousandFormatFloat64(r.maxSellMarketValue*-1))

	return
}

func NewChinaStockHoldingReport() *ChinaStockHoldingReport {
	return &ChinaStockHoldingReport{
		ReportDate: time.Now(),
		details:    newDetailsReport(),
	}
}

func (r *ChinaStockHoldingReport) Load() error {
	latestDate := TheLibrary.GetLatestHoldingDate()
	if latestDate.IsZero() {
		return errNoLatestDate
	}

	r.CurrentHolding = TheLibrary.GetHoldings(latestDate)
	holdingsList := TheLibrary.GetPreviousHoldings(latestDate, 1)
	r.PreviousHolding = holdingsList[0]

	return nil
}

/*
ARK持有的中概股总市值为xx美元，相比上个交易日增加/减少了xxx美元；
持有市值最多的是xx，共计xxx美元；
总计持有市值增长最多的是xxx，增加了xxx美元；
总计持有市值减少最多的是xxx，减少了xxx美元。

重点操作：
  ARKK建仓了xxx，市值xxx美元；
  ARKK建仓了xxx，市值xxx美元；
  ARKX清仓了xxx；

具体如下：
  TCEHY：持有xxx股，市值xxx美元，相比上个交易日增加了xxx股；
  JD：持有xxx股，市值xxx美元，相比上个交易日减少了xxx股；
  BIDU：持有xxx股，市值xxx美元，相比上个交易日没有变化；
*/

func (r *ChinaStockHoldingReport) TxtReport() string {
	var (
		report                                                 string
		firstBuyReport, soldOutReport, maxReport, detailReport string
	)

	for _, fund := range allARKTypes {
		holdings := r.CurrentHolding.GetFundStockHoldings(fund)
		previousHoldings := r.PreviousHolding.GetFundStockHoldings(fund)
		var (
			hitMap = make(map[string]bool)
		)
		for ticker, holding := range holdings.Holdings {
			if TheChinaStockManager.IsChinaStock(ticker) {
				hitMap[ticker] = true

				var (
					previousHolding *StockHolding
				)
				if previousHoldings == nil {
					previousHolding = nil
				} else {
					previousHolding = previousHoldings.GetStockHolding(ticker)
				}

				if previousHolding == nil || previousHolding.Shards == 0 {
					firstBuyReport += fmt.Sprintf("%s建仓了%s，买入%s股，市值%s美元；\n", fund, ticker,
						utils.ThousandFormatFloat64(holding.Shards), utils.ThousandFormatFloat64(holding.MarketValue))
				}

				r.details.Add(holding, previousHolding)

			}
		}

		for ticker, holding := range previousHoldings.Holdings {
			if TheChinaStockManager.IsChinaStock(ticker) {
				if !hitMap[ticker] {
					soldOutReport += fmt.Sprintf("%s清仓了%s，卖出%s股；\n", fund, ticker, utils.ThousandFormatFloat64(holding.Shards))
					r.details.Add(nil, holding)
				}
			}
		}
	}

	maxReport, detailReport = r.details.report()

	report = maxReport + "\n" + "重点操作：\n" + firstBuyReport + soldOutReport + "\n" + detailReport

	return report
}

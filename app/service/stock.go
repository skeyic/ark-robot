package service

import (
	"fmt"
	"github.com/grd/statistics"
	"github.com/skeyic/ark-robot/utils"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Compare current holdings with previous holdings to know if a stock is "buy for the first time", "sell all"
// Compare current tradings with previous tradings to know if a stock is continuously buy or sell

type TradeDirection string

const (
	TradeEmpty        TradeDirection = ""
	TradeBuy          TradeDirection = "Buy"
	TradeRelativeBuy  TradeDirection = "Relative Buy"
	TradeSell         TradeDirection = "Sell"
	TradeRelativeSell TradeDirection = "Relative Sell"
	TradeKeep         TradeDirection = "Keep"
	TradeDoNothing    TradeDirection = "DoNothing"
)

// StockHolding Get from downloaded CSV file
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

var (
	missingTickerMap = map[string]string{
		"CM LIFE SCIENCES II INC": "CMIIU",
	}
	updateTickerMap = map[string]string{
		"3690": "MPNGY",
		"1833": "PNGAY",
	}
)

func NewStockHoldingFromRecord(record []string) *StockHolding {
	date, _ := time.Parse("1/2/2006", record[0])
	shards, _ := strconv.ParseFloat(record[5], 64)
	marketValue, _ := strconv.ParseFloat(record[6], 64)
	weight, _ := strconv.ParseFloat(record[7], 64)

	ticker := record[3]
	if ticker == "" {
		newTicker, hit := missingTickerMap[record[2]]
		if hit {
			ticker = newTicker
		} else {
			ticker = strings.ReplaceAll(record[2], " ", "_")
			ticker = strings.TrimSuffix(ticker, " ")
		}
	}

	return &StockHolding{
		Date:        date,
		Fund:        record[1],
		Ticker:      ticker,
		Cusip:       record[4],
		Company:     record[2],
		Shards:      shards,
		MarketValue: marketValue,
		Weight:      weight,
	}
}

func (s *StockHolding) ESID() string {
	return fmt.Sprintf("f%s_s%s_d%s", strings.ToLower(s.Fund), strings.ToLower(s.Ticker), s.Date.Format(TheDateIDFormat))
}

func (s *StockHolding) ESBody() string {
	return fmt.Sprintf(
		`{ "date": "%s", "fund": "%s", "ticker": "%s", "cusip": "%s", "shards": %f, "market_value": %f, "weight": %f}`,
		s.Date.Format(TheDateFormat), s.Fund, s.Ticker, s.Cusip, s.Shards, s.MarketValue, s.Weight)
}

func (s *StockHolding) Merge(d *StockHolding) {
	s.Shards += d.Shards
	s.MarketValue += d.MarketValue
	s.Weight += d.Weight
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
		//glog.V(4).Infof("NewStockHoldings: %+v", holding)
		theStock := s.Holdings[holding.Ticker]
		if theStock == nil {
			s.Holdings[holding.Ticker] = holding
		} else {
			s.Holdings[holding.Ticker].Merge(holding)
		}
	}

	return s
}

func (h *StockHoldings) GetStockHolding(ticker string) *StockHolding {
	return h.Holdings[ticker]
}

func (h *StockHoldings) GetTop10() []*StockHolding {
	var (
		top10Holdings    []*StockHolding
		idx              = 1
		toReportHoldings = make(map[float64]*StockHolding)
		toSortWeight     sort.Float64Slice
	)

	for _, holding := range h.Holdings {
		if toSkipTicker(holding.Ticker) {
			continue
		}
		weight := holding.Weight
		if _, hit := toReportHoldings[weight]; hit {
			weight += 0.000001
		}
		toReportHoldings[weight] = holding
		toSortWeight = append(toSortWeight, weight)
	}

	sort.Sort(sort.Reverse(toSortWeight))

	for _, weight := range toSortWeight {
		// Only the top 10
		if idx > maxIdx {
			break
		}

		top10Holdings = append(top10Holdings, toReportHoldings[weight])

		idx++
	}

	return top10Holdings
}

func (h *StockHoldings) GenerateTrading(p *StockHoldings) *StockTradings {
	var (
		tradings = &StockTradings{
			Date:     h.Date,
			Fund:     h.Fund,
			Tradings: make(map[string]*StockTrading),
		}
	)

	for _, holding := range h.Holdings {
		var (
			pHolding *StockHolding
		)
		if p != nil && p.Holdings != nil {
			pHolding = p.Holdings[holding.Ticker]
		}
		trading := &StockTrading{
			Date:    h.Date,
			Fund:    h.Fund,
			Ticker:  holding.Ticker,
			Cusip:   holding.Cusip,
			Company: holding.Company,
			Holding: holding.Shards,
		}
		if pHolding == nil || pHolding.Shards == 0 {
			trading.Direction = TradeBuy
			trading.Shards = holding.Shards
			trading.MarketValue = holding.MarketValue
			trading.Percent = 100.0
			trading.PreviousHolding = 0
		} else {
			trading.PreviousHolding = pHolding.Shards
			trading.Shards = holding.Shards - pHolding.Shards
			trading.MarketValue = holding.MarketValue - pHolding.MarketValue
			if pHolding.Shards < holding.Shards {
				trading.Direction = TradeBuy
				trading.Percent = trading.Shards / pHolding.Shards * 100
			} else if pHolding.Shards > holding.Shards {
				trading.Direction = TradeSell
				trading.Percent = trading.Shards / pHolding.Shards * 100
			} else {
				trading.Direction = TradeDoNothing
			}
			//glog.V(4).Infof("Ticker: %s, previous shards: %f, current shards: %f, trading shards: %f, percent: %f, direction: %s,",
			//	holding.Ticker, pHolding.Shards, holding.Shards, trading.Shards, trading.Percent, trading.Direction)
		}
		tradings.AddTrade(trading)
	}

	if p != nil {
		for _, pHolding := range p.Holdings {
			holding := h.Holdings[pHolding.Ticker]
			if holding == nil {
				trading := &StockTrading{
					Date:    h.Date,
					Fund:    h.Fund,
					Ticker:  pHolding.Ticker,
					Cusip:   pHolding.Cusip,
					Company: pHolding.Company,

					Direction: TradeSell,
					Shards:    -pHolding.Shards,
					Percent:   -100.0,

					Holding:         0,
					PreviousHolding: pHolding.Shards,
				}
				tradings.AddTrade(trading)
			}
		}
	}

	return tradings
}

// StockTrading Analyse the holding and generate the trading list
type StockTrading struct {
	Date time.Time

	Direction TradeDirection

	Fund        string
	Ticker      string
	Cusip       string
	Company     string
	Shards      float64
	Percent     float64
	MarketValue float64

	Holding         float64
	PreviousHolding float64

	FixedDirection                TradeDirection // Buy, Sell or Keep
	FixedDirectionContinuouslyDay int
}

func (t *StockTrading) IsBuy() bool {
	return t.Direction == TradeBuy
}

func (t *StockTrading) IsSell() bool {
	return t.Direction == TradeSell
}

//func NewStockTradingFromRecord(record []string) *StockTrading {
//	date, _ := time.Parse("1/2/2006", record[1])
//	shards, _ := strconv.ParseFloat(record[6], 64)
//	weight, _ := strconv.ParseFloat(record[7], 64)
//
//	return &StockTrading{
//		Date:      date,
//		Direction: record[2],
//		Fund:      record[0],
//		Ticker:    record[3],
//		Cusip:     record[4],
//		Company:   record[5],
//		Shards:    shards,
//		Percent:   weight,
//	}
//}

type TradingList []*StockTrading

func (s TradingList) Len() int {
	return len(s)
}

// Buy > Relative Buy > Sell > Relative Sell > Do Nothing > Keep
// More > Less
var directionWeightMap = map[TradeDirection]float64{
	TradeBuy:          6,
	TradeRelativeBuy:  5,
	TradeSell:         4,
	TradeRelativeSell: 3,
	TradeDoNothing:    2,
	TradeKeep:         1,
}

func (s TradingList) Less(i, j int) bool {
	if directionWeightMap[s[i].FixedDirection] == directionWeightMap[s[j].FixedDirection] {
		return s[i].Percent > s[j].Percent
	}
	return directionWeightMap[s[i].FixedDirection] > directionWeightMap[s[j].FixedDirection]
}

func (s TradingList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func RemoveAbnormalData(pl statistics.Float64) statistics.Float64 {
	var (
		npl            = statistics.Float64{}
		needCheckAgain bool
	)
	mean := statistics.Mean(&pl)
	variance := statistics.Sd(&pl)
	for _, data := range pl {
		if data > mean-3*variance && data < mean+3*variance {
			npl = append(npl, data)
		} else {
			//glog.V(10).Infof("Remove abnormal data: %f", data)
			needCheckAgain = true
		}
	}
	if needCheckAgain {
		return RemoveAbnormalData(npl)
	}
	return npl
}

const theMaxVariance = 0.005

func PickAbnormalData(pl statistics.Float64) (statistics.Float64, statistics.Float64) {
	var (
		npl            = statistics.Float64{}
		apl            = statistics.Float64{}
		needCheckAgain bool
	)

	mean := statistics.Mean(&pl)
	variance := statistics.Sd(&pl)
	//glog.V(4).Infof("MEAN: %f, VARIANCE: %f", mean, variance)
	for _, data := range pl {
		if variance < theMaxVariance || (data > mean-3*variance && data < mean+3*variance) ||
			(data/mean > 0.99 && data/mean < 1.01) {
			npl = append(npl, data)
		} else {
			//glog.V(10).Infof("Remove abnormal data: %f", data)
			apl = append(apl, data)
			needCheckAgain = true
		}
	}
	if needCheckAgain {
		normalPL, abnormalPL := PickAbnormalData(npl)
		return normalPL, append(abnormalPL, apl...)
	}
	return npl, apl
}

//func (s TradingList) SetFixDirection() {
//	var (
//		positivePercents         = statistics.Float64{}
//		negativePercents         = statistics.Float64{}
//		positiveNum, negativeNum int
//	)
//
//	for _, trading := range s {
//		if trading.Direction == TradeSell {
//			negativePercents = append(negativePercents, trading.Percent*-1)
//			negativeNum++
//		} else if trading.Direction == TradeBuy {
//			positivePercents = append(positivePercents, trading.Percent)
//			positiveNum++
//		}
//	}
//
//	allThePercents := positivePercents
//	if positiveNum < negativeNum {
//		allThePercents = negativePercents
//	}
//
//	theNormalPercents, theAbnormalPercents := PickAbnormalData(allThePercents)
//	glog.V(10).Infof("NORMAL: %+v", theNormalPercents)
//	glog.V(10).Infof("ABNORMAL: %+v", theAbnormalPercents)
//	means := statistics.Mean(&theNormalPercents)
//	glog.V(10).Infof("MEANS: %f", means)
//
//	for _, trading := range s {
//		var (
//			isKeep      = false
//			thisPercent = trading.Percent
//		)
//
//		if (positiveNum < negativeNum && trading.Direction == TradeSell) ||
//			(positiveNum > negativeNum && trading.Direction == TradeBuy) {
//			for _, normalPercent := range theNormalPercents {
//				if trading.Direction == TradeSell {
//					thisPercent *= -1
//				}
//				if thisPercent == normalPercent {
//					isKeep = true
//					break
//				}
//			}
//		}
//
//		if isKeep {
//			trading.FixedDirection = TradeKeep
//		} else {
//			if means < 0 && trading.Direction == TradeSell && trading.Percent < means*-1 {
//				trading.FixedDirection = TradeRelativeBuy
//			} else if means > 0 && trading.Direction == TradeBuy && trading.Percent < means {
//				trading.FixedDirection = TradeRelativeSell
//			} else {
//				trading.FixedDirection = trading.Direction
//			}
//		}
//	}
//}

type StockTradings struct {
	Date        time.Time
	Fund        string
	Direction   TradeDirection
	Percent     float64
	Tradings    map[string]*StockTrading
	TradingList TradingList
}

func NewStockTradings(date time.Time, fund string, tradings []*StockTrading) *StockTradings {
	s := &StockTradings{
		Date:     date,
		Fund:     fund,
		Tradings: make(map[string]*StockTrading),
	}

	for _, trading := range tradings {
		s.Tradings[trading.Ticker] = trading
		s.TradingList = append(s.TradingList, trading)
	}

	return s
}

func (s *StockTradings) AddTrade(t *StockTrading) {
	s.Tradings[t.Ticker] = t
	s.TradingList = append(s.TradingList, t)
}

func (s *StockTradings) SortedTradingList() TradingList {
	p := s.TradingList
	sort.Sort(p)
	return p
}

func (s *StockTradings) SetFixDirection() {
	var (
		positivePercents         = statistics.Float64{}
		negativePercents         = statistics.Float64{}
		positiveNum, negativeNum int
	)

	for _, trading := range s.TradingList {
		if trading.Percent < 0 {
			negativePercents = append(negativePercents, trading.Percent)
			negativeNum++
		} else if trading.Percent > 0 {
			positivePercents = append(positivePercents, trading.Percent)
			positiveNum++
		} else if trading.Direction == TradeDoNothing {
			negativePercents = append(negativePercents, 0)
			positivePercents = append(positivePercents, 0)
		}
	}

	allThePercents := positivePercents
	if positiveNum < negativeNum {
		allThePercents = negativePercents
	}

	theNormalPercents, _ := PickAbnormalData(allThePercents)
	//theNormalPercents, theAbnormalPercents := PickAbnormalData(allThePercents)
	//glog.V(10).Infof("NORMAL: %+v", theNormalPercents)
	//glog.V(10).Infof("ABNORMAL: %+v", theAbnormalPercents)
	means := statistics.Mean(&theNormalPercents)
	//glog.V(10).Infof("MEANS: %f", means)

	if means < 0 {
		s.Direction = TradeSell
	} else if means > 0 {
		s.Direction = TradeBuy
	} else {
		s.Direction = TradeDoNothing
	}
	s.Percent = means

	for _, trading := range s.TradingList {
		var (
			isKeep      = false
			thisPercent = trading.Percent
		)

		if (positiveNum < negativeNum && trading.Direction == TradeSell) ||
			(positiveNum > negativeNum && trading.Direction == TradeBuy) {
			for _, normalPercent := range theNormalPercents {
				if thisPercent == normalPercent {
					isKeep = true
					break
				}
			}
		}

		if isKeep {
			trading.FixedDirection = TradeKeep
		} else {
			//glog.V(4).Infof("thisPercent: %f, means: %f", trading.Percent, means)
			if means < 0 && trading.Direction == TradeSell && trading.Percent > means {
				trading.FixedDirection = TradeRelativeBuy
			} else if means > 0 && trading.Direction == TradeBuy && trading.Percent < means {
				trading.FixedDirection = TradeRelativeSell
			} else {
				trading.FixedDirection = trading.Direction
			}
		}
	}
}

func (s *StockTradings) GetStockTrading(ticker string) *StockTrading {
	return s.Tradings[ticker]
}

type StockARKHoldings struct {
	Ticker string
	Date   time.Time
	ARKK   *StockHolding
	ARKQ   *StockHolding
	ARKW   *StockHolding
	ARKG   *StockHolding
	ARKF   *StockHolding
	ARKX   *StockHolding
}

func NewStockARKHoldings(ticker string, date time.Time) *StockARKHoldings {
	return &StockARKHoldings{
		Ticker: ticker,
		Date:   date,
	}
}

func NewLatestStockARKHoldingsFromLibrary(ticker string) *StockARKHoldings {
	var (
		latestDate = TheLibrary.GetLatestHoldingDate()
	)

	theHoldings := TheLibrary.GetHoldings(latestDate)

	arkHoldings := &StockARKHoldings{
		Ticker: ticker,
		Date:   latestDate,
		ARKK:   theHoldings.GetStockFundStockHolding("ARKK", ticker),
		ARKQ:   theHoldings.GetStockFundStockHolding("ARKQ", ticker),
		ARKW:   theHoldings.GetStockFundStockHolding("ARKW", ticker),
		ARKG:   theHoldings.GetStockFundStockHolding("ARKG", ticker),
		ARKF:   theHoldings.GetStockFundStockHolding("ARKF", ticker),
		ARKX:   theHoldings.GetStockFundStockHolding("ARKX", ticker),
	}

	return arkHoldings
}

func (s *StockARKHoldings) Add(t *StockHolding) {
	switch t.Fund {
	case "ARKK":
		s.ARKK = t
	case "ARKQ":
		s.ARKQ = t
	case "ARKW":
		s.ARKW = t
	case "ARKG":
		s.ARKG = t
	case "ARKF":
		s.ARKF = t
	case "ARKX":
		s.ARKX = t
	default:
		panic(errFundNotMatch)
	}
}

func (s *StockARKHoldings) GetFundHolding(fund string) *StockHolding {
	switch fund {
	case "ARKK":
		return s.ARKK
	case "ARKQ":
		return s.ARKQ
	case "ARKW":
		return s.ARKW
	case "ARKG":
		return s.ARKG
	case "ARKF":
		return s.ARKF
	case "ARKX":
		return s.ARKX
	default:
		panic(errFundNotMatch)
	}
}

func (s *StockARKHoldings) TxtReport() string {
	var (
		report           string
		fundNum          int
		totalShards      float64
		totalMarketValue float64
	)

	for _, fund := range allARKTypes {
		theHolding := s.GetFundHolding(fund)
		if theHolding != nil {
			report += fmt.Sprintf("\t%s持有%s股(比重%.2f%%)，市值为%s美元。\n", theHolding.Fund,
				utils.ThousandFormatFloat64(theHolding.Shards), theHolding.Weight,
				utils.ThousandFormatFloat64(theHolding.MarketValue))
			totalShards += theHolding.Shards
			totalMarketValue += theHolding.MarketValue
			fundNum++
		}
	}

	if fundNum == 1 {
		report = fmt.Sprintf("ARK仅有一款基金持有%s：\n", s.Ticker) + report
	} else {
		report = fmt.Sprintf("ARK共有%d款基金持有%s，共计%s股，总计市值%s美元。其中：\n", fundNum, s.Ticker,
			utils.ThousandFormatFloat64(totalShards), utils.ThousandFormatFloat64(totalMarketValue)) + report
	}
	//report += strings.TrimSuffix(report, "，") + "。" + "\n"

	report = fmt.Sprintf("基于ARK基金公开的截止%s（不含）的持仓数据，", s.Date.Format(TheDateFormat)) + report

	return report
}

type StockARKTradings struct {
	Ticker string
	Date   time.Time
	ARKK   *StockTrading
	ARKQ   *StockTrading
	ARKW   *StockTrading
	ARKG   *StockTrading
	ARKF   *StockTrading
	ARKX   *StockTrading
}

func NewStockARKTradings(ticker string, date time.Time) *StockARKTradings {
	return &StockARKTradings{
		Ticker: ticker,
		Date:   date,
	}
}

func NewLatestStockARKTradingsFromLibrary(ticker string) *StockARKTradings {
	var (
		latestDate = TheLibrary.GetLatestHoldingDate()
	)

	theTradings := TheLibrary.GetTradings(latestDate)

	arkTradings := &StockARKTradings{
		Ticker: ticker,
		Date:   latestDate,
		ARKK:   theTradings.GetStockFundStockTrading("ARKK", ticker),
		ARKQ:   theTradings.GetStockFundStockTrading("ARKQ", ticker),
		ARKW:   theTradings.GetStockFundStockTrading("ARKW", ticker),
		ARKG:   theTradings.GetStockFundStockTrading("ARKG", ticker),
		ARKF:   theTradings.GetStockFundStockTrading("ARKF", ticker),
		ARKX:   theTradings.GetStockFundStockTrading("ARKX", ticker),
	}

	return arkTradings
}

func (s *StockARKTradings) Add(t *StockTrading) {
	switch t.Fund {
	case "ARKK":
		s.ARKK = t
	case "ARKQ":
		s.ARKQ = t
	case "ARKW":
		s.ARKW = t
	case "ARKG":
		s.ARKG = t
	case "ARKF":
		s.ARKF = t
	case "ARKX":
		s.ARKX = t
	default:
		panic(errFundNotMatch)
	}
}

func (s *StockARKTradings) GetFundTrading(fund string) *StockTrading {
	switch fund {
	case "ARKK":
		return s.ARKK
	case "ARKQ":
		return s.ARKQ
	case "ARKW":
		return s.ARKW
	case "ARKG":
		return s.ARKG
	case "ARKF":
		return s.ARKF
	case "ARKX":
		return s.ARKX
	default:
		panic(errFundNotMatch)
	}
}

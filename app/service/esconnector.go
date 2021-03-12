package service

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	"github.com/olivere/elastic/v7"
	"github.com/skeyic/ark-robot/config"
	"strings"
)

const (
	holdingsIndex         = "holdings"
	tradingsIndex         = "tradings"
	holdingsIndexSettings = `
{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0
  },
  "mappings": {
    "properties": {
      "date": {
        "format": "yyyy-MM-dd",
        "type": "date"
      },
      "fund": {
        "type": "keyword"
      },
      "ticker": {
        "type": "keyword"
      },
      "cusip": {
        "type": "keyword"
      },
      "company": {
        "type": "keyword"
      },
      "shards": {
        "type": "long"
      },
      "market_value": {
        "type": "long"
      },
      "weight": {
        "type": "long"
      }
    }
  }
}`
	tradingsIndexSettings = `
{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0
  },
  "mappings": {
    "properties": {
      "date": {
        "format": "yyyy-MM-dd",
        "type": "date"
      },
      "fund": {
        "type": "keyword"
      },
      "ticker": {
        "type": "keyword"
      },
      "cusip": {
        "type": "keyword"
      },
      "company": {
        "type": "keyword"
      },
      "shards": {
        "type": "long"
      },
      "percent": {
        "type": "long"
      },
      "direction": {
        "type": "keyword"
      },
      "fixed_direction": {
        "type": "keyword"
      },
      "fund_direction": {
        "type": "keyword"
      },
      "fund_percent": {
        "type": "long"
      }
    }
  }
}`
)

var (
	TheESConnector = &ESConnector{
		url: config.Config.ESServer.URL,
	}
)

type ESConnector struct {
	url string
}

func (e *ESConnector) IndexStockHoldings(holdings *StockHoldings) error {
	var (
		ctx = context.Background()
	)

	if holdings == nil || holdings.Holdings == nil {
		glog.Errorf("Empty holdings: %v", holdings)
		return errEmptySourceToIndex
	}

	client, err := elastic.NewClient(
		elastic.SetURL(e.url),
		elastic.SetHealthcheck(false),
		elastic.SetSniff(false),
	)
	if err != nil {
		glog.Errorf("New ES client to %s failed, err: %v", e.url, err)
		return err
	}

	indexExistResp, err := client.IndexExists(holdingsIndex).Do(ctx)
	if err != nil {
		glog.Errorf("Check index %s failed, err: %v", holdingsIndex, err)
		return err
	}

	if !indexExistResp {
		indexCreateResp, err := client.CreateIndex(holdingsIndex).BodyString(holdingsIndexSettings).Do(ctx)
		if err != nil || !indexCreateResp.Acknowledged {
			glog.Errorf("Create index %s failed, err: %v", holdingsIndex, err)
			return err
		}
	}

	bulkProcessor, err := elastic.NewBulkProcessorService(client).Do(ctx)
	if err != nil {
		glog.Errorf("Create bulk processor failed, err: %v", err)
		return err
	}
	defer bulkProcessor.Close()

	for _, holding := range holdings.Holdings {
		data := ToHoldingData(holding)
		//bs, _ := json.Marshal(data)
		//glog.V(4).Infof("ID: %s, BODY: %s", data.ESID(), bs)
		bulkProcessor.Add(elastic.NewBulkUpdateRequest().
			Index(holdingsIndex).Type("_doc").Id(holding.ESID()).Doc(data).DocAsUpsert(true))
	}

	err = bulkProcessor.Start(ctx)
	if err != nil {
		glog.Errorf("Start bulk processor failed, err: %v", err)
		return err
	}

	err = bulkProcessor.Flush()
	if err != nil {
		glog.Errorf("Flush bulk processor failed, err: %v", err)
		return err
	}

	resp := bulkProcessor.Stats()
	glog.V(3).Infof("PROCESSOR STATS: %+v", resp)

	return nil
}

func (e *ESConnector) IndexStockTradings(tradings *StockTradings) error {
	var (
		ctx = context.Background()
	)

	if tradings == nil || tradings.Tradings == nil {
		glog.Errorf("Empty tradings: %v", tradings)
		return errEmptySourceToIndex
	}

	client, err := elastic.NewClient(
		elastic.SetURL(e.url),
		elastic.SetHealthcheck(false),
		elastic.SetSniff(false),
	)
	if err != nil {
		glog.Errorf("New ES client to %s failed, err: %v", e.url, err)
		return err
	}

	indexExistResp, err := client.IndexExists(tradingsIndex).Do(ctx)
	if err != nil {
		glog.Errorf("Check index %s failed, err: %v", tradingsIndex, err)
		return err
	}

	if !indexExistResp {
		indexCreateResp, err := client.CreateIndex(tradingsIndex).BodyString(tradingsIndexSettings).Do(ctx)
		if err != nil || !indexCreateResp.Acknowledged {
			glog.Errorf("Create index %s failed, err: %v", tradingsIndex, err)
			return err
		}
	}

	bulkProcessor, err := elastic.NewBulkProcessorService(client).Do(ctx)
	if err != nil {
		glog.Errorf("Create bulk processor failed, err: %v", err)
		return err
	}
	defer bulkProcessor.Close()

	for _, trading := range tradings.Tradings {
		data := ToTradingData(tradings, trading)
		//bs, _ := json.Marshal(data)
		//glog.V(4).Infof("ID: %s, BODY: %s", data.ESID(), bs)
		bulkProcessor.Add(elastic.NewBulkUpdateRequest().
			Index(tradingsIndex).Type("_doc").Id(data.ESID()).Doc(data).DocAsUpsert(true))
	}

	err = bulkProcessor.Start(ctx)
	if err != nil {
		glog.Errorf("Start bulk processor failed, err: %v", err)
		return err
	}

	err = bulkProcessor.Flush()
	if err != nil {
		glog.Errorf("Flush bulk processor failed, err: %v", err)
		return err
	}

	resp := bulkProcessor.Stats()
	glog.V(3).Infof("PROCESSOR STATS: %+v", resp)

	return nil
}

type TradingData struct {
	Date           string         `json:"date"`
	Direction      TradeDirection `json:"direction"`
	Fund           string         `json:"fund"`
	Ticker         string         `json:"ticker"`
	Cusip          string         `json:"cusip"`
	Company        string         `json:"company"`
	Shards         float64        `json:"shards"`
	Percent        float64        `json:"percent"`
	FixedDirection TradeDirection `json:"fix_direction"`
	FundDirection  TradeDirection `json:"fund_direction"`
	FundPercent    float64        `json:"fund_percent"`
}

func ToTradingData(tradings *StockTradings, trading *StockTrading) *TradingData {
	data := &TradingData{
		Date:           trading.Date.Format(TheDateFormat),
		Direction:      trading.Direction,
		Fund:           trading.Fund,
		Ticker:         trading.Ticker,
		Cusip:          trading.Cusip,
		Company:        trading.Company,
		Shards:         trading.Shards,
		Percent:        trading.Percent,
		FixedDirection: trading.FixedDirection,
		FundDirection:  tradings.Direction,
		FundPercent:    tradings.Percent,
	}
	if data.Direction == TradeSell {
		data.Shards *= -1
		data.Percent *= -1
	}
	return data
}

func (t *TradingData) ESID() string {
	return fmt.Sprintf("f%s_s%s_d%s", strings.ToLower(t.Fund), strings.ToLower(t.Ticker), t.Date)
}

type HoldingData struct {
	Date        string  `json:"date"`
	Fund        string  `json:"fund"`
	Ticker      string  `json:"ticker"`
	Cusip       string  `json:"cusip"`
	Company     string  `json:"company"`
	Shards      float64 `json:"shards"`
	MarketValue float64 `json:"market_value"`
	Weight      float64 `json:"weight"`
}

func ToHoldingData(holding *StockHolding) *HoldingData {
	return &HoldingData{
		Date:        holding.Date.Format(TheDateFormat),
		Fund:        holding.Fund,
		Ticker:      holding.Ticker,
		Cusip:       holding.Cusip,
		Company:     holding.Company,
		Shards:      holding.Shards,
		MarketValue: holding.MarketValue,
		Weight:      holding.Weight,
	}
}

func (t *HoldingData) ESID() string {
	return fmt.Sprintf("f%s_s%s_d%s", strings.ToLower(t.Fund), strings.ToLower(t.Ticker), t.Date)
}

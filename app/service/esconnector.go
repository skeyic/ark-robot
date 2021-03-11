package service

import "github.com/skeyic/ark-robot/config"

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
        "format": "epoch_second",
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
)

var (
	TheESConnector = &ESConnector{
		url: config.Config.ESServer.URL,
	}
)

type ESConnector struct {
	url string
}

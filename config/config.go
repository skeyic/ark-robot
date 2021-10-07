package config

import (
	"github.com/jinzhu/configor"
)

var Config = struct {
	DebugMode bool `default:"false" env:"DEBUG_MODE"`
	Port      int  `default:"7766" env:"PORT"`
	RpcPort   int  `default:"7767" env:"RPC_PORT"`
	StaleInit bool `default:"true" env:"STALE_INIT"`
	//DataFolder     string `default:"/Users/carrick/go/src/github.com/skeyic/ark-robot/data" env:"DATA_FOLDER"`
	//ResourceFolder string `default:"/Users/carrick/go/src/github.com/skeyic/ark-robot/resource" env:"RESOURCE_FOLDER"`

	DataFolder     string `default:"C:\\Users\\15902\\go\\src\\github.com\\skeyic\\ark-robot\\data" env:"DATA_FOLDER"`
	ResourceFolder string `default:"C:\\Users\\15902\\go\\src\\github.com\\skeyic\\ark-robot\\resource" env:"RESOURCE_FOLDER"`

	NeuronServer struct {
		URL  string `default:"http://www.tanglicai.xyz:7474" env:"NEURON_SERVER_URL"`
		User string `default:"79c721a6-4d0b-4b2b-bc7c-0050fe5484a2" env:"NEURON_SERVER_USER"`
	}

	ESServer struct {
		Force bool   `default:"false" env:"ES_SERVER_FORCE"`
		URL   string `default:"http://www.tanglicai.xyz:7222" env:"ES_SERVER_URL"`
	}

	SpiderServer struct {
		URL        string `default:"http://192.168.31.32:4444" env:"SPIDER_SERVER_URL"`
		DataFolder string `default:"\\\\cocoa\\ubuntu\\spider\\data" env:"SPIDER_SERVER_DATA_FOLDER"`
	}

	Report struct {
		SpecialTradingPercent float64 `default:"3" env:"SPECIAL_TRADING_PERCENT"`
		WithExcel             bool    `default:"false" env:"REPORT_WITH_EXCEL"`
	}
}{}

func init() {
	if err := configor.Load(&Config); err != nil {
		panic(err)
	}
}

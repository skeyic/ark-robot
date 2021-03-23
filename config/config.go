package config

import (
	"github.com/jinzhu/configor"
)

var Config = struct {
	DebugMode      bool   `default:"false" env:"DEBUG_MODE"`
	Port           int    `default:"7766" env:"PORT"`
	DataFolder     string `default:"/Users/carrick/go/src/github.com/skeyic/ark-robot/data" env:"DATA_FOLDER"`
	ResourceFolder string `default:"/Users/carrick/go/src/github.com/skeyic/ark-robot/resource" env:"RESOURCE_FOLDER"`

	//DataFolder string `default:"C:\\Users\\15902\\go\\src\\github.com\\skeyic\\ark-robot\\data" env:"DATA_FOLDER"`
	//ResourceFolder string `default:"C:\\Users\\15902\\go\\src\\github.com\\skeyic\\ark-robot\\resource" env:"RESOURCE_FOLDER"`

	NeuronServer struct {
		URL  string `default:"http://www.xiaxuanli.com:7474" env:"NEURON_SERVER_URL"`
		User string `default:"79c721a6-4d0b-4b2b-bc7c-0050fe5484a2" env:"NEURON_SERVER_USER"`
	}

	ESServer struct {
		URL string `default:"http://www.xiaxuanli.com:7222" env:"ES_SERVER_URL"`
	}

	Report struct {
		SpecialTradingPercent float64 `default:"1" env:"SPECIAL_TRADING_PERCENT"`
	}
}{}

func init() {
	if err := configor.Load(&Config); err != nil {
		panic(err)
	}
}

package service

import (
	"flag"
	"github.com/golang/glog"
	"testing"
	"time"
)

func Test_LibraryInit(t *testing.T) {
	flag.Set("logtostderr", "true")
	flag.Set("v", "10")
	flag.Parse()

	glog.V(4).Infof("The LIBRARY: %+v", TheLibrary)

	date, _ := time.Parse("2006-01-02", "2021-02-05")

	//for key, value := range TheLibrary.HistoryStockHoldings[date]["ARKW"].Holdings {
	//	glog.V(4).Infof("KEY: %s, VALUE: %+v", key, value)
	//}

	for key, value := range TheLibrary.CurrentStockHoldings["ARKW"].Holdings {
		glog.V(4).Infof("KEY: %s, VALUE: %+v", key, value)
	}

	glog.V(4).Infof("HUYA: %+v", TheLibrary.HistoryStockHoldings[date]["ARKW"].Holdings["HUYA"])

}

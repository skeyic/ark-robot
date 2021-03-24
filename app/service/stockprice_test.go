package service

import (
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/utils"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestStockPriceGet(t *testing.T) {
	utils.EnableGlogForTesting()

	var (
		url = "http://q.stock.sohu.com/hisHq?code=nasdaq_aapl&start=20150504&end=20210324&stat=1&order=D&period=d&callback=historySearchHandler&rt=jsonp&r=0.8391495715053367&0.9677250558488026"
	)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// handle err
		glog.Errorf("new request failed, err: %v", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		glog.Errorf("do request failed, err: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("read response body failed, err: %v", err)
		return
	}

	glog.V(4).Infof("BODY: %s", body)
}

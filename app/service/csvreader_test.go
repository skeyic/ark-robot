package service

import (
	"flag"
	"github.com/golang/glog"
	"testing"
)

func TestCSVRead(t *testing.T) {
	flag.Set("logtostderr", "true")
	flag.Set("v", "10")
	flag.Parse()

	r := NewCSVReader("C:\\Users\\15902\\go\\src\\github.com\\skeyic\\ark-robot\\data\\downloader\\2021-02-06-11-57-53-ARKK.csv")
	records, err := r.Load()
	if err != nil {
		glog.Errorf("failed to load, error: %v", err)
		return
	}
	glog.V(4).Infof("RECORDS: %+v", records)

	glog.V(4).Infof("CSV Validate: %v", ValidateARKCSV(records))

	for idx, record := range records[1:] {
		glog.V(4).Infof("RECORD IDX: %d, stockHolding: %+v, value: %v", idx, NewStockHoldingFromRecord(record), record)
	}
}

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

	r := &CSVReader{filepath: generateFilePath("ARKK")}
	records, err := r.Load()
	if err != nil {
		glog.Errorf("failed to load, error: %v", err)
		return
	}
	glog.V(4).Infof("RECORDS: %+v", records)

	if len(records) >= 2 {
		glog.V(4).Infof(records[1][0])
	}
}

package service

import (
	"flag"
	"github.com/golang/glog"
	"testing"
)

func TestDownloadCSV(t *testing.T) {
	flag.Set("logtostderr", "true")
	flag.Set("v", "10")
	flag.Parse()

	err := TheDownloader.DownloadAllARKCSVs()
	if err != nil {
		glog.Errorf("failed to download csv, err: %v", err)
		return
	}

	//<-make(chan struct{}, 1)
}

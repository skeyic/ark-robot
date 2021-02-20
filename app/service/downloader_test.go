package service

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"testing"
	"time"
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

func TestDownloadTime(t *testing.T) {
	fmt.Printf("NOW: %d", time.Now().UTC().Day())
}

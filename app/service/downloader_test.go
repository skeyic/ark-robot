package service

import (
	"github.com/golang/glog"
	"testing"
)

func TestDownloadCSV(t *testing.T) {
	var (
		arkType  = "ARKK"
		fileName = generateFilePath(arkType)
	)
	err := DownloadARKCSV(generateArkCSVURL(arkType), fileName)
	if err != nil {
		glog.Errorf("failed to download csv, err: %v", err)
		return
	}
	glog.V(4).Infof("download csv successfully, filename: %s", fileName)
}

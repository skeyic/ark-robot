package service

import (
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/utils"
	"testing"
)

func TestListAllCSVs(t *testing.T) {
	utils.EnableGlogForTesting()

	files, err := ThePorter.ListAllCSVs()
	if err != nil {
		glog.Errorf("failed to list all csv files, err: %v", err)
		return
	}

	glog.V(4).Infof("BEFORE READ: %+v", TheLibrary)

	for _, theFile := range files {
		glog.V(4).Infof("File: %s", theFile)
		err = ThePorter.ReadCSV(theFile)
		if err != nil {
			glog.Errorf("failed to read csv file %s, err: %v", theFile, err)
			return
		}
	}

	glog.V(4).Infof("AFTER READ: %+v", TheLibrary)
}

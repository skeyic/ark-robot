package service

import (
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/utils"
	"testing"
)

func TestListAllCSVs(t *testing.T) {
	utils.EnableGlogForTesting()

	//files, err := ThePorter.ListAllCSVs()
	//if err != nil {
	//	glog.Errorf("failed to list all csv files, err: %v", err)
	//	return
	//}
	//
	//glog.V(4).Infof("BEFORE READ: %+v", TheLibrary)
	//
	//for _, theFile := range files {
	//	glog.V(4).Infof("File: %s", theFile)
	//	_, err = ThePorter.ReadCSV(theFile)
	//	if err != nil {
	//		glog.Errorf("failed to read csv file %s, err: %v", theFile, err)
	//		return
	//	}
	//}
	//
	//glog.V(4).Infof("AFTER READ: %+v", TheLibrary)
}

func TestListAllDates(t *testing.T) {
	utils.EnableGlogForTesting()

	files, err := ThePorter.ListAllDates()
	if err != nil {
		glog.Errorf("failed to list all csv files, err: %v", err)
		return
	}

	for _, theFile := range files {
		glog.V(4).Infof("File: %s", theFile)
	}
}

func TestPorter_Catalog(t *testing.T) {
	ThePorter.Catalog("C:\\Users\\15902\\Downloads\\ARK105\\ARK_INNOVATION_ETF_ARKK_HOLDINGS.csv")
	ThePorter.Catalog("C:\\Users\\15902\\Downloads\\ARK105\\ARK_INNOVATION_ETF_ARKX_HOLDINGS.csv")
	ThePorter.Catalog("C:\\Users\\15902\\Downloads\\ARK105\\ARK_INNOVATION_ETF_ARKG_HOLDINGS.csv")
	ThePorter.Catalog("C:\\Users\\15902\\Downloads\\ARK105\\ARK_INNOVATION_ETF_ARKK_HOLDINGS.csv")
	ThePorter.Catalog("C:\\Users\\15902\\Downloads\\ARK105\\ARK_INNOVATION_ETF_ARKQ_HOLDINGS.csv")
	ThePorter.Catalog("C:\\Users\\15902\\Downloads\\ARK105\\ARK_INNOVATION_ETF_ARKW_HOLDINGS.csv")
}

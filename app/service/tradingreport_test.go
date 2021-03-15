package service

import (
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/utils"
	"testing"
)
import "github.com/360EntSecGroup-Skylar/excelize/v2"

func TestExcel(t *testing.T) {
	utils.EnableGlogForTesting()

	var (
		err  error
		date = "2021-03-01"
		r    = TradingsReport{
			Date: date,
		}
		fileName = r.ExcelPath()
	)

	r.InitExcelFromTemplate()
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		glog.Errorf("failed to open excel %s, err: %v", fileName, err)
		return
	}

	f.SetCellValue(tradingsSheet, "A4", "HUYA")
	f.SetCellValue(tradingsSheet, "B4", "ARKK")
	f.SetCellValue(tradingsSheet, "C4", "买入")
	f.SetCellValue(tradingsSheet, "E4", 1000)
	f.SetCellValue(tradingsSheet, "F4", "10%")

	err = f.Save()
	if err != nil {
		glog.Errorf("failed to save excel %s, err: %v", fileName, err)
		return
	}

}

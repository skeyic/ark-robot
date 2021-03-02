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
		r    = Report{
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

	f.SetCellValue(sheet, "A4", "HUYA")
	f.SetCellValue(sheet, "B4", "ARKK")
	f.SetCellValue(sheet, "C4", "买入")
	f.SetCellValue(sheet, "E4", 1000)
	f.SetCellValue(sheet, "F4", "10%")

	err = f.Save()
	if err != nil {
		glog.Errorf("failed to save excel %s, err: %v", fileName, err)
		return
	}

}

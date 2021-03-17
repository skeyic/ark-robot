package service

import (
	"fmt"
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

	f.SetCellValue(defaultSheet, "A4", "HUYA")
	f.SetCellValue(defaultSheet, "B4", "ARKK")
	f.SetCellValue(defaultSheet, "C4", "买入")
	f.SetCellValue(defaultSheet, "E4", 1000)
	f.SetCellValue(defaultSheet, "F4", "10%")

	err = f.Save()
	if err != nil {
		glog.Errorf("failed to save excel %s, err: %v", fileName, err)
		return
	}

}

func TestPrintf(t *testing.T) {
	fmt.Printf("%.2f\n", 12.345)
	fmt.Printf("%.2g\n", 12.345)
	fmt.Println(fmt.Sprintf("%.2g", 12.345))
	fmt.Println(floatToStringIntOnly(12.345))
}

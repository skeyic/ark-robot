package service

import (
	"errors"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/config"
	"github.com/skeyic/ark-robot/utils"
	"strconv"
)

var (
	reportPath    = config.Config.DataFolder + "/report"
	excelTemplate = reportPath + "/ARK.xlsx"
	sheet         = "sheet"
)

var (
	ErrInitReportFile = errors.New("failed to init report file")
)

func init() {
	utils.CheckFolder(reportPath)
	utils.CheckFile(excelTemplate)
}

type Report struct {
	Date         string
	StockReports []*StockReport
}

func (r *Report) ToExcel() error {
	var (
		err      error
		fileName = r.ExcelPath()
	)

	if len(r.StockReports) == 0 {
		glog.Warningf("Empty report")
		return nil
	}

	err = r.InitExcelFromTemplate()
	if err != nil {
		glog.Errorf("failed to init excel from template, err: %v", err)
		return err
	}

	f, err := excelize.OpenFile(fileName)
	if err != nil {
		glog.Errorf("failed to open excel %s, err: %v", fileName, err)
		return err
	}

	for idx, stockReport := range r.StockReports {
		// Leave the example to test
		line := strconv.Itoa(idx + 4)
		f.SetCellValue(sheet, "A"+line, stockReport.StockTicker)
		f.SetCellValue(sheet, "B"+line, stockReport.Fund)
		f.SetCellValue(sheet, "C"+line, stockReport.CurrentDirection)
		f.SetCellValue(sheet, "D"+line, stockReport.FixDirection)
		f.SetCellValue(sheet, "E"+line, stockReport.CurrentTradingShards)
		f.SetCellValue(sheet, "F"+line, stockReport.CurrentTradingPercent)
		f.SetCellValue(sheet, "G"+line, stockReport.CurrentHoldingShards)
		//f.SetCellValue(sheet, "H"+line, stockReport.CurrentDirection)
		//f.SetCellValue(sheet, "I"+line, stockReport.CurrentDirection)
		//f.SetCellValue(sheet, "J"+line, stockReport.CurrentDirection)
		f.SetCellValue(sheet, "K"+line, stockReport.FundDirection)
		f.SetCellValue(sheet, "L"+line, stockReport.FundTradingPercent)
	}

	err = f.Save()
	if err != nil {
		glog.Errorf("failed to save excel %s, err: %v", fileName, err)
		return err
	}

	return nil
}

func (r *Report) ExcelPath() string {
	return reportPath + "/" + r.ExcelName()
}

func (r *Report) ExcelName() string {
	return r.Date + ".xlsx"
}

func (r *Report) InitExcelFromTemplate() error {
	var fileName = r.ExcelPath()
	if utils.CheckFileExist(fileName) {
		utils.DeleteFile(fileName)
	}
	utils.CopyFile(excelTemplate, fileName)
	//var i = 0
	//for {
	//	time.Sleep(2 * time.Second)
	//	if utils.CheckFileExist(fileName) {
	//		break
	//	}
	//	i++
	//	if i > 5 {
	//		return ErrInitReportFile
	//	}
	//}
	glog.V(4).Infof("Init fileName: %s", fileName)
	return nil
}

type StockReport struct {
	Date                 string
	StockTicker          string
	Fund                 string
	CurrentHoldingShards float64
	CurrentHoldingWeight float64

	// Last 3 days
	HistoryShards [3]float64

	CurrentDirection      TradeDirection
	FixDirection          TradeDirection
	CurrentTradingShards  float64
	CurrentTradingPercent float64

	FundDirection      TradeDirection
	FundTradingPercent float64
}

package service

import (
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/config"
	"github.com/skeyic/ark-robot/utils"
)

var (
	reportPath    = config.Config.DataFolder + "/report"
	excelTemplate = reportPath + "/ARK.xlsx"
	sheet         = "sheet"
)

func init() {
	utils.CheckFolder(reportPath)
	utils.CheckFile(excelTemplate)
}

type Report struct {
	Date         string
	StockReports []*StockReport
}

func (r *Report) ToExcel() {

}

func (r *Report) ExcelPath() string {
	return reportPath + r.ExcelName()
}

func (r *Report) ExcelName() string {
	return r.Date + ".xlsx"
}

func (r *Report) InitExcelFromTemplate() {
	var fileName = r.ExcelPath()
	if utils.CheckFileExist(fileName) {
		utils.DeleteFile(fileName)
	}
	utils.CopyFile(excelTemplate, fileName)
	glog.V(4).Infof("FileName: %s", fileName)
}

type StockReport struct {
	Date                 string
	StockTicker          string
	Fund                 string
	CurrentHoldingShards float64

	// Last 3 days
	HistoryShards [3]float64

	CurrentDirection      TradeDirection
	FixDirection          TradeDirection
	CurrentTradingShards  float64
	CurrentTradingPercent float64

	FundDirection      TradeDirection
	FundTradingPercent float64
}

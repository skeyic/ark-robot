package service

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/config"
	"github.com/skeyic/ark-robot/utils"
	"strconv"
)

var (
	reportPath    = config.Config.DataFolder + "/report"
	excelTemplate = config.Config.ResourceFolder + "/ARK.xlsx"
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

func toSkipTrade(direction TradeDirection) bool {
	//return false
	return direction == TradeDoNothing || direction == TradeKeep
}

func toSkipTicker(ticker string) bool {
	return ticker == "MORGAN STANLEY GOVT INSTL 8035"
}

func toPercentString(percent float64) string {
	return fmt.Sprintf("%.3f", percent) + "%"
}

func (r *Report) ToExcel(full bool) error {
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

	var idx = 4
	for _, stockReport := range r.StockReports {
		if !full && (toSkipTrade(stockReport.FixDirection) ||
			toSkipTicker(stockReport.StockTicker)) {
			continue
		}

		previousHoldings := TheStockLibraryMaster.GetStockPreviousHoldings(stockReport.StockTicker, stockReport.Fund, 3)
		glog.V(4).Infof("PREVIOUSHOLDINGS: %+v", previousHoldings)
		getHoldingShards := func(holding *StockHolding) float64 {
			if holding == nil {
				return 0
			}
			return holding.Shards
		}

		// Leave the example to test
		line := strconv.Itoa(idx)
		f.SetCellValue(sheet, "A"+line, stockReport.StockTicker)
		f.SetCellValue(sheet, "B"+line, stockReport.Company)
		f.SetCellValue(sheet, "C"+line, stockReport.Cusip)
		f.SetCellValue(sheet, "D"+line, stockReport.Fund)
		f.SetCellValue(sheet, "E"+line, stockReport.CurrentDirection)
		f.SetCellValue(sheet, "F"+line, stockReport.FixDirection)
		f.SetCellValue(sheet, "G"+line, stockReport.CurrentTradingShards)
		f.SetCellValue(sheet, "H"+line, toPercentString(stockReport.CurrentTradingPercent))
		f.SetCellValue(sheet, "I"+line, stockReport.CurrentHoldingShards)
		f.SetCellValue(sheet, "J"+line, getHoldingShards(previousHoldings[0]))
		f.SetCellValue(sheet, "K"+line, getHoldingShards(previousHoldings[1]))
		f.SetCellValue(sheet, "L"+line, getHoldingShards(previousHoldings[2]))
		f.SetCellValue(sheet, "M"+line, stockReport.FundDirection)
		f.SetCellValue(sheet, "N"+line, toPercentString(stockReport.FundTradingPercent))

		idx++
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
	Company              string
	Cusip                string
	Fund                 string
	CurrentHoldingShards float64
	CurrentHoldingWeight float64

	CurrentDirection      TradeDirection
	FixDirection          TradeDirection
	CurrentTradingShards  float64
	CurrentTradingPercent float64

	FundDirection      TradeDirection
	FundTradingPercent float64
}

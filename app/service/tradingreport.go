package service

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/utils"
	"strconv"
	"time"
)

type TradingsReport struct {
	Date         string
	StockReports []*StockReport
}

func NewTradingsReport(date time.Time) *TradingsReport {
	var (
		r = &TradingsReport{
			Date: date.Format(TheDateFormat),
		}
	)
	tradings := TheLibrary.GetTradings(date)
	if tradings == nil {
		return r
	}

	utils.CheckFolder(r.ExcelFolder())

	for _, fund := range allARKTypes {
		tradings := tradings.GetFundStockTradings(fund)
		for _, trading := range tradings.SortedTradingList() {
			stockCurrentHoldings := TheStockLibraryMaster.GetStockCurrentHolding(trading.Ticker, trading.Fund)
			r.StockReports = append(r.StockReports, &StockReport{
				Date:                  trading.Date.Format(TheDateFormat),
				StockTicker:           trading.Ticker,
				Company:               trading.Company,
				Cusip:                 trading.Cusip,
				Fund:                  trading.Fund,
				CurrentHoldingShards:  stockCurrentHoldings.Shards,
				CurrentDirection:      trading.Direction,
				FixDirection:          trading.FixedDirection,
				CurrentTradingShards:  trading.Shards,
				CurrentTradingPercent: trading.Percent,
				FundDirection:         tradings.Direction,
				FundTradingPercent:    tradings.Percent,
			},
			)
		}
	}

	return r
}

func (r *TradingsReport) ToExcel(full bool) error {
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
		getHoldingShards := func(holding *StockHolding) float64 {
			if holding == nil {
				return 0
			}
			return holding.Shards
		}

		// Leave the example to test
		line := strconv.Itoa(idx)
		f.SetCellValue(tradingsSheet, "A"+line, stockReport.StockTicker)
		f.SetCellValue(tradingsSheet, "B"+line, stockReport.Company)
		f.SetCellValue(tradingsSheet, "C"+line, stockReport.Cusip)
		f.SetCellValue(tradingsSheet, "D"+line, stockReport.Fund)
		f.SetCellValue(tradingsSheet, "E"+line, stockReport.CurrentDirection)
		f.SetCellValue(tradingsSheet, "F"+line, stockReport.FixDirection)
		f.SetCellValue(tradingsSheet, "G"+line, stockReport.CurrentTradingShards)
		f.SetCellValue(tradingsSheet, "H"+line, floatToPercentString(stockReport.CurrentTradingPercent))
		f.SetCellValue(tradingsSheet, "I"+line, stockReport.CurrentHoldingShards)
		f.SetCellValue(tradingsSheet, "J"+line, getHoldingShards(previousHoldings[0]))
		f.SetCellValue(tradingsSheet, "K"+line, getHoldingShards(previousHoldings[1]))
		f.SetCellValue(tradingsSheet, "L"+line, getHoldingShards(previousHoldings[2]))
		f.SetCellValue(tradingsSheet, "M"+line, stockReport.FundDirection)
		f.SetCellValue(tradingsSheet, "N"+line, floatToPercentString(stockReport.FundTradingPercent))

		idx++
	}

	err = f.Save()
	if err != nil {
		glog.Errorf("failed to save excel %s, err: %v", fileName, err)
		return err
	}

	glog.V(4).Infof("TradingsReport %s is provided", fileName)

	return nil
}

func (r *TradingsReport) ExcelFolder() string {
	return reportPath + "/" + r.Date
}

func (r *TradingsReport) ExcelPath() string {
	return r.ExcelFolder() + "/" + r.ExcelName()
}

func (r *TradingsReport) ExcelName() string {
	return "tradings_" + r.Date + ".xlsx"
}

func (r *TradingsReport) InitExcelFromTemplate() error {
	var fileName = r.ExcelPath()
	if utils.CheckFileExist(fileName) {
		utils.DeleteFile(fileName)
	}
	utils.CopyFile(tradingsExcelTemplate, fileName)
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

package service

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/utils"
	"strconv"
	"time"
)

type TradingsReport struct {
	Date             string
	tradings         *ARKTradings
	holdings         *ARKHoldings
	previousHoldings []*ARKHoldings
}

func NewTradingsReport(date time.Time) *TradingsReport {
	var (
		r = &TradingsReport{
			Date:             date.Format(TheDateFormat),
			tradings:         TheLibrary.GetTradings(date),
			holdings:         TheLibrary.GetHoldings(date),
			previousHoldings: TheLibrary.GetPreviousHoldings(date, 3),
		}
	)

	utils.CheckFolder(r.ExcelFolder())

	return r
}

func (r *TradingsReport) ToExcel(full bool) error {
	var (
		err      error
		fileName = r.ExcelPath()
	)

	if r.tradings == nil {
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

	var (
		sheet = defaultSheet
		idx   = 3
	)

	getHoldingShards := func(holding *ARKHoldings, fund, ticker string) float64 {
		if holding == nil {
			return 0
		}
		if holding.GetFundStockHoldings(fund) == nil {
			return 0
		}
		stockHolding := holding.GetFundStockHoldings(fund).Holdings[ticker]
		if stockHolding == nil {
			return 0
		}
		return stockHolding.Shards
	}

	for _, fund := range allARKTypes {
		var (
			toReportTradings []*StockTrading
		)

		tradings := r.tradings.GetFundStockTradings(fund)
		if tradings == nil {
			return errEmptyReport
		}

		for _, trading := range tradings.Tradings {
			if toSkipTicker(trading.Ticker) {
				continue
			}
			if !full && toSkipTrade(trading.FixedDirection) {
				continue
			}
			toReportTradings = append(toReportTradings, trading)
		}

		for _, trading := range toReportTradings {
			line := strconv.Itoa(idx)
			f.SetCellValue(sheet, "A"+line, trading.Ticker)
			f.SetCellValue(sheet, "B"+line, trading.Company)
			f.SetCellValue(sheet, "C"+line, trading.Cusip)
			f.SetCellValue(sheet, "D"+line, trading.Fund)
			f.SetCellValue(sheet, "E"+line, trading.Direction)
			f.SetCellValue(sheet, "F"+line, trading.FixedDirection)
			f.SetCellValue(sheet, "G"+line, trading.Shards)
			f.SetCellValue(sheet, "H"+line, floatToPercentStringWithSign(trading.Percent))
			f.SetCellValue(sheet, "I"+line, getHoldingShards(r.holdings, fund, trading.Ticker))
			f.SetCellValue(sheet, "J"+line, getHoldingShards(r.previousHoldings[0], fund, trading.Ticker))
			f.SetCellValue(sheet, "K"+line, getHoldingShards(r.previousHoldings[1], fund, trading.Ticker))
			f.SetCellValue(sheet, "L"+line, getHoldingShards(r.previousHoldings[2], fund, trading.Ticker))
			f.SetCellValue(sheet, "M"+line, tradings.Direction)
			f.SetCellValue(sheet, "N"+line, floatToPercentStringWithSign(tradings.Percent))

			idx++
		}
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
	return prefixTradings + r.Date + ".xlsx"
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

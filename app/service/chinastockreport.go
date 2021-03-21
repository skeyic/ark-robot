package service

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/utils"
	"strconv"
	"time"
)

type ChinaStockTradingsReport struct {
	Date     string
	tradings *ARKTradings
	holdings *ARKHoldings
	//previousHoldings []*ARKHoldings
}

func NewChinaStockTradingsReport(date time.Time) *ChinaStockTradingsReport {
	var (
		r = &ChinaStockTradingsReport{
			Date:     date.Format(TheDateFormat),
			tradings: TheLibrary.GetTradings(date),
			holdings: TheLibrary.GetHoldings(date),
			//previousHoldings: TheLibrary.GetPreviousHoldings(date, 3),
		}
	)

	utils.CheckFolder(r.ExcelFolder())

	return r
}

func (r *ChinaStockTradingsReport) ToExcel() error {
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
		idx   = 2
	)

	getHoldingWeight := func(holding *ARKHoldings, fund, ticker string) float64 {
		if holding == nil {
			return 0
		}
		stockHolding := holding.GetFundStockHoldings(fund).Holdings[ticker]
		if stockHolding == nil {
			return 0
		}
		return stockHolding.Weight
	}

	for _, fund := range allARKTypes {
		var (
			toReportChinaStockTradings []*StockTrading
		)

		tradings := r.tradings.GetFundStockTradings(fund)
		if tradings == nil {
			return errEmptyReport
		}

		for _, trading := range tradings.Tradings {
			if !TheChinaStockManager.IsChinaStock(trading.Ticker) {
				continue
			}
			toReportChinaStockTradings = append(toReportChinaStockTradings, trading)
		}

		for _, trading := range toReportChinaStockTradings {
			line := strconv.Itoa(idx)
			f.SetCellValue(sheet, "A"+line, trading.Fund)
			f.SetCellValue(sheet, "B"+line, trading.Ticker)
			f.SetCellValue(sheet, "C"+line, trading.Company)
			f.SetCellValue(sheet, "D"+line, floatToPercentStringWithSign(trading.Percent))
			f.SetCellValue(sheet, "E"+line, floatToPercentString(getHoldingWeight(r.holdings, fund, trading.Ticker)))
			//f.SetCellValue(sheet, "J"+line, getHoldingShards(r.previousHoldings[0], fund, trading.Ticker))
			//f.SetCellValue(sheet, "K"+line, getHoldingShards(r.previousHoldings[1], fund, trading.Ticker))
			//f.SetCellValue(sheet, "L"+line, getHoldingShards(r.previousHoldings[2], fund, trading.Ticker))

			idx++
		}
	}

	err = f.Save()
	if err != nil {
		glog.Errorf("failed to save excel %s, err: %v", fileName, err)
		return err
	}

	glog.V(4).Infof("ChinaStockTradingsReport %s is provided", fileName)

	return nil
}

func (r *ChinaStockTradingsReport) ExcelFolder() string {
	return reportPath + "/" + r.Date
}

func (r *ChinaStockTradingsReport) ExcelPath() string {
	return r.ExcelFolder() + "/" + r.ExcelName()
}

func (r *ChinaStockTradingsReport) ExcelName() string {
	return prefixChinaStockTradings + r.Date + ".xlsx"
}

func (r *ChinaStockTradingsReport) InitExcelFromTemplate() error {
	var fileName = r.ExcelPath()
	if utils.CheckFileExist(fileName) {
		utils.DeleteFile(fileName)
	}
	utils.CopyFile(chinaStockExcelTradingsTemplate, fileName)
	glog.V(4).Infof("Init fileName: %s", fileName)
	return nil
}

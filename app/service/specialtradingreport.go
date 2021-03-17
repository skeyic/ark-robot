package service

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/utils"
	"strconv"
	"time"
)

type SpecialTradingsReport struct {
	Date             string
	tradings         *ARKTradings
	previousTradings []*ARKTradings
}

func NewSpecialTradingsReport(date time.Time) *SpecialTradingsReport {
	var (
		r = &SpecialTradingsReport{
			Date:             date.Format(TheDateFormat),
			tradings:         TheLibrary.GetTradings(date),
			previousTradings: TheLibrary.GetPreviousTradings(date, 2),
		}
	)

	utils.CheckFolder(r.ExcelFolder())

	return r
}

func (r *SpecialTradingsReport) ToExcel() error {
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

	getTradingPercent := func(trading *ARKTradings, fund, ticker string) float64 {
		if trading == nil {
			return 0
		}
		stockTrading := trading.GetFundStockTradings(fund).Tradings[ticker]
		if stockTrading == nil {
			return 0
		}
		return stockTrading.Percent
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
			if !IsSpecialTradings(trading) {
				continue
			}
			toReportTradings = append(toReportTradings, trading)
		}

		for _, trading := range toReportTradings {

			line := strconv.Itoa(idx)
			f.SetCellValue(sheet, "A"+line, trading.Fund)
			f.SetCellValue(sheet, "B"+line, trading.Ticker)
			f.SetCellValue(sheet, "C"+line, floatToStringIntOnlyWithSign(trading.Shards))
			f.SetCellValue(sheet, "D"+line, floatToPercentStringWithSign(trading.Percent))
			f.SetCellValue(sheet, "E"+line, floatToPercentStringWithSign(getTradingPercent(r.previousTradings[1], fund, trading.Ticker)))
			f.SetCellValue(sheet, "F"+line, floatToPercentStringWithSign(getTradingPercent(r.previousTradings[0], fund, trading.Ticker)))

			idx++
		}
	}

	err = f.Save()
	if err != nil {
		glog.Errorf("failed to save excel %s, err: %v", fileName, err)
		return err
	}

	glog.V(4).Infof("SpecialTradingsReport %s is provided", fileName)

	return nil
}

func (r *SpecialTradingsReport) ExcelFolder() string {
	return reportPath + "/" + r.Date
}

func (r *SpecialTradingsReport) ExcelPath() string {
	return r.ExcelFolder() + "/" + r.ExcelName()
}

func (r *SpecialTradingsReport) ExcelName() string {
	return prefixSpecialTradings + r.Date + ".xlsx"
}

func (r *SpecialTradingsReport) InitExcelFromTemplate() error {
	var fileName = r.ExcelPath()
	if utils.CheckFileExist(fileName) {
		utils.DeleteFile(fileName)
	}
	utils.CopyFile(specialTradingsExcelTemplate, fileName)
	glog.V(4).Infof("Init fileName: %s", fileName)
	return nil
}

func IsSpecialTradings(trading *StockTrading) bool {
	return trading.Percent >= 10 || trading.Percent <= -10
}

package service

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/utils"
	"sort"
	"strconv"
	"time"
)

type StockReport struct {
	Ticker   string
	FromDate time.Time
	EndDate  time.Time
}

func NewStockReport(ticker string, fromDate, endDate time.Time) *StockReport {
	return &StockReport{
		Ticker:   ticker,
		FromDate: fromDate,
		EndDate:  endDate,
	}
}

func (r *StockReport) ToExcel() error {
	var (
		err      error
		fileName = r.ExcelPath()
	)

	stock := TheStockLibraryMaster.StockLibraries[r.Ticker]
	if stock == nil {
		return errStockNotFound
	}

	var (
		dateList timeList
	)

	for theDate := range stock.HistoryStockTradings {
		if theDate.After(r.FromDate) && theDate.Before(r.EndDate) {
			dateList = append(dateList, theDate)
		}
	}
	if len(dateList) == 0 {
		return errNoDataInDateRange
	}

	sort.Sort(dateList)

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
		idx   = 21
	)

	for i := 0; i < len(dateList); i++ {
		var (
			shards float64
		)
		holdings := stock.HistoryStockHoldings[dateList[i]]
		for _, fund := range allARKTypes {
			holding := holdings[fund]
			if holding != nil {
				//glog.V(4).Infof("DATE: %s, FUND: %s, FD: %s, SHARDS: %f, PERCENT: %f", dateList[i], fund, fundTradings.FixedDirection, fundTradings.Shards, fundTradings.Percent)
				shards += holding.Shards
			}
		}

		line := strconv.Itoa(idx)
		f.SetCellValue(sheet, "A"+line, r.Ticker)
		f.SetCellValue(sheet, "B"+line, dateList[i].Format(TheDateFormat))
		f.SetCellValue(sheet, "C"+line, fmt.Sprintf("%.0f", shards))

		idx++
		glog.V(4).Infof("IDX: %d, TICKER: %s, DATE: %s, SHARDS: %.0f", idx, r.Ticker, dateList[i].Format(TheDateFormat), shards)
	}

	err = f.Save()
	if err != nil {
		glog.Errorf("failed to save excel %s, err: %v", fileName, err)
		return err
	}

	return nil
}

func (r *StockReport) ReportFolder() string {
	return stockReportPath
}

func (r *StockReport) ExcelPath() string {
	return r.ReportFolder() + "/" + r.ExcelName()
}

func (r *StockReport) ExcelName() string {
	return fmt.Sprintf("%s%s_from_%s_to_%s.xlsx", prefixStockReport, r.Ticker,
		r.FromDate.Format(TheDateFormat), r.EndDate.Format(TheDateFormat))
}

func (r *StockReport) InitExcelFromTemplate() error {
	var fileName = r.ExcelPath()
	if utils.CheckFileExist(fileName) {
		utils.DeleteFile(fileName)
	}
	utils.CopyFile(stockReportExcelTemplate, fileName)
	glog.V(4).Infof("Init fileName: %s", fileName)
	return nil
}

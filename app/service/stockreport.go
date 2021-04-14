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
	Details  *stockDetails
}

func NewStockReport(ticker string, fromDate, endDate time.Time) *StockReport {
	return &StockReport{
		Ticker:   ticker,
		FromDate: fromDate,
		EndDate:  endDate,
	}
}

type stockDetails struct {
	dateList    []time.Time
	fundList    []string
	dailyDetail map[time.Time]*stockDailyDetail
}

type stockDailyDetail struct {
	date     time.Time
	holdings map[string]*StockHolding
	tradings map[string]*StockTrading
}

func (r *StockReport) Load() error {
	var (
		stockDetails = &stockDetails{
			dailyDetail: make(map[time.Time]*stockDailyDetail),
		}
	)

	stock := TheStockLibraryMaster.StockLibraries[r.Ticker]
	if stock == nil {
		return errStockNotFound
	}

	var (
		dateList timeList
	)

	for theDate := range stock.HistoryStockHoldings {
		if theDate.After(r.FromDate) && theDate.Before(r.EndDate) {
			dateList = append(dateList, theDate)
		}
	}
	if len(dateList) == 0 {
		return errNoDataInDateRange
	}

	sort.Sort(dateList)

	var (
		fundList []string
	)

	for _, fund := range allARKTypes {
		for i := 0; i < len(dateList); i++ {
			holdings := stock.HistoryStockHoldings[dateList[i]]

			holding := holdings[fund]
			if holding != nil {
				fundList = append(fundList, fund)
				break
			}
		}
	}

	glog.V(4).Infof("%s was holding in %v from %s to %s", r.Ticker, fundList, r.FromDate.Format(TheDateFormat), r.EndDate.Format(TheDateFormat))
	stockDetails.fundList = fundList
	stockDetails.dateList = dateList

	for i := 0; i < len(dateList); i++ {
		var (
			theDate = dateList[i]
		)
		stockDetails.dailyDetail[theDate] = &stockDailyDetail{
			date:     theDate,
			holdings: stock.HistoryStockHoldings[theDate],
			tradings: stock.HistoryStockTradings[theDate],
		}
	}

	r.Details = stockDetails
	return nil
}

func (r *StockReport) ToExcel() error {
	var (
		err       error
		fileName  = r.ExcelPath()
		txtReport = `关于` + r.Ticker + ": "
	)

	err = r.Load()
	if err != nil {
		return err
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

	/*
			A           B           C           D           E           F
		25		        2021-03-29	2021-03-30	2021-03-31	2021-04-01	2021-04-02
		26	ARKW持仓变动	10013448 	102200000 	10013408 	10002000 	10013448
		27	ARKF持仓变动	100230000 	10001245 	10013448 	100014000 	100014000
		28	ARK总持仓变动	110243448 	112201245 	20026856 	110016000 	110027448
	*/

	var (
		sheet = defaultSheet

		dateIdxList = []string{"B", "C", "D", "E", "F", "G",
			"H", "I", "J", "K", "L", "M", "N", "O", "P",
			"Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
		maxDate = len(dateIdxList)
	)

	for dateIdx, theDate := range r.Details.dateList {
		if dateIdx > maxDate {
			glog.Warningf("MAX DATE RANGE REACHED, we have: %d, max: %d", len(r.Details.dateList), maxDate)
			break
		}

		var (
			totalShards float64
			fundIdx     = 25
			holdings    = r.Details.dailyDetail[theDate].holdings
			//tradings    = r.Details.dailyDetail[theDate].tradings
			txtDailyReport = fmt.Sprintf("%d月%d日，", theDate.Month(), theDate.Day())
			txtDailyTemp   string
		)

		for idx, fund := range r.Details.fundList {
			holding := holdings[fund]
			var currentShards float64

			line := strconv.Itoa(fundIdx)
			// Set the date column
			if idx == 0 {
				f.SetCellValue(sheet, dateIdxList[dateIdx]+line, theDate.Format(TheDateFormat))
				fundIdx++
			}

			line = strconv.Itoa(fundIdx)
			if dateIdx == 0 {
				f.SetCellValue(sheet, "A"+line, fund)
			}
			if holding != nil {
				totalShards += holding.Shards
				currentShards = holding.Shards
			}
			f.SetCellValue(sheet, dateIdxList[dateIdx]+line, fmt.Sprintf("%.0f", currentShards))
			txtDailyTemp = txtDailyTemp + fmt.Sprintf("%s持有%.0f股(比重%.2f%%)，", holding.Fund, holding.Shards, holding.Weight)
			fundIdx++

			// Set the total
			if idx == len(r.Details.fundList)-1 {
				line = strconv.Itoa(fundIdx)
				if dateIdx == 0 {
					f.SetCellValue(sheet, "A"+line, "TOTAL")
				}
				f.SetCellValue(sheet, dateIdxList[dateIdx]+line, fmt.Sprintf("%.0f", totalShards))
				txtDailyTemp = fmt.Sprintf("ARK共持有%.0f股，", totalShards) + txtDailyTemp
			}
		}

		txtDailyReport += txtDailyTemp
		txtReport += txtDailyReport
	}

	err = f.Save()
	if err != nil {
		glog.Errorf("failed to save excel %s, err: %v", fileName, err)
		return err
	}

	glog.V(4).Infof("%s", txtReport)

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

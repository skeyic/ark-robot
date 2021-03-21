package service

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/utils"
	"sort"
	"strconv"
	"time"
)

type Top10HoldingsReport struct {
	Date            string
	holdings        *ARKHoldings
	previousHolding *ARKHoldings
}

func NewTop10HoldingsReport(date time.Time) *Top10HoldingsReport {
	var (
		r = &Top10HoldingsReport{
			Date:            date.Format(TheDateFormat),
			holdings:        TheLibrary.GetHoldings(date),
			previousHolding: TheLibrary.GetPreviousHoldings(date, 1)[0],
		}
	)

	utils.CheckFolder(r.ExcelFolder())

	return r
}

func (r *Top10HoldingsReport) ToExcel() error {
	var (
		err      error
		fileName = r.ExcelPath()
	)

	if r.holdings == nil {
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

	for _, fund := range allARKTypes {
		var (
			idx              = 35
			sheet            = fund
			toReportHoldings = make(map[float64]*StockHolding)
			toSortWeight     sort.Float64Slice
		)
		holdings := r.holdings.GetFundStockHoldings(fund)
		if holdings == nil {
			return errEmptyReport
		}
		previousHoldings := r.previousHolding.GetFundStockHoldings(fund)

		for _, holding := range holdings.Holdings {
			if toSkipTicker(holding.Ticker) {
				continue
			}
			weight := holding.Weight
			if _, hit := toReportHoldings[weight]; hit {
				weight += 0.000001
			}
			toReportHoldings[weight] = holding
			toSortWeight = append(toSortWeight, weight)
		}

		sort.Sort(sort.Reverse(toSortWeight))

		for _, weight := range toSortWeight {
			// Only the top 10
			if idx == 45 {
				break
			}

			var (
				previousWeight float64
			)

			holding := toReportHoldings[weight]
			if previousHoldings != nil && previousHoldings.Holdings != nil {
				previousHolding := previousHoldings.Holdings[holding.Ticker]
				if previousHolding != nil {
					previousWeight = previousHolding.Weight
				}
			}

			line := strconv.Itoa(idx)
			f.SetCellValue(sheet, "A"+line, holding.Ticker)
			f.SetCellValue(sheet, "B"+line, holding.Company)
			f.SetCellValue(sheet, "C"+line, previousWeight/100)
			f.SetCellValue(sheet, "D"+line, holding.Weight/100)
			f.SetCellValue(sheet, "E"+line, floatToStringIntOnly(holding.Shards))
			f.SetCellValue(sheet, "F"+line, holding.MarketValue)

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

func (r *Top10HoldingsReport) ExcelFolder() string {
	return reportPath + "/" + r.Date
}

func (r *Top10HoldingsReport) ExcelPath() string {
	return r.ExcelFolder() + "/" + r.ExcelName()
}

func (r *Top10HoldingsReport) ExcelName() string {
	return prefixTop10Holdings + r.Date + ".xlsx"
}

func (r *Top10HoldingsReport) InitExcelFromTemplate() error {
	var fileName = r.ExcelPath()
	if utils.CheckFileExist(fileName) {
		utils.DeleteFile(fileName)
	}
	utils.CopyFile(top10ExcelTemplate, fileName)
	glog.V(4).Infof("Init fileName: %s", fileName)
	return nil
}
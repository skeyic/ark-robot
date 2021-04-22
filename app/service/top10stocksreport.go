package service

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/utils"
	"sort"
	"strconv"
	"time"
)

const (
	maxIdx = 10 //Top10
)

type Top10HoldingsReport struct {
	Date            string
	holdings        *ARKHoldings
	previousHolding *ARKHoldings
	Data            []*Top10HoldingsData
}

type Top10HoldingsData struct {
	Fund string
	Data []*RankData
}

type RankData struct {
	Ticker         string
	Company        string
	PreviousWeight float64
	CurrentWeight  float64
	Shards         float64
	MarketValue    float64
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

func (r *Top10HoldingsReport) Report() error {
	var (
		err      error
		fileName = r.ExcelPath()
	)

	err = r.Load()
	if err != nil {
		return err
	}

	err = r.ToExcel()
	if err != nil {
		return err
	}

	glog.V(4).Infof("Top10StockReport %s is provided", fileName)
	return nil
}

func (r *Top10HoldingsReport) Load() error {
	if r.holdings == nil {
		glog.Warningf("Empty report")
		return errEmptyReport
	}

	for _, fund := range allARKTypes {
		var (
			idx              = 1
			toReportHoldings = make(map[float64]*StockHolding)
			toSortWeight     sort.Float64Slice
			previousHoldings *StockHoldings
			top10HoldingData = &Top10HoldingsData{Fund: fund}
		)
		holdings := r.holdings.GetFundStockHoldings(fund)
		if holdings == nil {
			return errEmptyReport
		}
		if r.previousHolding != nil {
			previousHoldings = r.previousHolding.GetFundStockHoldings(fund)
		}

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
			if idx > maxIdx {
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

			top10HoldingData.Data = append(top10HoldingData.Data, &RankData{
				Ticker:         holding.Ticker,
				Company:        holding.Company,
				PreviousWeight: previousWeight / 100,
				CurrentWeight:  holding.Weight / 100,
				Shards:         holding.Shards,
				MarketValue:    holding.MarketValue,
			})

			idx++
		}

		r.Data = append(r.Data, top10HoldingData)
	}

	return nil
}

func (r *Top10HoldingsReport) ToExcel() error {
	var (
		err      error
		fileName = r.ExcelPath()
	)

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

	for _, data := range r.Data {
		var (
			idx   = 35
			sheet = data.Fund
		)

		for _, stockData := range data.Data {
			line := strconv.Itoa(idx)
			f.SetCellValue(sheet, "A"+line, stockData.Ticker)
			f.SetCellValue(sheet, "B"+line, stockData.Company)
			f.SetCellValue(sheet, "C"+line, stockData.PreviousWeight)
			f.SetCellValue(sheet, "D"+line, stockData.CurrentWeight)
			f.SetCellValue(sheet, "E"+line, floatToStringIntOnly(stockData.Shards))
			f.SetCellValue(sheet, "F"+line, stockData.MarketValue)

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

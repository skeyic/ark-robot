package service

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/utils"
	"os"
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

	utils.CheckFolder(r.ReportFolder())

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

	err = r.ToImage()
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
				PreviousWeight: previousWeight,
				CurrentWeight:  holding.Weight,
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
			f.SetCellValue(sheet, "C"+line, stockData.PreviousWeight/100)
			f.SetCellValue(sheet, "D"+line, stockData.CurrentWeight/100)
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

func (r *Top10HoldingsReport) ToImage() error {
	for _, data := range r.Data {
		var (
			stocks           []string
			currentHoldings  = make([]opts.BarData, 0)
			previousHoldings = make([]opts.BarData, 0)
		)

		for _, stockData := range data.Data {
			stocks = append(stocks, stockData.Ticker)
			currentHoldings = append(currentHoldings, opts.BarData{
				Name:  "Current",
				Value: stockData.CurrentWeight,
				Tooltip: &opts.Tooltip{
					Show: true,
				},
			})
			previousHoldings = append(previousHoldings, opts.BarData{
				Name:  "Previous",
				Value: stockData.PreviousWeight,
				Tooltip: &opts.Tooltip{
					Show: true,
				},
			})
		}

		// create a new bar instance
		bar := charts.NewBar()

		// set some global options like Title/Legend/ToolTip or anything else
		bar.SetGlobalOptions(
			charts.WithTitleOpts(opts.Title{
				Title: data.Fund + " TOP 10（百分占比）",
				//Top: "5%",
				//Bottom: "20%",
				Left: "center",
				//Right: "20%",

			}), charts.WithLegendOpts(opts.Legend{
				Show: true,
				Top:  "7%",
			}))

		bar.SetXAxis(stocks).
			AddSeries("当前持仓", currentHoldings).
			AddSeries("昨日持仓", previousHoldings)

		var (
			htmlPath = r.htmlPath(data.Fund)
			//imagePath = r.ImagePath(data.Fund)
		)
		f, err := os.Create(htmlPath)
		if err != nil {
			glog.Errorf("failed to create html file %s", htmlPath)
			return err
		}
		err = bar.Render(f)
		if err != nil {
			glog.Errorf("failed to render html file %s", htmlPath)
			return err
		}

		// TODO do not use chrome to generate image, will add another micro service install
		//err = utils.TheChartPainter.GenerateImage(htmlPath, imagePath)
		//if err != nil {
		//	glog.Errorf("failed to save image file %s", imagePath)
		//	return err
		//}
	}

	return nil
}

func (r *Top10HoldingsReport) ReportFolder() string {
	return reportPath + "/" + r.Date
}

func (r *Top10HoldingsReport) ExcelPath() string {
	return r.ReportFolder() + "/" + r.ExcelName()
}

func (r *Top10HoldingsReport) ExcelName() string {
	return prefixTop10Holdings + r.Date + ".xlsx"
}

func (r *Top10HoldingsReport) htmlPath(fund string) string {
	return r.ReportFolder() + "/" + r.htmlName(fund)
}

func (r *Top10HoldingsReport) htmlName(fund string) string {
	return prefixTop10Holdings + r.Date + fund + ".html"
}

func (r *Top10HoldingsReport) ImagePath(fund string) string {
	return r.ReportFolder() + "/" + r.ImageName(fund)
}

func (r *Top10HoldingsReport) ImageName(fund string) string {
	return prefixTop10Holdings + r.Date + fund + ".png"
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

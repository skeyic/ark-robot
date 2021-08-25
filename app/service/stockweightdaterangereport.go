package service

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/utils"
	"os"
	"strconv"
	"time"
)

// StockWeightDateRangeReport ...
type StockWeightDateRangeReport struct {
	Tickers    []string
	FromDate   time.Time
	EndDate    time.Time
	ReportTime time.Time
	TotalDays  int

	holdings []*ARKHoldings
}

func NewStockWeightDateRangeReport(tickers []string, fromDate, endDate time.Time) *StockWeightDateRangeReport {
	r := &StockWeightDateRangeReport{
		Tickers:    tickers,
		ReportTime: time.Now(),
		FromDate:   fromDate,
		EndDate:    endDate,
	}
	utils.CheckFolder(r.ReportFolder())
	return r
}

func NewStockWeightDateRangeReportFromDays(tickers []string, days int) *StockWeightDateRangeReport {
	r := &StockWeightDateRangeReport{
		Tickers:    tickers,
		ReportTime: time.Now(),
		TotalDays:  days,
	}
	utils.CheckFolder(r.ReportFolder())
	return r
}

func (r *StockWeightDateRangeReport) ReportFolder() string {
	return stockReportPath + "/" + prefixStockWeightDateRangeReport + r.ReportTime.Format("2006-01-02-15-04-05")
}

func (r *StockWeightDateRangeReport) HtmlPath() string {
	return r.ReportFolder() + "/" + r.HtmlName()
}

func (r *StockWeightDateRangeReport) HtmlName() string {
	return r.FileName() + ".html"
}

func (r *StockWeightDateRangeReport) WeightHtmlPath() string {
	return r.ReportFolder() + "/" + r.WeightHtmlName()
}

func (r *StockWeightDateRangeReport) WeightHtmlName() string {
	return r.FileName() + "Weight.html"
}

func (r *StockWeightDateRangeReport) ImagePath() string {
	return r.ReportFolder() + "/" + r.ImageName()
}

func (r *StockWeightDateRangeReport) ImageName() string {
	return r.FileName() + ".png"
}

func (r *StockWeightDateRangeReport) FileName() string {
	return fmt.Sprintf("%s%s_from_%s_to_%s", prefixStockWeightDateRangeReport,
		r.Tickers, r.FromDate.Format(TheDateFormat), r.EndDate.Format(TheDateFormat))
}

func (r *StockWeightDateRangeReport) Load() error {
	latestDate := TheLibrary.GetLatestHoldingDate()
	if latestDate.IsZero() {
		return errNoLatestDate
	}

	currentHolding := TheLibrary.GetHoldings(latestDate)
	var (
		holdingsList []*ARKHoldings
	)

	if r.TotalDays == 0 {
		initDays := int(r.EndDate.Sub(r.FromDate).Hours() / 24)
		tempHoldingsList := TheLibrary.GetPreviousHoldingsWithoutNil(latestDate, initDays)
		for _, theHolding := range tempHoldingsList {
			if !theHolding.Date.Before(r.FromDate) {
				holdingsList = append(holdingsList, theHolding)
			}
		}
		r.TotalDays = len(holdingsList)
	} else {
		holdingsList = TheLibrary.GetPreviousHoldingsWithoutNil(latestDate, r.TotalDays)
		r.EndDate = currentHolding.Date
		r.FromDate = holdingsList[0].Date
	}

	r.holdings = append(holdingsList, currentHolding)

	return nil
}

func (r *StockWeightDateRangeReport) ReportWeightImage() error {
	var (
		theX  []string
		TheYs = make(map[string][]opts.LineData)
	)

	tr := charts.NewLine()
	tr.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "ARK持仓中的中概股具体个股比重变化情况",
			Left:  "center",
			Top:   "2%",
		}),
		charts.WithSingleAxisOpts(opts.SingleAxis{
			Type:   "time",
			Bottom: "10%",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: true,
			Top:  "6%",
			Left: "center",
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Theme:  types.ThemeEssos,
			Width:  "1200px",
			Height: "800px",
		}),
		//charts.WithColorsOpts(utils.ColorsForPainter),
	)

	for _, theHoldings := range r.holdings {
		var (
			theMarketValues       = make(map[string]float64)
			totalStockMarketValue float64
		)
		for _, theFund := range allARKTypes {
			holdings := theHoldings.GetFundStockHoldings(theFund)
			if holdings == nil {
				continue
			}
			for _, theTicker := range r.Tickers {
				theHolding := holdings.Holdings[theTicker]
				if theHolding != nil {
					theMarketValues[theTicker] += theHolding.MarketValue
				}
			}
			for _, holding := range holdings.Holdings {
				totalStockMarketValue += holding.MarketValue
			}
		}

		theX = append(theX, theHoldings.Date.Format(TheDatePainterFormat))
		for _, theTicker := range r.Tickers {
			tempValue := fmt.Sprintf("%.6f", theMarketValues[theTicker]/totalStockMarketValue)
			value, _ := strconv.ParseFloat(tempValue, 10)
			TheYs[theTicker] = append(TheYs[theTicker], opts.LineData{Value: value})
		}
	}

	tr.SetXAxis(theX).SetSeriesOptions(
		charts.WithLabelOpts(opts.Label{Show: true}))

	for ticker, theY := range TheYs {
		tr.AddSeries(TheChinaStockManager.Translate(ticker), theY)
	}

	var (
		htmlPath = r.WeightHtmlPath()
		//imagePath = r.ImagePath()
	)
	f, err := os.Create(htmlPath)
	if err != nil {
		glog.Errorf("failed to create html file %s", htmlPath)
		return err
	}
	err = tr.Render(f)
	if err != nil {
		glog.Errorf("failed to render html file %s", htmlPath)
		return err
	}
	//}
	//
	//// TODO do not use chrome to generate image, will add another micro service install
	//err = utils.TheChartPainter.GenerateImage(htmlPath, imagePath)
	//if err != nil {
	//	glog.Errorf("failed to save image file %s", imagePath)
	//	return err
	//}

	return nil
}

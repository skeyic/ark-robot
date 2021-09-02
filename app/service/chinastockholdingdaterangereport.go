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

// ChinaStockHoldingDateRangeReport ...
type ChinaStockHoldingDateRangeReport struct {
	FromDate   time.Time
	EndDate    time.Time
	ReportTime time.Time
	TotalDays  int

	holdings []*ARKHoldings
}

func NewChinaStockHoldingDateRangeReport(fromDate, endDate time.Time) *ChinaStockHoldingDateRangeReport {
	r := &ChinaStockHoldingDateRangeReport{
		ReportTime: time.Now(),
		FromDate:   fromDate,
		EndDate:    endDate,
	}
	utils.CheckFolder(r.ReportFolder())
	return r
}

func NewChinaStockHoldingDateRangeReportFromDays(days int) *ChinaStockHoldingDateRangeReport {
	r := &ChinaStockHoldingDateRangeReport{
		ReportTime: time.Now(),
		TotalDays:  days,
	}
	utils.CheckFolder(r.ReportFolder())
	return r
}

func (r *ChinaStockHoldingDateRangeReport) ReportFolder() string {
	return stockReportPath + "/" + prefixChinaStockDateRangeReport + r.ReportTime.Format("2006-01-02-15-04-05")
}

func (r *ChinaStockHoldingDateRangeReport) HtmlPath() string {
	return r.ReportFolder() + "/" + r.HtmlName()
}

func (r *ChinaStockHoldingDateRangeReport) HtmlName() string {
	return r.FileName() + ".html"
}

func (r *ChinaStockHoldingDateRangeReport) WeightHtmlPath() string {
	return r.ReportFolder() + "/" + r.WeightHtmlName()
}

func (r *ChinaStockHoldingDateRangeReport) WeightHtmlName() string {
	return r.FileName() + "Weight.html"
}

func (r *ChinaStockHoldingDateRangeReport) ImagePath() string {
	return r.ReportFolder() + "/" + r.ImageName()
}

func (r *ChinaStockHoldingDateRangeReport) ImageName() string {
	return r.FileName() + ".png"
}

func (r *ChinaStockHoldingDateRangeReport) FileName() string {
	return fmt.Sprintf("%s_from_%s_to_%s", prefixChinaStockDateRangeReport,
		r.FromDate.Format(TheDateFormat), r.EndDate.Format(TheDateFormat))
}

func (r *ChinaStockHoldingDateRangeReport) Load() error {
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

func (r *ChinaStockHoldingDateRangeReport) ReportImage() error {
	var (
		trData []opts.ThemeRiverData
	)

	tr := charts.NewThemeRiver()
	tr.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "ARK持仓中的中概股市值变化情况",
			Left:  "center",
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
			//Theme:  types.ThemeEssos,
			Width:  "1200px",
			Height: "800px",
		}),
		charts.WithColorsOpts(utils.ColorsForPainter),
	)

	for _, theHoldings := range r.holdings {
		var (
			chinaStockMap = make(map[string]float64)
		)
		for _, theFund := range allARKTypes {
			holdings := theHoldings.GetFundStockHoldings(theFund)
			if holdings == nil {
				continue
			}
			for ticker, holding := range holdings.Holdings {
				if TheChinaStockManager.IsChinaStock(ticker) {
					previous, _ := chinaStockMap[ticker]
					chinaStockMap[ticker] = previous + holding.MarketValue
				}
			}
		}
		for ticker, marketValue := range chinaStockMap {
			trData = append(trData, opts.ThemeRiverData{
				Date:  theHoldings.Date.Format(TheDatePainterFormat),
				Value: marketValue,
				Name:  TheChinaStockManager.Translate(ticker),
			})
		}
	}
	tr.AddSeries("themeRiver", trData, charts.WithLabelOpts(opts.Label{
		//Position: "inside",
		Show: false,
	}))

	//for _, t := range []string{types.ThemeEssos, types.ThemeChalk, types.ThemeRoma, types.ThemeRomantic, types.ThemeInfographic, types.ThemeMacarons, types.ThemePurplePassion, types.ThemeShine, types.ThemeVintage, types.ThemeWalden, types.ThemeWesteros, types.ThemeWonderland} {
	var (
		htmlPath = r.HtmlPath()
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

func (r *ChinaStockHoldingDateRangeReport) ReportWeightImage() error {
	var (
		theX []string
		theY []opts.LineData
	)

	tr := charts.NewLine()
	tr.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "ARK持仓中的中概股比重变化情况",
			Left:  "center",
		}),
		charts.WithSingleAxisOpts(opts.SingleAxis{
			Type:   "time",
			Bottom: "10%",
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Theme:  types.ThemeEssos,
			Width:  "1200px",
			Height: "800px",
		}),
		charts.WithColorsOpts(utils.ColorsForPainter),
	)

	for _, theHoldings := range r.holdings {
		var (
			chinaStockMarketValue float64
			totalStockMarketValue float64
		)
		for _, theFund := range allARKTypes {
			holdings := theHoldings.GetFundStockHoldings(theFund)
			if holdings == nil {
				continue
			}
			for ticker, holding := range holdings.Holdings {
				if TheChinaStockManager.IsChinaStock(ticker) {
					chinaStockMarketValue += holding.MarketValue
				}
				totalStockMarketValue += holding.MarketValue
			}
		}

		theX = append(theX, theHoldings.Date.Format(TheDatePainterFormat))
		tempValue := fmt.Sprintf("%.3f", chinaStockMarketValue/totalStockMarketValue)
		value, _ := strconv.ParseFloat(tempValue, 10)
		theY = append(theY, opts.LineData{Value: value})
	}

	tr.SetXAxis(theX).AddSeries("Weight", theY).SetSeriesOptions(
		charts.WithLabelOpts(opts.Label{Show: true}))

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

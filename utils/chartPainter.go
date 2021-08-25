package utils

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/golang/glog"
	"io/ioutil"
)

var (
	TheChartPainter  = &ChartPainter{}
	ColorsForPainter = opts.Colors{
		"#5470c6", "#91cc75", "#fac858", "#ee6666", "#73c0de", "#3ba272", "#fc8452", "#9a60b4", "#ea7ccc", "#ef4464",
		"#929fff", "#e75840", "#50c48f", "#26ccd8", "#3685fe", "#9977ef", "#FFB6C1", "#2f4554", "#61a0a8", "#d48265",
		"#91c7ae", "#749f83", "#ca8622", "#bda29a", "#6e7074", "#546570", "#c4ccd3",
	}
)

type ChartPainter struct {
}

func (c *ChartPainter) GenerateImage(htmlPath, imagePath string) error {
	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		//chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	var buf []byte

	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot("file://"+htmlPath, 90, &buf)); err != nil {
		glog.Errorf("failed to take snapshot, html: %s, err: %v", htmlPath, err)
		return err
	}

	if CheckFileExist(imagePath) {
		DeleteFile(imagePath)
	}

	// save image
	if err := ioutil.WriteFile(imagePath, buf, 0x644); err != nil {
		glog.Errorf("failed to save image, image: %s, err: %v", imagePath, err)
		return err
	}

	return nil
}

func ToBarData(name string, data []float64) []opts.BarData {
	var (
		theData   []opts.BarData
		showLabel = len(data) <= 10
	)
	for _, value := range data {
		theData = append(theData,
			opts.BarData{
				Name:  name,
				Value: value,
				Label: &opts.Label{
					Show:     showLabel,
					Position: "insideTop",
				},
				//ItemStyle: nil,
				Tooltip: &opts.Tooltip{
					Show: true,
				},
			})
	}
	return theData
}

func ToLineData(name string, data []float64) []opts.LineData {
	var (
		theData []opts.LineData
	)
	for _, value := range data {
		theData = append(theData,
			opts.LineData{
				Name:  name,
				Value: value,
			})
	}
	return theData
}

func ToPercentLineData(name string, data []float64, percent float64) []opts.LineData {
	var (
		theData []opts.LineData
	)
	for _, value := range data {
		theData = append(theData,
			opts.LineData{
				Name:  name,
				Value: value / percent,
			})
	}
	return theData
}

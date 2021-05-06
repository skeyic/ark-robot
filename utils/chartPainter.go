package utils

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/golang/glog"
	"io/ioutil"
)

var (
	TheChartPainter = &ChartPainter{}
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
		theData []opts.BarData
	)
	for _, value := range data {
		theData = append(theData,
			opts.BarData{
				Name:  name,
				Value: value,
				Label: &opts.Label{
					Show: true,
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

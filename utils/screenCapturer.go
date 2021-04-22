package utils

import (
	"context"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/golang/glog"
	"io/ioutil"
	"math"
)

var (
	TheScreenCapture = &ScreenCapture{}
)

const (
	imageURL = "http://localhost:8081/"
	//imageURL = "https://www.google.com"
	//selector = `#main`
)

type ScreenCapture struct {
}

func (s *ScreenCapture) GenerateImage(imagePath string) {
	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		//chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	var buf []byte
	glog.V(4).Info("START CAPTURE")
	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot(imageURL, 90, &buf)); err != nil {
		glog.Errorf("failed to run full screenshot, url: %s, err: %v", imageURL, err)
		panic(err)
	}
	glog.V(4).Info("START SAVE")
	if err := ioutil.WriteFile(imagePath, buf, 0o644); err != nil {
		glog.Errorf("failed to save full screenshot %s, err: %v", imagePath, err)
		panic(err)
	}
}

func fullScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, cssContentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(cssContentSize.Width)), int64(math.Ceil(cssContentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      cssContentSize.X,
					Y:      cssContentSize.Y,
					Width:  cssContentSize.Width,
					Height: cssContentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}

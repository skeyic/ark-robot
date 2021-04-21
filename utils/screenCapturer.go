package utils

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/golang/glog"
	"io/ioutil"
	"log"
)

var (
	TheScreenCapture = &ScreenCapture{}
)

const (
	//imageURL = "http://localhost:8081/"
	imageURL = "https://www.google.com"
	selector = `#main`
)

type ScreenCapture struct {
}

func (s *ScreenCapture) GenerateImage(imgPath string) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	glog.V(4).Info("NAVIGATE URL BEFORE")
	chromedp.Navigate("http://www.baidu.com")

	glog.V(4).Info("NAVIGATE URL")
	if err := chromedp.Run(ctx, elementScreenshot(imageURL, selector, &buf)); err != nil {
		log.Fatal(err)
	}
	glog.V(4).Info("AFTER RUN")
	// 写入文件
	if err := ioutil.WriteFile(imgPath, buf, 0644); err != nil {
		log.Fatal(err)
	}
}

func elementScreenshot(url, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		// 打开url指向的页面
		chromedp.Navigate(url),

		//// 等待待截图的元素渲染完成
		//chromedp.WaitVisible(sel, chromedp.ByID),

		// 也可以等待一定的时间
		//chromedp.Sleep(time.Duration(3) * time.Second),

		// 执行截图
		chromedp.Screenshot(sel, res, chromedp.NodeVisible, chromedp.ByID),
	}
}

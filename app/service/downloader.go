package service

import (
	"bytes"
	"fmt"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/config"
	"github.com/skeyic/ark-robot/utils"
	"github.com/tebeka/selenium"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	csvBaseURL = "https://ark-funds.com/wp-content/uploads/funds-etf-csv/%s.csv"
)

var (
	arkMap = map[string]string{
		"ARKK": "ARK_INNOVATION_ETF_ARKK_HOLDINGS",
		"ARKQ": "ARK_AUTONOMOUS_TECHNOLOGY_&_ROBOTICS_ETF_ARKQ_HOLDINGS",
		"ARKW": "ARK_NEXT_GENERATION_INTERNET_ETF_ARKW_HOLDINGS",
		"ARKG": "ARK_GENOMIC_REVOLUTION_MULTISECTOR_ETF_ARKG_HOLDINGS",
		"ARKF": "ARK_FINTECH_INNOVATION_ETF_ARKF_HOLDINGS",
		"ARKX": "ARK_SPACE_EXPLORATION_&_INNOVATION_ETF_ARKX_HOLDINGS", // added 2021/3/30
	}
	downloaderFolder        = config.Config.DataFolder + "/downloader/"
	downloaderUTCStartHour  = 0 // UTC 00:00
	downloaderCheckInterval = 4 * time.Hour
	downloaderRetryWaitTime = 5 * time.Minute
)

var (
	TheDownloader = NewDownloader()
)

func generateArkCSVURL(arkType string) string {
	arkName, hit := arkMap[arkType]
	if !hit {
		glog.Fatalf("incorrect ark type: %s", arkType)
	}

	return fmt.Sprintf(csvBaseURL, arkName)
}

func generateDownloaderFilePath(arkType string) string {
	return config.Config.DataFolder + "/downloader/" + time.Now().Format("2006-01-02-15-04-05-") + arkType + ".csv"
}

type Downloader struct {
}

func NewDownloader() *Downloader {
	d := &Downloader{}
	return d
}

func (d *Downloader) Init() error {
	utils.CheckFolder(downloaderFolder)
	glog.V(4).Infof("downloader init completed")
	return nil
}

func (d *Downloader) process() {

}

func screenshot(driver selenium.WebDriver, filename string) {
	bytes, err := driver.Screenshot()
	if nil != err {
		glog.Errorf("take screenshot error, err: %s", err)
		return
	}

	err = ioutil.WriteFile(config.Config.DataFolder+"/"+filename, bytes, 0666)
	if nil != err {
		glog.Errorf("save screenshot error, err: %s", err)
		return
	}
}

func screenshotV2(driver selenium.WebDriver) {
	pic := fmt.Sprintf("%s.jpg", time.Now().Format("15-04-05"))
	screenshot(driver, pic)
	glog.V(4).Infof("TAKE PICTURE %s", pic)
}

func (d *Downloader) DownloadAllARKCSVsV2() error {
	var (
		fileNames   []string
		arkHoldings = &ARKHoldings{}
	)

	var (
		gridURL     = config.Config.SpiderServer.URL
		browserName = "chrome"
		URL         = "https://ark-funds.com/download-fund-materials/"
		driver      selenium.WebDriver
		fundNameMap = map[string]string{
			"ARK Innovation ETF":                     "ARK_INNOVATION_ETF_ARKK_HOLDINGS",
			"ARK Genomic Revolution ETF":             "ARK_GENOMIC_REVOLUTION_ETF_ARKG_HOLDINGS",
			"ARK Next Generation Internet ETF":       "ARK_NEXT_GENERATION_INTERNET_ETF_ARKW_HOLDINGS",
			"ARK Autonomous Tech. & Robotics ETF":    "ARK_AUTONOMOUS_TECH._&_ROBOTICS_ETF_ARKQ_HOLDINGS",
			"ARK Fintech Innovation ETF":             "ARK_FINTECH_INNOVATION_ETF_ARKF_HOLDINGS",
			"ARK Space Exploration & Innovation ETF": "ARK_SPACE_EXPLORATION_&_INNOVATION_ETF_ARKX_HOLDINGS",
		}
	)

	caps := selenium.Capabilities{"browserName": browserName}
	webDriver, err := selenium.NewRemote(caps, gridURL)
	if nil != err {
		panic(err)
	}
	driver = webDriver

	// teardown
	defer driver.Quit()

	err = driver.Get(URL)
	if nil != err {
		glog.Errorf("page open error, err: %s", err)
		return errDownloadCSV
	}

	//screenshot(driver, "1.jpg")
	err = driver.ResizeWindow("", 3200, 2600)
	if err != nil {
		glog.Errorf("Failed to resize window, err: %+v", err)
		return errDownloadCSV
	}

	//searchDoc, err := driver.FindElement(selenium.ByID, "searchDoc")
	//if err != nil {
	//	screenshotV2(driver)
	//	glog.Errorf("Failed to find search, err: %+v", err)
	//	return errDownloadCSV
	//}
	//err = searchDoc.Click()
	//if err != nil {
	//	screenshotV2(driver)
	//	glog.Errorf("Failed to click search, err: %+v", err)
	//	return errDownloadCSV
	//}
	//err = searchDoc.SendKeys("csv")
	//if err != nil {
	//	screenshotV2(driver)
	//	glog.Errorf("Failed to search, err: %+v", err)
	//	return errDownloadCSV
	//}
	//err = searchDoc.SendKeys(selenium.EnterKey)
	//if err != nil {
	//	screenshotV2(driver)
	//	glog.Errorf("Failed to search, err: %+v", err)
	//	return errDownloadCSV
	//}
	////screenshot(driver, "2.jpg")
	//driver.WaitWithTimeoutAndInterval(func(driver selenium.WebDriver) (bool, error) {
	//	_, err := driver.FindElement(selenium.ByLinkText, "Fund Holdings CSV")
	//	if err != nil {
	//		screenshot(driver, fmt.Sprintf("%s.jpg", time.Now().Format("15-04-05")))
	//	}
	//	return err == nil, nil
	//}, 100*time.Second, 10*time.Second)
	//
	screenshotV2(driver)
	navigate, err := driver.FindElement(selenium.ByLinkText, "Fund Holdings CSV")
	if err != nil {
		glog.Errorf("Failed to find element Fund Holdings CSV, err: %+v", err)
		screenshotV2(driver)
		return errDownloadCSV
	}

	err = navigate.Click()
	if err != nil {
		glog.Errorf("Failed to click element Fund Holdings CSV, err: %+v", err)
		screenshotV2(driver)
		return errDownloadCSV
	}

	for fileType, fileName := range fundNameMap {
		for i := 0; i < 3; i++ {
			driver.WaitWithTimeoutAndInterval(func(driver selenium.WebDriver) (bool, error) {
				_, err := driver.FindElement(selenium.ByXPATH, "//div[contains(text(),'"+fileType+"')]/../../../div[2]//button")
				if err != nil {
					screenshotV2(driver)
				}
				return err == nil, nil
			}, 100*time.Second, 10*time.Second)

			e, err := driver.FindElement(selenium.ByXPATH, "//div[contains(text(),'"+fileType+"')]/../../../div[2]//button")
			if err != nil {
				glog.Errorf("Failed to find element, err: %+v", err)
				if i == 2 {
					return errDownloadCSV
				}
				time.Sleep(5 * time.Second)
				continue
			}

			err = e.Click()
			if err != nil {
				glog.Errorf("Failed to click element, err: %+v", err)
				if i == 2 {
					return errDownloadCSV
				}
				time.Sleep(5 * time.Second)
				continue
			}
			break
		}

		time.Sleep(5 * time.Second)

		latestFileName, err := GetLatestFileName(config.Config.SpiderServer.DataFolder, fileName, ".csv")
		if err != nil {
			glog.Errorf("failed to get latest file name, err: %v", err)
			return errDownloadCSV
		}
		fileNames = append(fileNames, latestFileName)
		glog.V(4).Infof("Downloaded %s", latestFileName)
	}

	// Wait 30 seconds to make sure the download is finished
	time.Sleep(30 * time.Second)

	defer func(files []string) {
		for _, file := range files {
			glog.V(4).Infof("DELETE FILE %s", file)
			os.Remove(file)
		}
	}(fileNames)

	for _, fileName := range fileNames {
		stockHoldings, err := ThePorter.Catalog(fileName)
		if err != nil {
			return err
		}

		err = arkHoldings.AddStockHoldings(stockHoldings)
		if err != nil {
			return err
		}
	}

	if !arkHoldings.Validation() {
		return errValidateFail
	}

	if TheLibrary.GetLatestHoldingDate() == arkHoldings.Date {
		glog.V(4).Infof("No need to update, latest date: %s", arkHoldings.Date)
		return nil
	}

	TheLibrary.GenerateCurrentTrading(arkHoldings)
	TheLibrary.AddStockHoldings(arkHoldings)
	TheStockLibraryMaster.AddStockHoldings(arkHoldings)
	TheTop10HoldingsReportMaster.Refresh()

	glog.V(4).Infof("Add ark holdings of %s at %s to library", arkHoldings.Date, time.Now())
	utils.SendAlertV2("Add to library", fmt.Sprintf("Add ark holdings of %s at %s to library", arkHoldings.Date, time.Now()))

	err = TheMaster.ReportLatestTrading(true)
	if err != nil {
		glog.Errorf("report latest trading failed, err: %v", err)
		return err
	}

	if config.Config.ESServer.Force {
		err = TheMaster.IndexLatestToES()
		if err != nil {
			glog.Errorf("index latest data to ES failed, err: %v", err)
			return err
		}
	}

	glog.V(4).Infof("TradingsReport latest trading of %s at %s to library", arkHoldings.Date, time.Now())
	//if config.Config.DebugMode {
	//	utils.SendAlertV2("Add to library", fmt.Sprintf("TradingsReport latest trading of %s at %s to library", arkHoldings.Date, time.Now()))
	//}

	return nil
}

func GetLatestFileName(dirPath, prefix, suffix string) (string, error) {
	var (
		filesets        []string
		latestTimestamp int64
		latestFile      string
	)

	Listfunc := func(path string, f os.FileInfo, err error) error {
		//ostype := os.Getenv("GOOS") // windows, linux

		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}

		if prefix != "" {
			ok := strings.HasPrefix(f.Name(), prefix)
			if !ok {
				return nil
			}
		}

		if suffix != "" {
			ok := strings.HasSuffix(f.Name(), suffix)
			if !ok {
				return nil
			}
		}

		filesets = append(filesets, path)

		return nil
	}

	err := filepath.Walk(dirPath, Listfunc)
	if err != nil {
		glog.Errorf("failed to walk dir: %s, err: %v", dirPath, err)
		return "", errDownloadCSV
	}

	glog.V(4).Infof("The FILE LIST: %v", filesets)

	if len(filesets) == 1 {
		return filesets[0], nil
	}

	for _, theFile := range filesets {
		theList := strings.Split(theFile, "_")
		if len(theList) == 6 {
			theTS := strings.TrimSuffix(theList[5], suffix)
			theTSInt, _ := strconv.ParseInt(theTS, 10, 64)
			if theTSInt > latestTimestamp {
				latestTimestamp = theTSInt
				latestFile = theFile
			}
		}
	}

	return latestFile, nil
}

func (d *Downloader) DownloadAllARKCSVs() error {
	var (
		fileNames   []string
		arkHoldings = &ARKHoldings{}
	)

	// Make sure we have downloaded all funds
	for _, theType := range allARKTypes {
		theFile, err := d.DownloadARKCSV(theType)
		if err != nil {
			glog.Errorf("download ARK %s CSV failed, err: %v", theType, err)
			return err
		}
		fileNames = append(fileNames, theFile)
	}

	defer func(files []string) {
		for _, file := range files {
			glog.V(4).Infof("DELETE FILE %s", file)
			//os.Remove(file)
		}
	}(fileNames)

	for _, fileName := range fileNames {
		stockHoldings, err := ThePorter.Catalog(fileName)
		if err != nil {
			return err
		}

		err = arkHoldings.AddStockHoldings(stockHoldings)
		if err != nil {
			return err
		}
	}

	if !arkHoldings.Validation() {
		return errValidateFail
	}

	if TheLibrary.GetLatestHoldingDate() == arkHoldings.Date {
		glog.V(4).Infof("No need to update, latest date: %s", arkHoldings.Date)
		return nil
	}

	TheLibrary.GenerateCurrentTrading(arkHoldings)
	TheLibrary.AddStockHoldings(arkHoldings)
	TheStockLibraryMaster.AddStockHoldings(arkHoldings)
	TheTop10HoldingsReportMaster.Refresh()

	glog.V(4).Infof("Add ark holdings of %s at %s to library", arkHoldings.Date, time.Now())
	utils.SendAlertV2("Add to library", fmt.Sprintf("Add ark holdings of %s at %s to library", arkHoldings.Date, time.Now()))

	err := TheMaster.ReportLatestTrading(true)
	if err != nil {
		glog.Errorf("report latest trading failed, err: %v", err)
		return err
	}

	if config.Config.ESServer.Force {
		err = TheMaster.IndexLatestToES()
		if err != nil {
			glog.Errorf("index latest data to ES failed, err: %v", err)
			return err
		}
	}

	glog.V(4).Infof("TradingsReport latest trading of %s at %s to library", arkHoldings.Date, time.Now())
	//if config.Config.DebugMode {
	//	utils.SendAlertV2("Add to library", fmt.Sprintf("TradingsReport latest trading of %s at %s to library", arkHoldings.Date, time.Now()))
	//}

	return nil
}

func (d *Downloader) DownloadARKCSV(arkType string) (string, error) {
	var (
		url      = generateArkCSVURL(arkType)
		fileName = generateDownloaderFilePath(arkType)
	)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authority", "ark-funds.com")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", "https://ark-funds.com/investor-resources")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Cookie", "_ga=GA1.2.1970799815.1612344474; __cfduid=dd1ea544054408d2ddd9a60fe5981e7191615726396; PHPSESSID=ihegc2qttn6rg1oifupl91kmkl; _gid=GA1.2.1642418799.1615726420; _gat=1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		glog.Errorf("download CSV failed, url: %s, err: %v", url, err)
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	offset := bytes.Index(body, []byte(",,,,,,,"))
	if offset == -1 {
		glog.Errorf("incorrect csv format")
		return "", errDownloadCSV
	}

	err = ioutil.WriteFile(fileName, body[:offset], os.ModePerm)
	if err != nil {
		glog.Errorf("copy resp data to file failed, err: %v", err)
		return "", err
	}

	glog.V(4).Infof("download CSV %s completed", fileName)

	return fileName, nil
}

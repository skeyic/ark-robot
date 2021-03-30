package service

import (
	"bytes"
	"fmt"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/config"
	"github.com/skeyic/ark-robot/utils"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	csvBaseURL = "https://ark-funds.com/wp-content/fundsiteliterature/csv/%s.csv"
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

	glog.V(4).Infof("Add ark holdings of %s at %s to library", arkHoldings.Date, time.Now())
	if config.Config.DebugMode {
		utils.SendAlertV2("Add to library", fmt.Sprintf("Add ark holdings of %s at %s to library", arkHoldings.Date, time.Now()))
	}

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
	if config.Config.DebugMode {
		utils.SendAlertV2("Add to library", fmt.Sprintf("TradingsReport latest trading of %s at %s to library", arkHoldings.Date, time.Now()))
	}

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

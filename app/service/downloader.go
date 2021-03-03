package service

import (
	"bytes"
	"errors"
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
	}
	downloaderFolder        = config.Config.DataFolder + "/downloader/"
	downloaderUTCStartHour  = 0 // UTC 00:00
	downloaderCheckInterval = 4 * time.Hour
	downloaderRetryWaitTime = 5 * time.Minute
)

var (
	errDownloadCSV      = errors.New("download csv failed")
	errFileAlreadyExist = errors.New("file already exist")
	errRenameFile       = errors.New("rename file failed")
	errDateNotMatch     = errors.New("date not match")
	errFundNotMatch     = errors.New("fund not match")
	errValidateFail     = errors.New("validate failed")
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
	d.init()
	return d
}

func (d *Downloader) init() {
	utils.CheckFolder(downloaderFolder)
	glog.V(4).Infof("downloader init completed")
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

	TheLibrary.AddStockHoldings(arkHoldings)
	TheStockLibraryMaster.AddStockHoldings(arkHoldings)

	glog.V(4).Infof("Add ark holdings of %s at %s to library", arkHoldings.Date, time.Now())
	if config.Config.DebugMode {
		utils.SendAlertV2("Add to library", fmt.Sprintf("Add ark holdings of %s at %s to library", arkHoldings.Date, time.Now()))
	}

	return nil
}

func (d *Downloader) DownloadARKCSV(arkType string) (string, error) {
	var (
		url      = generateArkCSVURL(arkType)
		fileName = generateDownloaderFilePath(arkType)
	)
	resp, err := http.Get(url)
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

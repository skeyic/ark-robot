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
	allARKTypes = []string{"ARKK", "ARKQ", "ARKW", "ARKG", "ARKF"}
	arkMap      = map[string]string{
		"ARKK": "ARK_INNOVATION_ETF_ARKK_HOLDINGS",
		"ARKQ": "ARK_AUTONOMOUS_TECHNOLOGY_&_ROBOTICS_ETF_ARKQ_HOLDINGS",
		"ARKW": "ARK_NEXT_GENERATION_INTERNET_ETF_ARKW_HOLDINGS",
		"ARKG": "ARK_GENOMIC_REVOLUTION_MULTISECTOR_ETF_ARKG_HOLDINGS",
		"ARKF": "ARK_FINTECH_INNOVATION_ETF_ARKF_HOLDINGS",
	}
	downloaderFolder = config.Config.DataFolder + "/downloader/"
)

var (
	errDownloadCSV = errors.New("download csv failed")
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

func (d *Downloader) DownloadAllARKCSVs() error {
	for _, theType := range allARKTypes {
		err := d.DownloadARKCSV(theType)
		if err != nil {
			glog.Errorf("download ARK %s CSV failed, err: %v", theType, err)
			return err
		}
	}
	return nil
}

func (d *Downloader) DownloadARKCSV(arkType string) error {
	var (
		url      = generateArkCSVURL(arkType)
		filename = generateDownloaderFilePath(arkType)
	)
	resp, err := http.Get(url)
	if err != nil {
		glog.Errorf("download CSV failed, url: %s, err: %v", url, err)
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	offset := bytes.Index(body, []byte(",,,,,,,"))
	if offset == -1 {
		glog.Errorf("incorrect csv format")
		return errDownloadCSV
	}

	err = ioutil.WriteFile(filename, body[:offset], os.ModePerm)
	if err != nil {
		glog.Errorf("copy resp data to file failed, err: %v", err)
		return err
	}

	glog.V(4).Infof("download CSV %s completed", filename)
	ThePorter.Catalog(filename)

	return nil
}

package service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/app/config"
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
	downloaderFolder = config.Config.DataFolder + "/downloader/"
)

var (
	errDownloadCSV = errors.New("download csv failed")
)

func init() {
	_, err := os.Stat(downloaderFolder)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(downloaderFolder, 0777)
			if err != nil {
				panic(err)
			}
		}
	}
}

func generateArkCSVURL(arkType string) string {
	arkName, hit := arkMap[arkType]
	if !hit {
		glog.Fatalf("Incorrect ark type: %s", arkType)
	}

	return fmt.Sprintf(csvBaseURL, arkName)
}

func generateFilePath(arkType string) string {
	return config.Config.DataFolder + "/downloader/" + time.Now().Format("20060102") + arkType + ".csv"
}

func DownloadARKCSV(url string, filename string) error {
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

	go ThePorter.Catalog(filename)

	return nil
}

type Downloader struct {
}

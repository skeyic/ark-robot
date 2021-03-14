package service

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/utils"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestDownloadCSV(t *testing.T) {
	flag.Set("logtostderr", "true")
	flag.Set("v", "10")
	flag.Parse()

	err := TheDownloader.DownloadAllARKCSVs()
	if err != nil {
		glog.Errorf("failed to download csv, err: %v", err)
		return
	}

	//<-make(chan struct{}, 1)
}

func TestDownloadTime(t *testing.T) {
	fmt.Printf("NOW: %d", time.Now().UTC().Day())
}

func Test_DownloadWithMoreInfo(t *testing.T) {
	utils.EnableGlogForTesting()

	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	// curl "https://ark-funds.com/wp-content/fundsiteliterature/csv/ARK_NEXT_GENERATION_INTERNET_ETF_ARKW_HOLDINGS.csv" ^
	//   -H "authority: ark-funds.com" ^
	//   -H "upgrade-insecure-requests: 1" ^
	//   -H "user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36" ^
	//   -H "accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9" ^
	//   -H "sec-fetch-site: same-origin" ^
	//   -H "sec-fetch-mode: navigate" ^
	//   -H "sec-fetch-user: ?1" ^
	//   -H "sec-fetch-dest: document" ^
	//   -H "referer: https://ark-funds.com/investor-resources" ^
	//   -H "accept-language: zh-CN,zh;q=0.9,en;q=0.8" ^
	//   -H "cookie: _ga=GA1.2.1970799815.1612344474; __cfduid=dd1ea544054408d2ddd9a60fe5981e7191615726396; PHPSESSID=ihegc2qttn6rg1oifupl91kmkl; _gid=GA1.2.1642418799.1615726420; _gat=1" ^
	//   --compressed

	req, err := http.NewRequest("GET", "https://ark-funds.com/wp-content/fundsiteliterature/csv/ARK_NEXT_GENERATION_INTERNET_ETF_ARKW_HOLDINGS.csv", nil)
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
		// handle err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("ERR: %v", err)
	}
	glog.V(4).Infof("BODY: %s", b)
}

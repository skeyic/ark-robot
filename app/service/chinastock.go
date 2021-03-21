package service

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/utils"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	theChinaStockManagerFileStore = libraryFolder + "TheChinaStockManager"
	TheChinaStockManager          = &ChinaStockManager{
		fileStore: utils.NewFileStoreSvc(theChinaStockManagerFileStore),
		stocks:    make(map[string]*ChinaStock),
	}
)

type ChinaStockManager struct {
	fileStore *utils.FileStoreSvc
	stocks    map[string]*ChinaStock
}

func (m *ChinaStockManager) IsChinaStock(ticker string) bool {
	glog.V(4).Infof("TICKER: %s, STOCK: %+v", ticker, m.stocks[ticker])
	return m.stocks[ticker] != nil
}

type ChinaStock struct {
	Ticker string `json:"symbol"`
	Name   string `json:"name"`
}

func (m *ChinaStockManager) Init() error {
	var (
		err error
	)
	err = m.LoadFromFileStore()
	if err != nil {
		return err
	}
	glog.V(4).Infof("ChinaStockManager LoadFromFileStore num: %d", len(m.stocks))
	if len(m.stocks) == 0 {
		err = m.LoadAllChinaStock()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *ChinaStockManager) LoadFromFileStore() error {
	theBytes, err := m.fileStore.Read()
	if err != nil {
		if os.IsNotExist(err) {
			glog.V(4).Info("No saved file for library")
			return nil
		}
		glog.Errorf("failed to load china stock manager from the saved file")
		return err
	}

	err = json.Unmarshal(theBytes, &m)
	if err != nil {
		glog.Errorf("failed to unmarshal the saved file to library")
		return err
	}

	return nil
}

func (m *ChinaStockManager) Save() error {
	uByte, err := json.Marshal(m)
	if err != nil {
		glog.Errorf("failed to marshal the china stock manager, err: %v", err)
		return err
	}
	err = m.fileStore.Save(uByte)
	if err != nil {
		glog.Errorf("failed to save the china stock manager, err: %v", err)
		return err
	}
	return nil
}

func (m *ChinaStockManager) MustSave() {
	err := m.Save()
	if err != nil {
		panic(err)
	}
}

func (m *ChinaStockManager) LoadAllChinaStock() error {
	var (
		totalList []*ChinaStock
	)

	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	// curl "http://money.finance.sina.com.cn/q/api/jsonp_v2.php/IO.XSRV2.CallbackList^\['WH4iFBaO9ImEnmgC'^\]/US_ChinaStockService.getData?page=1&num=60&sort=&asc=0&market=&concept=0" ^
	//   -H "Connection: keep-alive" ^
	//   -H "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36" ^
	//   -H "Accept: */*" ^
	//   -H "Referer: http://finance.sina.com.cn/" ^
	//   -H "Accept-Language: zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7,ja;q=0.6" ^
	//   -H "Cookie: SINAGLOBAL=125.71.132.165_1586436411.836210; SCF=AlZv4lISPgOYj2vOcivksz9J57wav3YV6hdJIQ2Kny2bKQCDd9aMLaCLDqH1lcPdq_QxaZbH-ld7STMQpUunKv0.; sso_info=v02m6alo5qztKWRk5yljoOEpZCUiKWRk6SlkKSUpY6TpKWRk5iljpSUpY6UjKadlqWkj5OEt42DhLSOg5CzjJOQwA==; UOR=www.baidu.com,t.cj.sina.com.cn,; U_TRS1=00000066.9a137f6a.5eb16b2d.1ac880a6; SGUID=1593784977574_71636759; FINA_V_S_2=sh600900,sh600298; visited_uss=gb_ipo; UM_distinctid=177670316e8243-0025a6697a500d-31346d-384000-177670316e9bef; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WWAwC7p9sVwj.AoNlVlBJ2S5NHD95QpS0BpShnXe02XWs4DqcjsMc_LwgLo; ALF=1647529384; U_TRS2=00000017.5e0244bb.6056c047.c0a7bab3; SUB=_2A25NUrAZDeRhGedJ71MV-CrPyjiIHXVuKabRrDV_PUJbm9AfLVCskW9NVgfZhhui1YTpmqKLMtUviWqyNdrN43dK; SessionID=4g244a8293pksq92nhbldn4bs4; SINABLOGNUINFO=1741484314.67ccf11a.; ULV=1616313312569:15:1:1::1613206450825; Apache=222.209.173.142_1616313311.173365; MONEY-FINANCE-SINA-COM-CN-WEB5=; lxlrttp=1578733570" ^
	//   --compressed ^
	//   --insecure

	var (
		i = 1
	)
	for {
		if i == 6 {
			break
		}
		// TODO: This is insecure; use only in dev environments.
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}

		url := fmt.Sprintf("http://money.finance.sina.com.cn/q/api/jsonp_v2.php/IO.XSRV2.CallbackList^['WH4iFBaO9ImEnmgC'^]/US_ChinaStockService.getData?page=%d&num=60&sort=&asc=0&market=&concept=0", i)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return errGetChinaStock
		}
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36")
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Referer", "http://finance.sina.com.cn/")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7,ja;q=0.6")
		req.Header.Set("Cookie", "SINAGLOBAL=125.71.132.165_1586436411.836210; SCF=AlZv4lISPgOYj2vOcivksz9J57wav3YV6hdJIQ2Kny2bKQCDd9aMLaCLDqH1lcPdq_QxaZbH-ld7STMQpUunKv0.; sso_info=v02m6alo5qztKWRk5yljoOEpZCUiKWRk6SlkKSUpY6TpKWRk5iljpSUpY6UjKadlqWkj5OEt42DhLSOg5CzjJOQwA==; UOR=www.baidu.com,t.cj.sina.com.cn,; U_TRS1=00000066.9a137f6a.5eb16b2d.1ac880a6; SGUID=1593784977574_71636759; FINA_V_S_2=sh600900,sh600298; visited_uss=gb_ipo; UM_distinctid=177670316e8243-0025a6697a500d-31346d-384000-177670316e9bef; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WWAwC7p9sVwj.AoNlVlBJ2S5NHD95QpS0BpShnXe02XWs4DqcjsMc_LwgLo; ALF=1647529384; U_TRS2=00000017.5e0244bb.6056c047.c0a7bab3; SUB=_2A25NUrAZDeRhGedJ71MV-CrPyjiIHXVuKabRrDV_PUJbm9AfLVCskW9NVgfZhhui1YTpmqKLMtUviWqyNdrN43dK; SessionID=4g244a8293pksq92nhbldn4bs4; SINABLOGNUINFO=1741484314.67ccf11a.; ULV=1616313312569:15:1:1::1613206450825; Apache=222.209.173.142_1616313311.173365; MONEY-FINANCE-SINA-COM-CN-WEB5=; lxlrttp=1578733570")

		resp, err := client.Do(req)
		if err != nil {
			// handle err
			glog.Errorf("Get china stock failed, err: %v", err)
			return errGetChinaStock
		}
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		glog.V(4).Infof("BODY: %s", body)

		bodyStr := string(body)
		startIdx := strings.Index(bodyStr, "([{")
		endIdx := strings.Index(bodyStr, "}])")

		glog.V(4).Infof("FIX BODY: %s", bodyStr[startIdx+1:endIdx+2])

		var (
			chinaStockList []*ChinaStock
		)

		err = json.Unmarshal([]byte(bodyStr[startIdx+1:endIdx+2]), &chinaStockList)
		if err != nil {
			glog.Errorf("Unmarshal failed, err: %v", err)
			return errGetChinaStock
		}

		for _, stock := range chinaStockList {
			totalList = append(totalList, stock)
		}
		i++
	}

	glog.V(4).Infof("Total china stock: %d", len(totalList))
	for idx, stock := range totalList {
		glog.V(4).Infof("IDX: %d, STOCK: %+v", idx, stock)
		m.stocks[stock.Ticker] = stock
	}

	m.MustSave()
	return nil
}

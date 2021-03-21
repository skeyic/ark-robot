package service

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/config"
	"github.com/skeyic/ark-robot/utils"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Data folder
// - Porter
//   - ARKK
//     - 20210128ARKK.csv
//     - 20210129ARKK.csv
//   - ARKQ
//     - 20210128ARKQ.csv
//  ...

// Data folder
// - ARK
//   - 20210128
//     - 20210128ARKK.csv
//     - 20210128ARKF.csv
//     - 20210128ARKG.csv
//     - 20210128ARKQ.csv
//     - 20210128ARKW.csv
//   - 20210129
//     - 20210129ARKQ.csv
//     - 20210129ARKF.csv
//     - 20210129ARKG.csv
//     - 20210129ARKQ.csv
//     - 20210129ARKW.csv
//  ...

var (
	porterFolder = config.Config.DataFolder + "/ARK/"
)

var (
	ThePorter = NewPorter()
)

func generatePorterFilePath(fund string, date time.Time) string {
	folderPath := porterFolder + date.Format("20060102")
	utils.CheckFolder(folderPath)
	return folderPath + "/" + date.Format("20060102") + fund + ".csv"
}

func generatePorterCurrentFilePath(fund string) string {
	return porterFolder + fund + "/" + "current.csv"
}

type Porter struct {
}

func NewPorter() *Porter {
	p := &Porter{}
	return p
}

func (p *Porter) Init() error {
	utils.CheckFolder(porterFolder)
	glog.V(4).Infof("porter init completed")
	return nil
}

var (
	ARKCSVTitle = strings.Join([]string{"date", "fund", "company", "ticker", "cusip", "shares", "market value($)", "weight(%)"}, ",")
)

func ValidateARKCSV(records [][]string) bool {
	if len(records) < 2 {
		glog.Errorf("csv does not match, len of records: %d", len(records))
		return false
	}

	if strings.Join(records[0], ",") != ARKCSVTitle {
		glog.Errorf("csv does not match, title: >>%s<<, expect: >>%s<<", strings.Join(records[0], ","), ARKCSVTitle)
		return false
	}

	return true
}

// Read the CSV file, add the data in researcher, move and remove the file to correct place
func (p *Porter) Catalog(csvFileName string) (*StockHoldings, error) {
	var (
		stockHolding []*StockHolding
	)
	records, err := NewCSVReader(csvFileName).Load()
	if err != nil {
		panic(fmt.Sprintf("failed to read csv file: %s, err: %v", csvFileName, err))
	}

	if !ValidateARKCSV(records) {
		panic(fmt.Sprintf("failed to validate csv file: %s", csvFileName))
	}

	for _, record := range records[1:] {
		stockHolding = append(stockHolding, NewStockHoldingFromRecord(record))
	}

	if len(stockHolding) == 0 {
		panic("We will never reach here, but the compiler does not think so")
	}

	theDate := stockHolding[0].Date
	theFund := stockHolding[0].Fund
	newPath := generatePorterFilePath(theFund, theDate)

	// return if the file already exists
	_, err = os.Stat(newPath)
	if !os.IsNotExist(err) {
		glog.Warningf("file %s already exists, return", newPath)
		//_ = os.Remove(newPath)
	} else {
		err = os.Rename(csvFileName, newPath)
		glog.V(4).Infof("Rename %s to %s", csvFileName, newPath)
		if err != nil {
			glog.Errorf("failed to rename csv file: %s, new path: %s, err: %v", csvFileName, newPath, err)
			return nil, errRenameFile
		}
	}

	return NewStockHoldings(theDate, theFund, stockHolding), nil

	//TheLibrary.AddStockHoldings(holdings)
	//TheStockLibraryMaster.AddStockHoldings(holdings)
	//glog.V(4).Infof("Add %s at %s to library", theFund, theDate)
	//if config.Config.DebugMode {
	//	utils.SendAlertV2("Add to library", fmt.Sprintf("Add %s at %s to library", theFund, theDate))
	//}
}

func (p *Porter) ListAllDates() (files []string, err error) {
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if info.IsDir() {
			names := strings.Split(info.Name(), "/")
			name := names[len(names)-1]
			glog.V(4).Infof("NAME: %s", name)

			_, err := time.Parse("20060102", name)
			if err != nil {
				glog.Errorf("Skip incorrect file name: %s", name)
				return nil
			}
			files = append(files, path)

			return nil
		}
		return nil
	}

	err = filepath.Walk(porterFolder, walkFunc)
	if err != nil {
		return
	}

	return
}

func (p *Porter) ListAllCSVs(folderPath string) (files []string, err error) {
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(info.Name(), "csv") {
			files = append(files, path)
		}
		return nil
	}

	err = filepath.Walk(folderPath, walkFunc)
	if err != nil {
		return
	}

	return
}

func (p *Porter) ReadCSV(csvFileName string) (*StockHoldings, error) {
	var (
		stockHolding []*StockHolding
	)
	records, err := NewCSVReader(csvFileName).Load()
	if err != nil {
		glog.Errorf("failed to read csv file: %s, err: %v", csvFileName, err)
		return nil, err
	}

	if !ValidateARKCSV(records) {
		glog.Errorf("failed to validate csv file: %s", csvFileName)
		return nil, err
	}

	for _, record := range records[1:] {
		stockHolding = append(stockHolding, NewStockHoldingFromRecord(record))
	}

	if len(stockHolding) == 0 {
		panic("We will never reach here, but the compiler does not think so")
	}

	theDate := stockHolding[0].Date
	theFund := stockHolding[0].Fund

	// TODO
	return NewStockHoldings(theDate, theFund, stockHolding), nil
	//TheLibrary.AddStockHoldings(holding)
	//TheStockLibraryMaster.AddStockHoldings(holding)

	//newPath := "/Users/carrick/go/src/github.com/skeyic/ark-robot/data/ARK"
	//utils.CheckFolder("/Users/carrick/go/src/github.com/skeyic/ark-robot/data/ARK")
	//utils.CheckFolder("/Users/carrick/go/src/github.com/skeyic/ark-robot/data/ARK/" + theDate.Format("20060102"))
	//
	//names := strings.Split(csvFileName, "/")
	//err = os.Rename(csvFileName, "/Users/carrick/go/src/github.com/skeyic/ark-robot/data/ARK/"+
	//	theDate.Format("20060102")+"/"+names[len(names)-1])
	//if err != nil {
	//	panic(fmt.Sprintf("failed to rename csv file: %s, new path: %s, err: %v", csvFileName, newPath, err))
	//}

}

func (p *Porter) LoadFromDirectory() (err error) {
	// Load all holdings
	dates, err := ThePorter.ListAllDates()
	if err != nil {
		glog.Errorf("failed to list all csv files, err: %v", err)
		return
	}

	for _, dateFolder := range dates {
		glog.V(10).Infof("DATE_FOLDER: %s", dateFolder)
		arkHoldings, err := NewARKHoldingsFromDirectory(dateFolder)
		if err != nil {
			return err
		}
		TheLibrary.AddStockHoldings(arkHoldings)
		TheStockLibraryMaster.AddStockHoldings(arkHoldings)
	}

	return nil
}

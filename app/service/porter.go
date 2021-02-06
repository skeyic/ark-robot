package service

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/config"
	"github.com/skeyic/ark-robot/utils"
	"os"
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

var (
	porterFolder = config.Config.DataFolder + "/directory/"
)

var (
	ThePorter = NewPorter()
)

func generatePorterFilePath(fund string, date time.Time) string {
	return porterFolder + fund + "/" + date.Format("20060102") + fund + ".csv"
}

func generatePorterCurrentFilePath(fund string) string {
	return porterFolder + fund + "/" + "current.csv"
}

type Porter struct {
}

func NewPorter() *Porter {
	p := &Porter{}
	p.init()
	return p
}

func (p *Porter) init() {
	utils.CheckFolder(porterFolder)
	for _, arkType := range allARKTypes {
		utils.CheckFolder(porterFolder + arkType)
	}
	glog.V(4).Infof("porter init completed")
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
		glog.Errorf("csv does not match, title: %s, expect: %s", strings.Join(records[0], ","), ARKCSVTitle)
		return false
	}

	return true
}

// Read the CSV file, add the data in researcher, move and remove the file to correct place
func (p *Porter) Catalog(csvFileName string) {
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
	err = os.Rename(csvFileName, newPath)
	if err != nil {
		panic(fmt.Sprintf("failed to rename csv file: %s, new path: %s, err: %v", csvFileName, newPath, err))
	}

	TheResearcher.AddStockHoldings(NewStockHoldings(theDate, theFund, stockHolding))
}

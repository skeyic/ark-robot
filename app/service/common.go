package service

import (
	"errors"
	"github.com/skeyic/ark-robot/config"
)

const (
	TheDateFormat   = "2006-01-02"
	TheDateIDFormat = "2006_01_02"
)

var (
	errDownloadCSV        = errors.New("download csv failed")
	errFileAlreadyExist   = errors.New("file already exist")
	errRenameFile         = errors.New("rename file failed")
	errDateNotMatch       = errors.New("date not match")
	errFundNotMatch       = errors.New("fund not match")
	errValidateFail       = errors.New("validate failed")
	errNoLatestDate       = errors.New("no latest date")
	errEmptySourceToIndex = errors.New("empty to index")
	errInitReportFile     = errors.New("failed to init report file")
	errEmptyReport        = errors.New("nothing to report")
)

var (
	reportPath            = config.Config.DataFolder + "/report"
	tradingsExcelTemplate = config.Config.ResourceFolder + "/ARK.xlsx"
	top10ExcelTemplate    = config.Config.ResourceFolder + "/top_10_stocks_in_funds.xlsx"
	tradingsSheet         = "sheet"
)

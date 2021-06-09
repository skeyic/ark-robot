package service

import (
	"errors"
	"strings"
	"time"
)

const (
	TheDateFormat   = "2006-01-02"
	TheDateIDFormat = "2006_01_02"
)

var (
	allARKTypes     = []string{"ARKF", "ARKG", "ARKK", "ARKQ", "ARKW", "ARKX"}
	theFirstARKType = "ARKF"
	theLastARKType  = "ARKX"
	arkxDate, _     = time.Parse(TheDateFormat, "2021-03-30")

	AllFunds = strings.Join(allARKTypes, ",")
	TheTotal = "TOTAL"
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
	errGetChinaStock      = errors.New("get china stock failed")
	errStockNotFound      = errors.New("stock not found")
	errStockNotHold       = errors.New("stock not hold")
	errNoDataInDateRange  = errors.New("no data in date range")
)

var IgnoreTickers = []string{
	"MORGAN_STANLEY_GOVT_INSTL_8035",
	"JAPANESE_YEN",
	"HONG_KONG_DOLLAR",
	"DREYFUS_GOVT_CASH_MAN_INS",
	"CANADIAN_DOLLAR",
	"SOUTH_AFRICAN_RAND",
}

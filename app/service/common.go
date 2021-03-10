package service

import "errors"

const (
	TheDateFormat = "2006-01-02"
)

var (
	errDownloadCSV      = errors.New("download csv failed")
	errFileAlreadyExist = errors.New("file already exist")
	errRenameFile       = errors.New("rename file failed")
	errDateNotMatch     = errors.New("date not match")
	errFundNotMatch     = errors.New("fund not match")
	errValidateFail     = errors.New("validate failed")
	errNoLatestDate     = errors.New("no latest date")
)

package service

import (
	"fmt"
	"github.com/skeyic/ark-robot/utils"
)

func init() {
	utils.CheckFolder(reportPath)
	utils.CheckFile(tradingsExcelTemplate)
	utils.CheckFile(top10ExcelTemplate)
}

func toSkipTrade(direction TradeDirection) bool {
	//return false
	return direction == TradeDoNothing || direction == TradeKeep
}

func toSkipTicker(ticker string) bool {
	return ticker == "MORGAN_STANLEY_GOVT_INSTL_8035"
}

func floatToPercentString(percent float64) string {
	return fmt.Sprintf("%.3f", percent) + "%"
}

func floatToString(percent float64) string {
	return fmt.Sprintf("%.3f", percent)
}

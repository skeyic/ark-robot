package service

import (
	"fmt"
	"github.com/skeyic/ark-robot/config"
	"github.com/skeyic/ark-robot/utils"
	"math"
)

const (
	prefixTradings                           = "tradings_"
	prefixTop10Holdings                      = "top_10_stocks_"
	prefixSpecialTradings                    = "special_tradings_"
	prefixSpecialTradingsHigherThan10        = "special_tradings_higher_than_10_"
	prefixSpecialTradingsContinuousDirection = "special_tradings_continues_direction_"
	prefixChinaStockTradings                 = "china_stock_tradings_"
	prefixChinaStockNotKeepTradings          = "china_stock_not_keep_tradings_"
	prefixCurrentStockReport                 = "current_stock_report_"
	prefixDateRangeStockReport               = "date_range_stock_report_"
	defaultSheet                             = "sheet"
)

var (
	reportPath                      = config.Config.DataFolder + "/report"
	stockReportPath                 = config.Config.DataFolder + "/stock_report"
	tradingsExcelTemplate           = config.Config.ResourceFolder + "/TEMPLATE_ARK.xlsx"
	top10ExcelTemplate              = config.Config.ResourceFolder + "/TEMPLATE_top_10_stocks.xlsx"
	specialTradingsExcelTemplate    = config.Config.ResourceFolder + "/TEMPLATE_special_tradings.xlsx"
	chinaStockExcelTradingsTemplate = config.Config.ResourceFolder + "/TEMPLATE_china_stock_tradings.xlsx"
	stockReportExcelTemplate        = config.Config.ResourceFolder + "/TEMPLATE_stock_report.xlsx"
)

func init() {
	utils.CheckFolder(reportPath)
	utils.CheckFolder(stockReportPath)

	utils.CheckFile(tradingsExcelTemplate)
	utils.CheckFile(top10ExcelTemplate)
	utils.CheckFile(stockReportExcelTemplate)
}

func toSkipTrade(direction TradeDirection) bool {
	//return false
	return direction == TradeDoNothing || direction == TradeKeep
}

func toSkipTicker(ticker string) bool {
	return ticker == "MORGAN_STANLEY_GOVT_INSTL_8035" ||
		ticker == "JAPANESE_YEN" ||
		ticker == "HONG_KONG_DOLLAR" ||
		ticker == "DREYFUS_GOVT_CASH_MAN_INS" ||
		ticker == "CANADIAN_DOLLAR"
}

func floatToPercentString(percent float64) string {
	return floatToString(percent) + "%"
}

func floatToString(percent float64) (result string) {
	result += fmt.Sprintf("%.2f", percent)
	return
}

func floatToPercentStringWithSign(percent float64) string {
	return floatToStringWithSign(percent) + "%"
}

func floatToStringWithSign(percent float64) (result string) {
	if percent > 0 {
		result += "+"
	}
	result += fmt.Sprintf("%.2f", percent)
	return
}

func floatToStringIntOnly(data float64) (result string) {
	result += fmt.Sprintf("%.0f", math.Ceil(data))
	return
}

func floatToStringIntOnlyWithSign(data float64) (result string) {
	if data > 0 {
		result += "+"
	}
	result += fmt.Sprintf("%.0f", math.Ceil(data))
	return
}

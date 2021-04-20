package control

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/app/service"
	"github.com/skeyic/ark-robot/config"
	"github.com/skeyic/ark-robot/utils"
	"strconv"
	"strings"
	"time"
)

// @Summary Download
// @Tags Admin
// @Description let the master download latest data
// @Accept json
// @Produce json
// @Success 200 {object} utils.WebResponse "Ok"
// @Failure 400 {object} utils.WebResponse "Bad request"
// @Failure 500 {object} utils.WebResponse "Internal error"
// @Router /admin/download [post]
func DoDownload(c *gin.Context) {
	var (
		err error
	)

	err = service.TheDownloader.DownloadAllARKCSVs()
	if err != nil {
		msg := fmt.Sprintf("failed to download, err: %v", err)
		glog.Error(msg)
		utils.NewBadRequestError(c, msg)
		return
	}

	utils.NewOkResponse(c, "downloaded")
}

// @Summary Report
// @Tags Report
// @Description let the master report special date
// @Accept json
// @Produce json
// @Param date query string true "The report date"
// @Param full query bool false "Full report or not"
// @Param percent query float64 false "Special trading percent"
// @Success 200 {object} utils.WebResponse "Ok"
// @Failure 400 {object} utils.WebResponse "Bad request"
// @Failure 500 {object} utils.WebResponse "Internal error"
// @Router /report/report [post]
func DoReport(c *gin.Context) {
	var (
		err                   error
		date                  time.Time
		full                  = true
		specialTradingPercent = config.Config.Report.SpecialTradingPercent
	)

	dateInput := c.Query("date")
	date, err = time.Parse(service.TheDateFormat, dateInput)
	if err != nil {
		msg := fmt.Sprintf("Incorrect date: %s, should be like %s, err: %v", dateInput, service.TheDateFormat, err)
		glog.Error(msg)
		utils.NewBadRequestError(c, msg)
		return
	}

	fullInput := c.Query("full")
	if fullInput != "" {
		full, err = strconv.ParseBool(fullInput)
		if err != nil {
			msg := fmt.Sprintf("Incorrect full: %s, should be a bool, err: %v", fullInput, err)
			glog.Error(msg)
			utils.NewBadRequestError(c, msg)
			return
		}
	}

	specialTradingPercentInput := c.Query("percent")
	if specialTradingPercentInput != "" {
		specialTradingPercent, err = strconv.ParseFloat(specialTradingPercentInput, 64)
		if err != nil {
			msg := fmt.Sprintf("Incorrect specialTradingPercent: %s, should be a float64, err: %v", specialTradingPercentInput, err)
			glog.Error(msg)
			utils.NewBadRequestError(c, msg)
			return
		}
	}

	err = service.TheMaster.Report(date, full, specialTradingPercent)
	if err != nil {
		msg := fmt.Sprintf("failed to report date: %s, full: %v, err: %v", dateInput, fullInput, err)
		glog.Error(msg)
		utils.NewBadRequestError(c, msg)
		return
	}

	utils.NewOkResponse(c, fmt.Sprintf("report finished, date: %s", dateInput))
}

// @Summary ReportStock
// @Tags Report
// @Description let the master report special stock in a date range
// @Accept json
// @Produce json
// @Param stock query string true "The stock ticker"
// @Param from_date query string true "The report from date"
// @Param end_date query string true "The report end date"
// @Success 200 {object} utils.WebResponse "Ok"
// @Failure 400 {object} utils.WebResponse "Bad request"
// @Failure 500 {object} utils.WebResponse "Internal error"
// @Router /report/report_stock [post]
func DoReportStock(c *gin.Context) {
	var (
		err               error
		fromDate, endDate time.Time
		ticker            string
	)

	fromDateInput := c.Query("from_date")
	fromDate, err = time.Parse(service.TheDateFormat, fromDateInput)
	if err != nil {
		msg := fmt.Sprintf("Incorrect fromDate: %s, should be like %s, err: %v", fromDateInput, service.TheDateFormat, err)
		glog.Error(msg)
		utils.NewBadRequestError(c, msg)
		return
	}

	endDateInput := c.Query("end_date")
	endDate, err = time.Parse(service.TheDateFormat, endDateInput)
	if err != nil {
		msg := fmt.Sprintf("Incorrect endDate: %s, should be like %s, err: %v", endDateInput, service.TheDateFormat, err)
		glog.Error(msg)
		utils.NewBadRequestError(c, msg)
		return
	}

	ticker = c.Query("stock")
	if ticker == "" {
		msg := fmt.Sprintf("Empty stock")
		glog.Error(msg)
		utils.NewBadRequestError(c, msg)
		return
	}

	tickers := strings.Split(ticker, ",")
	for _, theTicker := range tickers {
		err = service.TheMaster.ReportStock(theTicker, fromDate, endDate)
		if err != nil {
			msg := fmt.Sprintf("failed to report stock %s, fromDate: %s, endDate: %s, err: %v", ticker, fromDateInput, endDateInput, err)
			glog.Error(msg)
			utils.NewBadRequestError(c, msg)
			return
		}
	}

	utils.NewOkResponse(c, fmt.Sprintf("report finished, stock %s, fromDate: %s, endDate: %s", ticker, fromDateInput, endDateInput))
}

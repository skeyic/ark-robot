package control

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/app/service"
	"github.com/skeyic/ark-robot/utils"
	"strings"
)

// ReportStockCurrent
// @Summary ReportStockCurrent
// @Tags Data
// @Description let the master report special stock current status
// @Accept json
// @Produce json
// @Param ticker path string true "The stock ticker"
// @Success 200 {object} utils.WebResponse "Ok"
// @Failure 400 {object} utils.WebResponse "Bad request"
// @Failure 500 {object} utils.WebResponse "Internal error"
// @Router /data/reports/stocks/{ticker}/current [post]
func ReportStockCurrent(c *gin.Context) {
	var (
		err    error
		report string
		ticker string
	)

	ticker = c.Param("ticker")
	if ticker == "" {
		msg := fmt.Sprintf("Empty stock")
		glog.Error(msg)
		utils.NewBadRequestError(c, msg)
		return
	}

	report, err = service.TheMaster.ReportStockCurrent(ticker)
	if err != nil {
		msg := fmt.Sprintf("failed to report stock %s, err: %v", ticker, err)
		glog.Error(msg)
		utils.NewBadRequestError(c, msg)
		return
	}

	utils.NewOkResponse(c, report)
}

// ReportFundTop10
// @Summary ReportFundTop10
// @Tags Data
// @Description let the master report special fund top 10
// @Accept json
// @Produce json
// @Param fund path string true "The fund name"
// @Success 200 {object} utils.WebResponse "Ok"
// @Failure 400 {object} utils.WebResponse "Bad request"
// @Failure 500 {object} utils.WebResponse "Internal error"
// @Router /data/reports/funds/{fund}/top10 [post]
func ReportFundTop10(c *gin.Context) {
	var (
		err    error
		report string
		fund   string
	)

	fund = c.Param("fund")
	if fund == "" {
		msg := fmt.Sprintf("Empty fund")
		glog.Error(msg)
		utils.NewBadRequestError(c, msg)
		return
	}

	isFund := service.IsFund(fund)
	if !isFund {
		msg := fmt.Sprintf("No such fund: %s, we only support %s", fund, service.AllFunds)
		glog.Error(msg)
		utils.NewNotFoundError(c, msg)
		return
	}

	report, err = service.TheMaster.ReportFundTop10(fund)
	if err != nil {
		msg := fmt.Sprintf("failed to report fund %s top 10, err: %v", fund, err)
		glog.Error(msg)
		utils.NewBadRequestError(c, msg)
		return
	}

	utils.NewOkResponse(c, report)
}

// GetAllTickers
// @Summary GetAllTickers
// @Tags Data
// @Description let the master report special stock current status
// @Accept json
// @Produce json
// @Success 200 {object} utils.WebResponse "All tickers"
// @Failure 400 {object} utils.WebResponse "Bad request"
// @Failure 500 {object} utils.WebResponse "Internal error"
// @Router /data/tickers [get]
func GetAllTickers(c *gin.Context) {
	tickers := service.TheMaster.GetAllTickers()

	utils.NewOkResponse(c, strings.Join(tickers, ","))
}

// IsTicker
// @Summary IsTicker
// @Tags Data
// @Description let the master report special stock current status
// @Accept json
// @Produce json
// @Param ticker path string true "The stock ticker"
// @Success 200 {object} utils.WebResponse "OK"
// @Failure 404 {object} utils.WebResponse "Not found"
// @Failure 400 {object} utils.WebResponse "Bad request"
// @Failure 500 {object} utils.WebResponse "Internal error"
// @Router /data/tickers/{ticker} [get]
func IsTicker(c *gin.Context) {
	var (
		ticker string
	)

	ticker = c.Param("ticker")
	if ticker == "" {
		msg := fmt.Sprintf("Empty stock")
		glog.Error(msg)
		utils.NewBadRequestError(c, msg)
		return
	}

	hit := service.TheMaster.IsTicker(ticker)
	if hit {
		utils.NewOkResponse(c, "ticker "+ticker+" found")
	} else {
		utils.NewNotFoundError(c, "ticker "+ticker+" not found")
	}
}

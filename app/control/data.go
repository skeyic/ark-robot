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
// @Router /data/reports/{ticker}/current [post]
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

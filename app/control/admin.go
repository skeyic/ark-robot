package control

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/app/service"
	"github.com/skeyic/ark-robot/utils"
	"strconv"
	"time"
)

// @Summary Download
// @Tags User
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
// @Tags User
// @Description let the master report special date
// @Accept json
// @Produce json
// @Param date query string true "The report date"
// @Param full query bool false "Full report or not"
// @Success 200 {object} utils.WebResponse "Ok"
// @Failure 400 {object} utils.WebResponse "Bad request"
// @Failure 500 {object} utils.WebResponse "Internal error"
// @Router /admin/report [post]
func DoReport(c *gin.Context) {
	var (
		err  error
		date time.Time
		full bool
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
	full, err = strconv.ParseBool(fullInput)
	if err != nil {
		msg := fmt.Sprintf("Incorrect full: %s, should be a bool, err: %v", fullInput, err)
		glog.Error(msg)
		utils.NewBadRequestError(c, msg)
		return
	}

	err = service.TheMaster.Report(date, full)
	if err != nil {
		msg := fmt.Sprintf("failed to report date: %s, full: %v, err: %v", dateInput, fullInput, err)
		glog.Error(msg)
		utils.NewBadRequestError(c, msg)
		return
	}

	utils.NewOkResponse(c, fmt.Sprintf("report finished, date: %s", dateInput))
}

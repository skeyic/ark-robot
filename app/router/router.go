package router

import (
	"github.com/gin-gonic/gin"
	"github.com/skeyic/ark-robot/app/control"
	"github.com/skeyic/ark-robot/config"
	_ "github.com/skeyic/ark-robot/docs"
)

// InitRouter ...
func InitRouter() *gin.Engine {
	if config.Config.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Basic
	r.GET("/", control.Index)

	// Admin
	admin := r.Group("/admin")
	{
		admin.POST("/download", control.DoDownload)
	}

	// Report
	report := r.Group("/report")
	{
		report.POST("/report", control.DoReport)
		report.POST("/report_stock", control.DoReportStock)
		report.POST("/report_stock_by_days", control.DoReportStockByDays)
		report.POST("/report_stock_current", control.DoReportStockCurrent)
	}

	return r
}

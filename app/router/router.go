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

	// Action
	actions := r.Group("/actions")
	{
		{
			actions.POST("/download", control.Download)
		}

		reports := actions.Group("/reports")
		{
			reports.POST("", control.Report)
			reports.POST("/stock", control.ReportStock)
			reports.POST("/stock_by_days", control.ReportStockByDays)
		}
	}

	// Data
	data := r.Group("/data")
	{
		reports := data.Group("/reports")
		{
			reports.POST("/stocks/:ticker/current", control.ReportStockCurrent)
			reports.POST("/funds/:fund/top10", control.ReportFundTop10)
			reports.GET("/funds/all/continue3days", control.ReportContinue3Days)
			reports.GET("/funds/all/bigswings", control.ReportBigSwings)
			reports.GET("/funds/all/chinastock", control.ReportChinaStock)
		}

		tickers := data.Group("/tickers")
		{
			tickers.GET("", control.GetAllTickers)
			tickers.GET(":ticker", control.IsTicker)
		}
	}

	//Debug
	//debug := r.Group("/debug", control.Debug)

	return r
}

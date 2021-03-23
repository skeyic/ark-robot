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

	// User
	admin := r.Group("/admin")
	{
		admin.POST("/download", control.DoDownload)
		admin.POST("/report", control.DoReport)
	}

	return r
}

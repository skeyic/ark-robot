package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/app/router"
	"github.com/skeyic/ark-robot/app/rpc"
	"github.com/skeyic/ark-robot/app/service"
	"github.com/skeyic/ark-robot/config"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title The ARK ROBOT Server
// @version 1.0
// @description The ARK ROBOT REST API

// @host
// @BasePath /

func main() {
	var (
		err error
	)

	flag.Parse()
	err = service.TheMaster.FreshInit()
	if err != nil {
		panic(fmt.Sprintf("master failed to fresh init, err: %v", err))
	}

	go rpc.TheServer.Start()

	go service.TheMaster.StartDownload()
	go service.TheMaster.ReportLatestTrading(true)

	r := router.InitRouter()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	glog.V(4).Info("ARK ROBOT Server starts...")
	glog.Fatal(r.Run(fmt.Sprintf(":%d", config.Config.Port)))

	glog.V(4).Info("ARK ROBOT starts...")
}

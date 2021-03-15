package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/app/service"
)

func main() {
	//var (
	//	err error
	//)

	flag.Parse()
	err := service.TheMaster.FreshInit()
	if err != nil {
		panic(fmt.Sprintf("master failed to fresh init, err: %v", err))
	}
	service.TheMaster.StartDownload()

	glog.V(4).Info("Ark robot starts...")
}

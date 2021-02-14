package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/app/service"
)

func main() {
	//var (
	//	err error
	//)

	flag.Parse()
	service.TheMaster.Start()

	glog.V(4).Info("Ark robot starts...")
}

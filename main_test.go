package main

import (
	"fmt"
	"github.com/skeyic/ark-robot/app/rpc"
	"github.com/skeyic/ark-robot/app/service"
	"github.com/skeyic/ark-robot/config"
	"github.com/skeyic/ark-robot/utils"
	"testing"
	"time"
)

func TestGRPCServer(t *testing.T) {
	var (
		err error
	)
	utils.EnableGlogForTesting()

	err = service.TheMaster.StaleInit()
	if err != nil {
		panic(fmt.Sprintf("master failed to fresh init, err: %v", err))
	}

	go rpc.TheServer.Start()
	time.Sleep(5 * time.Second)

	TheClient := rpc.NewClient(fmt.Sprintf("localhost:%d", config.Config.RpcPort))
	TheClient.GetCurrentStockReport("JD")
}

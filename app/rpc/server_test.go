package rpc

import (
	"fmt"
	"github.com/skeyic/ark-robot/config"
	"github.com/skeyic/ark-robot/utils"
	"testing"
	"time"
)

func TestGRPC(t *testing.T) {
	utils.EnableGlogForTesting()

	go TheServer.Start()

	time.Sleep(3 * time.Second)
	TheClient := &Client{server: fmt.Sprintf("localhost:%d", config.Config.RpcPort)}
	TheClient.GetCurrentStockReport()
}

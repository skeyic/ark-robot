package rpc

import (
	"fmt"
	"github.com/golang/glog"
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
	report, err := TheClient.GetCurrentStockReport("JD")
	glog.V(4).Infof("REPORT: %s", report)
	glog.V(4).Infof("ERROR: %v", err)
	err = TheClient.Hello()
	glog.V(4).Infof("HELLO: %v", err)
}

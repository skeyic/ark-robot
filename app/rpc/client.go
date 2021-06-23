package rpc

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/app/rpc/cocoa"
	"google.golang.org/grpc"
)

type Client struct {
	server string
}

func (c *Client) GetCurrentStockReport() {
	conn, err := grpc.Dial(c.server, grpc.WithInsecure())
	if err != nil {
		panic(fmt.Sprintf("failed to connect %s", c.server))
	}

	defer conn.Close()

	t := cocoa.NewWaiterClient(conn)

	tr, err := t.GetCurrentStockReport(context.Background(), &cocoa.Req{JsonStr: "TEST COCCA"})
	if err != nil {
		panic(fmt.Sprintf("failed to get current stock report: %v", err))
	}

	glog.V(4).Infof("RESPONSE: %s", tr.BackJson)
}

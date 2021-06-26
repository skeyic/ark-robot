package rpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/app/rpc/cocoa"
	"google.golang.org/grpc"
	"time"
)

var (
	ErrConnectServer = errors.New("failed to connect rpc server")
	ErrHelloMismatch = errors.New("the message in hello mismatch")
)

type Client struct {
	server string
}

func NewClient(server string) *Client {
	return &Client{server: server}
}

func (c *Client) GetCurrentStockReport(ticker string) (string, error) {
	conn, err := grpc.Dial(c.server, grpc.WithInsecure())
	if err != nil {
		panic(fmt.Sprintf("GRPC: failed to connect %s", c.server))
	}

	defer conn.Close()

	t := cocoa.NewWaiterClient(conn)

	tr, err := t.GetCurrentStockReport(context.Background(), &cocoa.Req{JsonStr: ticker})
	if err != nil {
		glog.Errorf("GRPC: failed to get current stock report: %v", err)
		return "", err
	}

	glog.V(4).Infof("GRPC: RESPONSE: %s", tr.BackJson)
	return tr.BackJson, err
}

func (c *Client) Hello() error {
	conn, err := grpc.Dial(c.server, grpc.WithInsecure())
	if err != nil {
		panic(fmt.Sprintf("GRPC: failed to connect %s", c.server))
	}

	defer conn.Close()

	var (
		message = time.Now().String()
	)

	t := cocoa.NewWaiterClient(conn)

	tr, err := t.Hello(context.Background(), &cocoa.Req{JsonStr: message})
	if err != nil {
		glog.Errorf("GRPC failed to hello, err: %v", err)
		return err
	}

	if tr.BackJson != message {
		glog.Errorf("GRPC: message mismatch, expect: %s, actual: %s", message, tr.BackJson)
		return ErrHelloMismatch
	}
	return nil
}

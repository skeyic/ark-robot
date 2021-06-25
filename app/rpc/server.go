package rpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/app/rpc/cocoa"
	"github.com/skeyic/ark-robot/app/service"
	"github.com/skeyic/ark-robot/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

var (
	TheServer = &Server{port: config.Config.RpcPort}
)

var (
	ErrStartGRPCServerFailed = errors.New("failed to start GRPC server")
)

type Server struct {
	port int
}

func (s *Server) Start() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		glog.Fatalf("grpc failed to listen %d", s.port)
	}

	ts := grpc.NewServer()

	cocoa.RegisterWaiterServer(ts, &Server{})

	reflection.Register(ts)

	glog.V(4).Infof("START SERVER")
	err = ts.Serve(lis)
	if err != nil {
		glog.Fatalf("grpc failed to serve: %v", err)
	}
}

func (s *Server) GetCurrentStockReport(ctx context.Context, in *cocoa.Req) (*cocoa.Res, error) {
	var (
		err    error
		report string
	)

	report, err = service.TheMaster.ReportStockCurrent(in.JsonStr)
	if err != nil {
		glog.Errorf("GRPC SERVER: failed to report stock current, err: %v", err)
	}
	return &cocoa.Res{
		BackJson: report,
	}, err
}

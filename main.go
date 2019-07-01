package main

import (
	"flag"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	configFlag = flag.String("config", "", "path to config file")
)

func main() {
	flag.Parse()
	logger, _ := zap.NewDevelopment()

	service, err := NewRegBox(*configFlag)
	if err != nil {
		logger.Fatal("Can not load config", zap.String("path", *configFlag), zap.String("error", err.Error()))
	}

	listener, err := net.Listen("tcp", service.Address)
	if err != nil {
		logger.Fatal("Can not listen", zap.String("address", service.Address), zap.String("error", err.Error()))
	}

	var server = grpc.NewServer()
	RegisterRegBoxServer(server, service)
	err = server.Serve(listener)
	if err != nil {
		logger.Fatal("Can not serve", zap.String("error", err.Error()))
	}
}

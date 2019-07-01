package main

import (
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	address = "0.0.0.0:23400"
)

func main() {
	logger, _ := zap.NewDevelopment()

	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatal("Can not listen", zap.String("address", address), zap.String("error", err.Error()))
	}
	var s = grpc.NewServer()
	RegisterRegBoxServer(s, &RegBox{})
	err = s.Serve(listener)
	if err != nil {
		logger.Fatal("Can not serve", zap.String("error", err.Error()))
	}
}

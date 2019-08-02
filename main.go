package main

import (
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	logger, _ := zap.NewDevelopment()

	service, err := NewRegBox()
	if err != nil {
		logger.Fatal("Can not create service", zap.String("error", err.Error()))
	}

	listener, err := net.Listen("tcp", service.Address)
	if err != nil {
		logger.Fatal("Can not listen", zap.String("address", service.Address), zap.String("error", err.Error()))
	}

	creds, err := credentials.NewServerTLSFromFile("assets/regbox.crt", "assets/regbox.key")
	if err != nil {
		logger.Fatal("Can not create credentials", zap.String("cert path", "assets/regbox.crt"), zap.String("key path", "assets/regbox.key"), zap.String("error", err.Error()))
	}

	var server = grpc.NewServer(grpc.Creds(creds))
	RegisterRegBoxServer(server, service)

	if err = server.Serve(listener); err != nil {
		logger.Fatal("Can not serve", zap.String("error", err.Error()))
	}
}

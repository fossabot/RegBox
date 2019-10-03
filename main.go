package main

import (
	"net"

	"github.com/Aded175/RegBox/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// logger, _ := zap.NewDevelopment()

	var svc RegBoxService
	svc, _ = NewRegBoxService()

	var grpcServer = NewGRPCServer(Endpoints{
		Register:     makeRegisterEndpoint(svc),
		Authenticate: makeAuthenticateEndpoint(svc),
	})

	var baseServer = grpc.NewServer()
	reflection.Register(baseServer)
	pb.RegisterRegBoxServer(baseServer, grpcServer)

	grpcListener, _ := net.Listen("tcp", "0.0.0.0:23400")

	_ = baseServer.Serve(grpcListener)
}

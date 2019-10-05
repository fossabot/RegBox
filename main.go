package main

import (
	"context"
	"net"

	"github.com/Aded175/RegBox/pb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	conn, _ := mongo.Connect(context.Background(),
		options.Client().SetHosts([]string{"127.0.0.1:27017"}).
			SetAuth(options.Credential{
				AuthSource: "regbox",
				Username:   "regbox",
				Password:   "P@ssw0rd",
			}),
	)

	var svc RegBoxService
	svc, _ = NewRegBoxService(conn)

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

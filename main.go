package main

import (
	"context"
	"net"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	logger         *zap.Logger
	collection     *mongo.Collection
	initCollection = make(chan bool)
	err            error
)

func main() {
	logger, _ = zap.NewDevelopment()

	go func() {
		collection, err = GetCollection()
		if err != nil {
			logger.Fatal("Can not create connection to mongodb", zap.String("err", err.Error()))
		}
		initCollection <- true
		logger.Info("Connected to mongodb")
	}()

	var service = NewRegBox()

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
	<-initCollection
	err = server.Serve(listener)
	if err != nil {
		logger.Fatal("Can not serve", zap.String("error", err.Error()))
	}
}

func GetCollection() (*mongo.Collection, error) {
	client, err := mongo.Connect(
		context.Background(),
		options.Client().SetHosts([]string{"127.0.0.1:27017"}).
			SetAuth(options.Credential{
				AuthSource: "regbox",
				Username:   "regbox",
				Password:   "P@ssw0rd",
			}))
	if err != nil {
		return nil, err
	}
	return client.Database("regbox").Collection("creds"), nil
}

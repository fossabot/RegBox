package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type accountRepositoryService struct {
	conn *mongo.Client
	col  *mongo.Collection
}

func NewAccountRepositoryService() (*accountRepositoryService, error) {
	conn, err := mongo.Connect(context.Background(),
		options.Client().
			SetHosts([]string{"127.0.0.1:27017"}).
			SetAuth(options.Credential{
				AuthSource: "regbox",
				Username:   "regbox",
				Password:   "P@ssw0rd",
			}))
	if err != nil {
		return nil, err
	}

	var col = conn.Database("regbox").Collection("creds")

	return &accountRepositoryService{
		conn: conn,
		col:  col,
	}, nil
}

func (s accountRepositoryService) AddAccount(login, hash, salt []byte) (err error) {
	_, err = s.col.InsertOne(
		context.Background(),
		bson.M{
			"login":  login,
			"passwd": hash,
			"salt":   salt,
		})
	return
}

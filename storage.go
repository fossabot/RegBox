package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type storageService struct {
	account *mongo.Collection
}

func NewStorageService(conn *mongo.Client) *storageService {
	return &storageService{
		account: conn.Database("regbox").Collection("creds"),
	}
}

func (s storageService) AddAccount(login, hash, salt []byte) (err error) {
	_, err = s.account.InsertOne(
		context.Background(),
		bson.M{
			"login":  login,
			"passwd": hash,
			"salt":   salt,
		},
	)
	return
}

func (s storageService) GetHashSalt(login []byte) (hash, salt []byte, err error) {
	cursor, err := s.account.Find(
		context.Background(),
		bson.M{
			"login": login,
		},
		options.Find().SetProjection(
			bson.M{
				"_id":    0,
				"passwd": 1,
				"salt":   1,
			},
		),
	)
	if err != nil {
		return
	}
	var response map[string][]byte
	err = cursor.Decode(&response)
	if err != nil {
		return
	}
	return response["passwd"], response["salt"], nil
}

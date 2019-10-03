package main

import (
	"context"
	"io"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type repositoryService struct {
	conn *mongo.Client
	col  *mongo.Collection
}

func NewRepositoryService() (*repositoryService, error) {
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

	return &repositoryService{
		conn: conn,
		col:  col,
	}, nil
}

func (s repositoryService) CountLogins(login []byte) (int, error) {
	cursor, err := s.col.Aggregate(context.Background(),
		bson.A{
			bson.M{
				"$match": bson.M{
					"login": bson.M{
						"$eq": login,
					},
				},
			},
			bson.M{
				"$count": "logins",
			},
		},
	)
	if err != nil {
		return 0, err
	}

	_ = cursor.Next(context.Background())

	var count map[string]int
	err = cursor.Decode(&count)
	if err == io.EOF {
		// Cursor empty == login free
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return count["logins"], nil
}

func (s repositoryService) AddAccount(login, hash, salt []byte) (err error) {
	_, err = s.col.InsertOne(
		context.Background(),
		bson.M{
			"login":  login,
			"passwd": hash,
			"salt":   salt,
		})
	return
}

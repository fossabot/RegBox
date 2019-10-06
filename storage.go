package main

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type storageService struct {
	account       *mongo.Collection
	accessTokens  *mongo.Collection
	refreshTokens *mongo.Collection
}

func NewStorageService(client *mongo.Client) *storageService {
	return &storageService{
		account:       client.Database("regbox").Collection("creds"),
		accessTokens:  client.Database("regbox").Collection("accessTokens"),
		refreshTokens: client.Database("regbox").Collection("refreshTokens"),
	}
}

func (s storageService) AddAccount(ctx context.Context, login string, hash, salt []byte) (err error) {
	_, err = s.account.InsertOne(
		ctx,
		bson.M{
			"login":  login,
			"passwd": hash,
			"salt":   salt,
		},
	)
	return
}

func (s storageService) GetHashSalt(ctx context.Context, login string) (hash, salt []byte, err error) {
	cursor, err := s.account.Find(
		ctx,
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

	_ = cursor.Next(ctx)

	var response map[string][]byte
	err = cursor.Decode(&response)
	if err != nil {
		return
	}
	return response["passwd"], response["salt"], nil
}

type TokenType uint8

const (
	TOKEN_ACCESS TokenType = iota
	TOKEN_REFRESH
)

var ErrTokenType = errors.New("Invalid token type")

func (s storageService) AddToken(ctx context.Context, tt TokenType, t *Token) (err error) {
	var col *mongo.Collection
	switch tt {
	case TOKEN_ACCESS:
		col = s.accessTokens
	case TOKEN_REFRESH:
		col = s.refreshTokens
	default:
		return ErrTokenType
	}
	tmpID, err := t.id.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = col.InsertOne(
		ctx,
		bson.M{
			"expireAt": t.expire,
			"id":       tmpID,
			"login":    t.login,
			"token":    t.ss,
		},
	)
	return
}

func (s storageService) AddTokenPair(ctx context.Context, at *Token, rt *Token) (err error) {
	err = s.AddToken(ctx, TOKEN_ACCESS, at)
	if err != nil {
		return err
	}
	err = s.AddToken(ctx, TOKEN_REFRESH, rt)
	if err != nil {
		return err
	}
	return nil
}

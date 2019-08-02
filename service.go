package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"io"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/argon2"
)

type RegBox struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
	Address     string
	Collection  *mongo.Collection
}

func NewRegBox() (*RegBox, error) {
	conn, err := mongo.Connect(
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
	return &RegBox{
		Memory:      64 * 1024, // 64 Mib
		Iterations:  10,
		Parallelism: 4,
		SaltLength:  16,
		KeyLength:   32,
		Address:     "0.0.0.0:23400",
		Collection:  conn.Database("regbox").Collection("creds"),
	}, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (s *RegBox) GenerateHash(password []byte, salt []byte) []byte {
	return argon2.IDKey(password, salt, s.Iterations, s.Memory, s.Parallelism, s.KeyLength)
}

func (s *RegBox) Register(ctx context.Context, in *AcccountRequest) (*Response, error) {
	var login = []byte(in.GetLogin())
	if err := s.Check(login); err != nil {
		return &Response{Error: err.Error()}, nil
	}

	salt, err := generateRandomBytes(s.SaltLength)
	if err != nil {
		return &Response{Error: err.Error()}, nil
	}
	var hash = s.GenerateHash([]byte(in.GetPassword()), salt)

	_, err = s.Collection.InsertOne(
		context.Background(),
		bson.M{"login": login,
			"passwd": hash,
			"salt":   salt,
		})
	if err != nil {
		return &Response{Error: err.Error()}, nil
	}
	return &Response{}, nil
}

var (
	ErrLoginUsed = errors.New("Login is used")
	ErrDupLogins = errors.New("Duplicated logins")
)

func (s *RegBox) Check(login []byte) error {
	cursor, err := s.Collection.Aggregate(context.Background(),
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
		return err
	}

	_ = cursor.Next(context.Background())

	var count map[string]int
	err = cursor.Decode(&count)
	if err == io.EOF {
		// Cursor empty == login free
		return nil
	}
	if err != nil {
		return err
	}

	switch count["logins"] {
	case 1:
		return ErrLoginUsed
	default:
		return ErrDupLogins
	}
}

func (s *RegBox) Authenticate(ctx context.Context, in *AcccountRequest) (*Response, error) {
	var login = []byte(in.GetLogin())
	if s.Check(login) != ErrLoginUsed {
		return &Response{Error: "Login not found"}, nil
	}

	cursor, err := s.Collection.Find(context.Background(),
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
		return &Response{Error: err.Error()}, nil
	}

	_ = cursor.Next(context.Background())

	var response map[string][]byte
	err = cursor.Decode(&response)
	if err != nil {
		return &Response{Error: err.Error()}, nil
	}

	var hash = s.GenerateHash([]byte(in.GetPassword()), response["salt"])
	if !bytes.Equal(hash, response["passwd"]) {
		return &Response{Error: "Passwords do not match"}, nil
	}

	return &Response{}, nil
}

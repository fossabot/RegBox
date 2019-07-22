package main

import (
	"context"
	"crypto/rand"
	"runtime"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"
)

type RegBox struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
	Address     string
}

func NewRegBox() *RegBox {
	return &RegBox{
		Memory:      64 * 1024, // 64 Mib
		Iterations:  10,
		Parallelism: uint8(runtime.NumCPU()),
		SaltLength:  16,
		KeyLength:   32,
		Address:     "0.0.0.0:23400",
	}
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (s *RegBox) GenerateHash(password string) (hash []byte, salt []byte, err error) {
	salt, err = generateRandomBytes(s.SaltLength)
	if err != nil {
		return nil, nil, err
	}
	hash = argon2.IDKey([]byte(password), salt, s.Iterations, s.Memory, s.Parallelism, s.KeyLength)
	return
}

func (s *RegBox) Register(ctx context.Context, in *RegisterRequest) (*RegisterResponse, error) {
	hash, salt, err := s.GenerateHash(in.GetPassword())
	if err != nil {
		logger.Warn("Can generate hash", zap.String("err", err.Error()))
		return &RegisterResponse{Successful: false}, nil
	}
	_, err = collection.InsertOne(
		context.Background(),
		bson.M{"login": []byte(in.GetLogin()),
			"passwd": hash,
			"salt":   salt,
		})
	if err != nil {
		logger.Warn("Can not insert one", zap.String("err", err.Error()))
		return &RegisterResponse{Successful: false}, nil
	}
	return &RegisterResponse{Successful: true}, nil
}

package main

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type RegBoxService interface {
	Register(string, string) (string, error)
	Authenticate(string, string) (string, string, error)
}

type regBoxService struct {
	stor   *storageService
	crypto *cryptoService
}

func NewRegBoxService(client *mongo.Client) (*regBoxService, error) {
	var s = NewStorageService(client)
	c, err := NewCryptoService()
	if err != nil {
		return nil, err
	}
	return &regBoxService{
		stor:   s,
		crypto: c,
	}, nil
}

func (s regBoxService) Register(login string, password string) (string, error) {
	salt, err := s.crypto.GenerateSalt()
	if err != nil {
		return "", err
	}
	var hash = s.crypto.GenerateHash([]byte(password), salt)

	err = s.stor.AddAccount(context.Background(), login, hash, salt)
	if err != nil {
		return "", err
	}
	return login, nil
}

func (s regBoxService) Authenticate(login string, password string) (string, string, error) {
	hash, salt, err := s.stor.GetHashSalt(context.Background(), login)
	if err != nil {
		return "", "", err
	}
	var actualHash = s.crypto.GenerateHash([]byte(password), salt)

	if !bytes.Equal(hash, actualHash) {
		return "", "", ErrPasswdIncorrect
	}

	var (
		now           = time.Now()
		expireAccess  = now.Add(3 * time.Minute)
		expireRefresh = now.Add(6 * time.Minute)
	)

	id, err := uuid.NewRandom()
	if err != nil {
		return "", "", err
	}
	at, err := generateToken(id, login, now, expireAccess)
	if err != nil {
		return "", "", err
	}
	rt, err := generateToken(id, login, now, expireRefresh)
	if err != nil {
		return "", "", err
	}
	err = s.stor.AddTokenPair(context.Background(), at, rt)
	if err != nil {
		return "", "", err
	}

	return at.ss, rt.ss, nil
}

var (
	ErrPasswdIncorrect = errors.New("Passwords do not match")
	JWTSecret          = []byte("MyAwesomeSecret")
)

type Token struct {
	id     uuid.UUID
	login  string
	expire time.Time
	ss     string
}

func generateToken(id uuid.UUID, login string, now, expire time.Time) (*Token, error) {
	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		Id:        id.String(),
		NotBefore: now.Unix(),
		ExpiresAt: expire.Unix(),
	})
	ss, err := token.SignedString(JWTSecret)
	if err != nil {
		return nil, err
	}
	return &Token{
		id:     id,
		login:  login,
		expire: expire,
		ss:     ss,
	}, nil
}

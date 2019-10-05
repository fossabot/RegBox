package main

import "go.mongodb.org/mongo-driver/mongo"

type RegBoxService interface {
	Register(string, string) (string, error)
	Authenticate(string, string) (string, string, error)
}

type regBoxService struct {
	stor   *storageService
	crypto *cryptoService
}

func NewRegBoxService(conn *mongo.Client) (*regBoxService, error) {
	var s = NewStorageService(conn)
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

	err = s.stor.AddAccount([]byte(login), hash, salt)
	if err != nil {
		return "", err
	}
	return login, nil
}

func (s regBoxService) Authenticate(login string, password string) (string, string, error) {
	panic("not implemented rpc method")
}

package main

import "errors"

type RegBoxService interface {
	Register(string, string) (string, error)
	Authenticate(string, string) (string, error)
}

type regBoxService struct {
	repo   *accountRepositoryService
	crypto *cryptoService
}

func NewRegBoxService() (*regBoxService, error) {
	r, err := NewAccountRepositoryService()
	if err != nil {
		return nil, err
	}
	c, err := NewCryptoService()
	if err != nil {
		return nil, err
	}
	return &regBoxService{
		repo:   r,
		crypto: c,
	}, nil
}

var (
	ErrLoginUsed = errors.New("Login already used")
)

func (s regBoxService) Register(l string, p string) (string, error) {
	var login = []byte(l)
	var password = []byte(p)

	salt, err := s.crypto.GenerateSalt()
	if err != nil {
		return "", err
	}
	var hash = s.crypto.GenerateHash(password, salt)
	err = s.repo.AddAccount(login, hash, salt)
	if err != nil {
		return "", err
	}
	return l, nil
}

func (s regBoxService) Authenticate(login string, password string) (auth string, err error) {
	panic("not implemented rpc method")
}

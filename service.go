package main

import (
	"context"
	"crypto/rand"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/argon2"
	"gopkg.in/yaml.v2"
)

var (
	ErrFileNotExists = errors.New("File does not exists")
	ErrFileFormat    = errors.New("Unsupported file format")
)

func VerifyConfigPath(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return ErrFileNotExists
	}
	if !strings.HasSuffix(filename, ".yml") && !strings.HasSuffix(filename, ".yaml") {
		return ErrFileFormat
	}
	return nil
}

type RegBox struct {
	Memory      uint32 `yaml:"memory"`
	Iterations  uint32 `yaml:"iterations"`
	Parallelism uint8  `yaml:"parallelism"`
	SaltLength  uint32 `yaml:"salt_length"`
	KeyLength   uint32 `yaml:"key_length"`
	Address     string `yaml:"address"`
}

func NewRegBox(filename string) (service *RegBox, err error) {
	err = VerifyConfigPath(filename)
	if err != nil {
		return nil, err
	}
	rawConfig, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	service = new(RegBox)
	err = yaml.Unmarshal(rawConfig, service)
	if err != nil {
		return nil, err
	}
	return
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
		return &RegisterResponse{Successful: false}, nil
	}
	log.Printf("%x\t%x\n", hash, salt)
	return &RegisterResponse{Successful: true}, nil
}

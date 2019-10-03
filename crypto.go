package main

import (
	"crypto/rand"
	"errors"

	"golang.org/x/crypto/argon2"
)

type cryptoService struct {
	saltLen uint32
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

var (
	DefaultSaltLength  uint32 = 16
	DefaultParallelism uint32 = 1
	DefaultMemory      uint32 = 64 * 1024
	DefaultThreads     uint8  = 4
	DefaultKeyLength   uint32 = 32
)

func NewCryptoService(options ...Argon2OptionFunc) (*cryptoService, error) {
	var s = &cryptoService{
		saltLen: DefaultSaltLength,
		time:    DefaultParallelism,
		memory:  DefaultMemory,
		threads: DefaultThreads,
		keyLen:  DefaultKeyLength,
	}

	for _, option := range options {
		if err := option(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

type Argon2OptionFunc func(*cryptoService) error

// SetTime setup number of iterations over memory
func SetTime(num uint32) Argon2OptionFunc {
	return func(s *cryptoService) error {
		if num > 1 {
			s.time = num
			return nil
		}
		return ErrIterNum
	}
}

// SetMemory setup consumed memory
func SetMemory(amount uint32) Argon2OptionFunc {
	return func(s *cryptoService) error {
		if amount > 0 {
			s.memory = amount
			return nil
		}
		return ErrMemoryAmount
	}
}

func SetThreads(num uint8) Argon2OptionFunc {
	return func(s *cryptoService) error {
		if num > 0 {
			s.threads = num
			return nil
		}
		return ErrThreadsNum
	}
}

func SetKeyLen(len uint32) Argon2OptionFunc {
	return func(s *cryptoService) error {
		if len > 0 {
			s.keyLen = len
			return nil
		}
		return ErrKeyLen
	}
}

func SetSaltLen(len uint32) Argon2OptionFunc {
	return func(s *cryptoService) error {
		if len > 0 {
			s.saltLen = len
			return nil
		}
		return ErrSaltLen
	}
}

var (
	ErrIterNum      = errors.New("Number of iterations must be greater then 1")
	ErrMemoryAmount = errors.New("Memory must be greater then 0")
	ErrThreadsNum   = errors.New("At least 1 thread required")
	ErrKeyLen       = errors.New("Key length must be greater then 0")
	ErrSaltLen      = errors.New("Salt length must be greater then 0")
)

func (s cryptoService) GenerateHash(password, salt []byte) []byte {
	return argon2.IDKey(password, salt, s.time, s.memory, s.threads, s.keyLen)
}

func (s cryptoService) GenerateSalt() ([]byte, error) {
	var b = make([]byte, s.saltLen)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

package main

import (
	"context"
)

type RegBox struct{}

func (*RegBox) Register(ctx context.Context, in *RegisterRequest) (*RegisterResponse, error) {
	return &RegisterResponse{Successful: true}, nil
}

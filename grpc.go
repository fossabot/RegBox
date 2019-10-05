package main

import (
	"context"

	"github.com/Aded175/RegBox/pb"
	grpckit "github.com/go-kit/kit/transport/grpc"
)

type GRPCServer struct {
	register     grpckit.Handler
	authenticate grpckit.Handler
}

func NewGRPCServer(endpoints Endpoints) pb.RegBoxServer {
	return &GRPCServer{
		register: grpckit.NewServer(
			endpoints.Register,
			decodeAcccountRequest,
			encodeAccountResponse,
		),
		authenticate: grpckit.NewServer(
			endpoints.Authenticate,
			decodeAcccountRequest,
			encodeTokenResponse,
		),
	}
}

func (s *GRPCServer) Register(ctx context.Context, req *pb.AcccountRequest) (*pb.AcccountResponse, error) {
	_, resp, err := s.register.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.AcccountResponse), nil
}

func (s *GRPCServer) Authenticate(ctx context.Context, req *pb.AcccountRequest) (*pb.TokenResponse, error) {
	_, resp, err := s.authenticate.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.TokenResponse), nil
}

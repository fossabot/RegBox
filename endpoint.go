package main

import (
	"context"

	"github.com/Aded175/RegBox/pb"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	Register     endpoint.Endpoint
	Authenticate endpoint.Endpoint
}

type AccountRequest struct {
	login    string
	password string
}

type AcccountResponse struct {
	login string
	err   error
}

type TokenResponse struct {
	accessToken  string
	refreshToken string
	err          error
}

func makeRegisterEndpoint(svc RegBoxService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		var request = req.(AccountRequest)
		login, err := svc.Register(request.login, request.password)
		return AcccountResponse{
			login: login,
			err:   err,
		}, nil
	}
}

func decodeAcccountRequest(_ context.Context, r interface{}) (interface{}, error) {
	var req = r.(*pb.AcccountRequest)
	return AccountRequest{
		login:    req.GetLogin(),
		password: req.GetPassword(),
	}, nil
}

func encodeAccountResponse(_ context.Context, r interface{}) (interface{}, error) {
	var req = r.(AcccountResponse)
	return &pb.AcccountResponse{
		Login: req.login,
		Error: err2str(req.err),
	}, nil
}

func makeAuthenticateEndpoint(svc RegBoxService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		var request = req.(AccountRequest)
		at, rt, err := svc.Authenticate(request.login, request.password)
		return TokenResponse{
			accessToken:  at,
			refreshToken: rt,
			err:          err,
		}, nil
	}
}

func encodeTokenResponse(_ context.Context, r interface{}) (interface{}, error) {
	var req = r.(TokenResponse)
	return &pb.TokenResponse{
		AccessToken:  req.accessToken,
		RefreshToken: req.refreshToken,
		Error:        err2str(req.err),
	}, nil
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

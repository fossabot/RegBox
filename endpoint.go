package main

import (
	"context"

	"github.com/Aded175/RegBox/pb"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type Endpoints struct {
	Register     endpoint.Endpoint
	Authenticate endpoint.Endpoint
	Refresh      endpoint.Endpoint
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

type TokenRequest struct {
	uuid         string
	refreshToken string
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

func decodeTokenRequest(_ context.Context, r interface{}) (interface{}, error) {
	var req = r.(*pb.TokenRequest)
	return TokenRequest{
		uuid:         req.GetUuid(),
		refreshToken: req.GetRefreshToken(),
	}, nil
}

func makeRefreshEndpoint(svc RegBoxService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		var request = req.(TokenRequest)
		id, err := uuid.Parse(request.uuid)
		if err != nil {
			return nil, err
		}
		at, rt, err := svc.Refresh(id, request.refreshToken)
		return TokenResponse{
			accessToken:  at,
			refreshToken: rt,
			err:          err,
		}, nil
	}
}

package main

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type loggingMiddleware struct {
	logger *zap.Logger
	next   RegBoxService
}

func (mw loggingMiddleware) Register(login string, password string) (registered string, err error) {
	defer func() {
		mw.logger.Info("logging register",
			zap.String("login", login),
			zap.String("err", err2str(err)),
		)
	}()
	return mw.next.Register(login, password)
}

func (mw loggingMiddleware) Authenticate(login string, password string) (at string, rt string, err error) {
	defer func() {
		mw.logger.Info("logging auth",
			zap.String("login", login),
			zap.String("AT", at),
			zap.String("RT", rt),
			zap.String("err", err2str(err)),
		)
	}()
	return mw.next.Authenticate(login, password)
}

func (mw loggingMiddleware) Refresh(id uuid.UUID, rt string) (string, string, error) {
	return mw.next.Refresh(id, rt)
}

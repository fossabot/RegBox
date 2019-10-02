package main

import (
	"crypto/tls"
	"net"

	"github.com/Aded175/RegBox/pb"
	"github.com/gobuffalo/packr"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func main() {
	var assets = packr.NewBox("./assets")

	logger, _ := zap.NewDevelopment()

	service, err := NewRegBoxServer()
	if err != nil {
		logger.Fatal("Can not create service", zap.String("error", err.Error()))
	}

	listener, err := net.Listen("tcp", service.Address)
	if err != nil {
		logger.Fatal("Can not listen", zap.String("address", service.Address), zap.String("error", err.Error()))
	}

	cert, err := loadCertFromAssets(&assets, "regbox.crt", "regbox.key")
	if err != nil {
		logger.Fatal("Can not load cert from assets", zap.String("cert path", "regbox.crt"), zap.String("key path", "regbox.key"), zap.String("error", err.Error()))
	}

	var creds = credentials.NewServerTLSFromCert(cert)

	var server = grpc.NewServer(grpc.Creds(creds))
	reflection.Register(server)
	pb.RegisterRegBoxServer(server, service)

	err = server.Serve(listener)
	if err != nil {
		logger.Fatal("Can not serve", zap.String("error", err.Error()))
	}
}

func loadCertFromAssets(assets *packr.Box, certPath, keyPath string) (*tls.Certificate, error) {
	crt, err := assets.Find(certPath)
	if err != nil {
		return nil, err
	}
	key, err := assets.Find(keyPath)
	if err != nil {
		return nil, err
	}
	cert, err := tls.X509KeyPair(crt, key)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

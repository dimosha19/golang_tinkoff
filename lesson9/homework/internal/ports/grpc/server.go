package grpc

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"homework9/internal/app"
)

type Service struct {
	a      app.App
	Server *grpc.Server
}

func NewService(a app.App) *Service {
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(SimpleLogger,
		recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler))))
	grpcService := &Service{a, grpcServer}
	RegisterAdServiceServer(grpcServer, grpcService)

	return grpcService
}

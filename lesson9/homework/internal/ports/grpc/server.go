package grpc

import (
	"google.golang.org/grpc"
	"homework9/internal/app"
)

type GRPCService struct {
	a      app.App
	Server *grpc.Server
}

func NewService(a app.App) *GRPCService {
	grpcServer := grpc.NewServer()
	grpcService := &GRPCService{a, grpcServer}
	RegisterAdServiceServer(grpcServer, grpcService)

	return grpcService
}

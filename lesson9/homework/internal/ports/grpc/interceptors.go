package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func SimpleLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println(info.FullMethod)

	return handler(ctx, req)
}

func grpcPanicRecoveryHandler(p any) error {
	log.Println(p)
	return status.Error(codes.Unknown, fmt.Sprintf("{%s} gives an unexpected error", p))
}

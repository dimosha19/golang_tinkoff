package grpc

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

func ExampleUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println(info.FullMethod)

	return handler(ctx, req)
}

func grpcPanicRecoveryHandler(p any) error {
	log.Println(p)
	return errors.New(fmt.Sprintf("%s gives panic", p))
}

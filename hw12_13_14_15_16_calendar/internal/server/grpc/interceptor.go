package internalgrpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

func unaryInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		logger.Info(fmt.Sprintf("Request: %s", info.FullMethod))

		return handler(ctx, req)
	}
}

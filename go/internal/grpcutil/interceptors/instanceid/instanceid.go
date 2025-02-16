package instanceid

import (
	"context"

	"github.com/ahmedalhulaibi/loggy"
	"google.golang.org/grpc"
)

const ContextKey = "instance_id"

// InstanceIdUnaryServerInterceptor returns a new unary server interceptors that injects the instance id from the context.
func InstanceIdUnaryServerInterceptor(logger *loggy.Logger, instanceID string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx, _ = logger.With(ctx, ContextKey, instanceID)
		return handler(ctx, req)
	}
}

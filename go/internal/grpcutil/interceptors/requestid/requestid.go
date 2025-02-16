package requestid

import (
	"context"

	"github.com/ahmedalhulaibi/loggy"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const ContextKey = "request_id"

// RequestIdUnaryServerInterceptor returns a new unary server interceptors that injects a request id into the context.
func RequestIdUnaryServerInterceptor(logger *loggy.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var requestId string

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			requestIdMd := md.Get(ContextKey)

			if len(requestIdMd) == 0 {
				requestId = uuid.NewString()
			} else {
				requestId = requestIdMd[0]
			}
		}

		ctx, _ = logger.With(ctx, ContextKey, requestId)

		return handler(ctx, req)
	}
}

package userid

import (
	"context"

	"github.com/ahmedalhulaibi/loggy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const ContextKey = "user_id"

// UserIdUnaryServerInterceptor returns a new unary server interceptor that extracts the user id from the context.
func UserIdUnaryServerInterceptor(logger *loggy.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			userIdMd := md.Get(ContextKey)

			if len(userIdMd) != 0 {
				userId := userIdMd[0]
				ctx, _ = logger.With(ctx, ContextKey, userId)
			}
		}
		return handler(ctx, req)
	}
}

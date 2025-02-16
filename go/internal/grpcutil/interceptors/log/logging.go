package grpcutil

import (
	"context"
	"time"

	"github.com/ahmedalhulaibi/loggy"
	"google.golang.org/grpc"
)

// LoggerUnaryServerInterceptor returns a new unary server interceptors that logs
func LoggerUnaryServerInterceptor(logger *loggy.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()

		logger.Infow(ctx, "", "method", info.FullMethod, "start", startTime)
		resp, err := handler(ctx, req)

		endTime := time.Now()
		durationVal := endTime.Sub(startTime)

		logger.Infow(ctx, "", "method", info.FullMethod, "start", startTime, "end", endTime, "duration", durationVal)
		return resp, err
	}
}

// LoggerStreamServerInterceptor returns a new stream server interceptors that logs
func LoggerStreamServerInterceptor(logger *loggy.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := time.Now()
		logger.Infow(stream.Context(), "", "method", info.FullMethod, "start", startTime)
		err := handler(srv, stream)

		endTime := time.Now()
		durationVal := endTime.Sub(startTime)

		logger.Infow(stream.Context(), "", "method", info.FullMethod, "start", startTime, "end", endTime, "duration", durationVal)
		return err
	}
}

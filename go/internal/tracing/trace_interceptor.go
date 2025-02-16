package tracing

import (
	"context"

	"github.com/ahmedalhulaibi/hello-world-api/internal/grpcutil/interceptors/instanceid"
	"github.com/ahmedalhulaibi/hello-world-api/internal/grpcutil/interceptors/requestid"
	"github.com/ahmedalhulaibi/hello-world-api/internal/grpcutil/interceptors/userid"
	"github.com/ahmedalhulaibi/loggy"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// NewOpenCensusTraceInterceptor creates a new OpenCensusTraceInterceptor
func NewOpenCensusTraceInterceptor(logger *loggy.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var span *trace.Span

		ctx, span = trace.StartSpan(
			ctx,
			info.FullMethod,
			trace.WithSampler(trace.AlwaysSample()),
		)

		// Extract request id from metadata
		if reqid, ok := extractRequestID(ctx); ok {
			span.AddAttributes(trace.StringAttribute(requestid.ContextKey, reqid))
		}

		// Extract user id from metadata
		if userID, ok := extractUserID(ctx); ok {
			span.AddAttributes(trace.StringAttribute(userid.ContextKey, userID))
		}

		// Extract instance id from metadata
		if instanceID, ok := extractInstanceID(ctx); ok {
			span.AddAttributes(trace.StringAttribute(instanceid.ContextKey, instanceID))
		}

		ctx, _ = logger.With(ctx, "trace_id", span.SpanContext().TraceID.String())
		defer span.End()
		return handler(ctx, req)
	}
}

func extractRequestID(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	reqid := getValueFromIncomingContextMetadata(md, requestid.ContextKey)
	return reqid, reqid != ""
}

func extractUserID(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	userID := getValueFromIncomingContextMetadata(md, userid.ContextKey)
	return userID, userID != ""
}

func extractInstanceID(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	instanceID := getValueFromIncomingContextMetadata(md, instanceid.ContextKey)
	return instanceID, instanceID != ""
}

func getValueFromIncomingContextMetadata(md metadata.MD, key string) string {
	mdValue := md.Get(key)

	if len(mdValue) != 0 {
		val := mdValue[0]
		return val
	}

	return ""
}

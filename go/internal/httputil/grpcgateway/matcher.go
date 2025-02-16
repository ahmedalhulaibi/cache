package httputilgrpcgateway

import (
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/ahmedalhulaibi/hello-world-api/internal/grpcutil/interceptors/requestid"
	"github.com/ahmedalhulaibi/hello-world-api/internal/grpcutil/interceptors/userid"
	"github.com/ahmedalhulaibi/hello-world-api/internal/httputil"
	"github.com/ahmedalhulaibi/hello-world-api/internal/tracing"
)

func CustomMatcher(key string) (string, bool) {
	switch strings.ToLower(key) {
	case strings.ToLower(httputil.XUserUUID):
		return userid.ContextKey, true
	case strings.ToLower(httputil.XRequestID):
		return requestid.ContextKey, true
		// Match B3 headers
	case strings.ToLower(tracing.TraceIDHeader):
		return tracing.B3TraceID, true
	case strings.ToLower(tracing.SpanIDHeader):
		return tracing.B3SpanID, true
	case strings.ToLower(tracing.SampledHeader):
		return tracing.B3Sampled, true
	}

	return runtime.DefaultHeaderMatcher(key)
}

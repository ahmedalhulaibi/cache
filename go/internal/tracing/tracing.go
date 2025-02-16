package tracing

import (
	"go.opencensus.io/plugin/ochttp/propagation/b3"
)

const (
	B3TraceID     = "b3-traceid"
	B3SpanID      = "b3-spanid"
	B3Sampled     = "b3-sampled"
	TraceIDHeader = b3.TraceIDHeader
	SpanIDHeader  = b3.SpanIDHeader
	SampledHeader = b3.SampledHeader
)

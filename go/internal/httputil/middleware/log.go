package middleware

import (
	"net/http"
	"time"

	"github.com/ahmedalhulaibi/loggy"
)

// Logger middleware.
type Logger struct {
	h      http.Handler
	logger *loggy.Logger
}

// SetLogger sets the logger to `log`. If you have used logger.New(), you can use this to set your
// logger. Alternatively, if you already have your log.Logger, then you can just call logger.NewLogger() directly.
func (l *Logger) SetLogger(logger *loggy.Logger) {
	l.logger = logger
}

// wrapper to capture status.
type wrapper struct {
	http.ResponseWriter
	written int
	status  int
}

// capture status.
func (w *wrapper) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// capture written bytes.
func (w *wrapper) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.written += n
	return n, err
}

// NewLogger logger middleware with the given log.Logger.
func NewLogger(logger *loggy.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &Logger{
			logger: logger,
			h:      h,
		}
	}
}

// ServeHTTP implementation.
func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	res := &wrapper{w, 0, 200}

	// get the context since we'll use it a few times
	ctx := r.Context()

	// continue to the next middleware
	l.h.ServeHTTP(res, r.WithContext(ctx))

	// log the request.end
	l.logger.Infow(ctx,
		"request ended",
		"method", r.Method,
		"uri", r.RequestURI,
		"headers", r.Header,
		"status", res.status,
		"size", res.written,
		"duration", time.Since(start),
		"evt", "request.end",
	)
}

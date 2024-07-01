package middleware

import (
	"net/http"
	"time"
)

type Logging struct {
	Handler http.Handler
}

func NewLoggingMiddleware(handler http.Handler) *Logging {
	return &Logging{Handler: handler}
}

func (l *Logging) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	customWriter := &responseWriter{ResponseWriter: w, status: http.StatusOK}

	l.Handler.ServeHTTP(customWriter, r)

	duration := time.Since(start).Seconds()

	logEntry := LogEntry{
		Time:    time.Now().Format(time.RFC3339),
		Level:   "info",
		Message: "request completed",
		Request: LogInfo{
			Method:       r.Method,
			URL:          r.URL.String(),
			TimeStamp:    start.Format(time.RFC3339),
			ResponseTime: duration,
			Status:       customWriter.status,
		},
	}

	Logger.Info("request completed", "", logEntry.LogValue())
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

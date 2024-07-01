package middleware

import "log/slog"

var Logger *slog.Logger

type LogInfo struct {
	Method       string
	URL          string
	TimeStamp    string
	ResponseTime float64
	Status       int
}

type LogEntry struct {
	Time    string
	Level   string
	Message string
	Request LogInfo
}

func (le LogEntry) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("time", le.Time),
		slog.String("level", le.Level),
		slog.String("message", le.Message),
		slog.Group("request",
			slog.String("method", le.Request.Method),
			slog.String("url", le.Request.URL),
			slog.String("timestamp", le.Request.TimeStamp),
			slog.Float64("responseTime", le.Request.ResponseTime),
			slog.Int("status", le.Request.Status),
		))
}

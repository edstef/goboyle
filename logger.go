package main

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/exp/slog"

	"github.com/go-chi/chi/v5/middleware"
)

func NewStructuredLogger(handler slog.Handler) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{Logger: handler})
}

type StructuredLogger struct {
	Logger slog.Handler
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	var logFields []slog.Attr
	logFields = append(logFields, slog.String("ts", time.Now().UTC().Format(time.RFC1123)))

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields = append(logFields, slog.String("req_id", reqID))
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	handler := l.Logger.WithAttrs(append(logFields,
		// slog.String("http_scheme", scheme),
		// slog.String("http_proto", r.Proto),
		slog.String("HTTP_METHOD", r.Method),
		// slog.String("remote_addr", r.RemoteAddr),
		slog.String("USER_AGENT", r.UserAgent()),
		slog.String("URI", fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI))))

	entry := StructuredLoggerEntry{Logger: slog.New(handler)}

	entry.Logger.LogAttrs(r.Context(), slog.LevelInfo, "INCOMING REQUEST")

	return &entry
}

type StructuredLoggerEntry struct {
	Logger *slog.Logger
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger.LogAttrs(nil, slog.LevelInfo, "OUTGOING RESPONSE",
		slog.Int("RESPONSE_CODE", status),
		slog.Int("RESPONSE_LENGTH", bytes),
		slog.Float64("RESPONSE_DURATION", float64(elapsed.Nanoseconds())/1000000.0),
	)
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger.LogAttrs(nil, slog.LevelInfo, "",
		slog.String("stack", string(stack)),
		slog.String("panic", fmt.Sprintf("%+v", v)),
	)
}

func GetLogEntry(r *http.Request) *slog.Logger {
	entry := middleware.GetLogEntry(r).(*StructuredLoggerEntry)
	return entry.Logger
}

func LogEntrySetField(r *http.Request, key string, value interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.Logger = entry.Logger.With(key, value)
	}
}

func LogEntrySetFields(r *http.Request, fields map[string]interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		for k, v := range fields {
			entry.Logger = entry.Logger.With(k, v)
		}
	}
}

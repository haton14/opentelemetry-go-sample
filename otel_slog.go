package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"

	slogotel "github.com/remychantenay/slog-otel"
)

func SetLogger() {
	logger := slog.New(
		slogotel.OtelHandler{
			Next: NewStackTraceLogHandler(
				slog.NewJSONHandler(os.Stderr, nil),
			),
		},
	)
	slog.SetDefault(logger)
}

type StackTraceLogHandler struct {
	slog.Handler
}

func NewStackTraceLogHandler(h slog.Handler) StackTraceLogHandler {
	return StackTraceLogHandler{h}
}

func (h StackTraceLogHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level != slog.LevelError {
		return h.Handler.Handle(ctx, r)
	}
	// fileの行数を追加
	pt, file, line, _ := runtime.Caller(3)
	funcName := runtime.FuncForPC(pt).Name()
	r.AddAttrs(slog.String("call", fmt.Sprintf("%s:%d %s", file, line, funcName)))
	return h.Handler.Handle(ctx, r)
}

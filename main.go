package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
)

func SetLogger() {
	logger := slog.New(
		NewStackTraceLogHandler(
			slog.NewJSONHandler(os.Stderr, nil),
		),
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

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	SetLogger()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		_ = r.Context()
		slog.Info("GET /")
		w.WriteHeader(http.StatusNoContent)
	})
	srv := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: mux,
	}

	go func() {
		slog.Info("ListenAndServe", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("ListenAndServe", "err", err)
		}
	}()
	<-ctx.Done()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Shutdown", "err", err)
	}
}

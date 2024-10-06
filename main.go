package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	exp, err := newExporter(ctx)
	if err != nil {
		slog.Error("newExporter", "err", err)
	}
	tp := newTraceProvider(exp)
	defer func() {
		tp.Shutdown(ctx)
	}()

	otel.SetTracerProvider(tp)
	tracer = tp.Tracer("example.io/package/name")

	SetLogger()

	mux := http.NewServeMux()
	// wrap the handler function with otelhttp.WithRouteTag
	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		// Configure the "http.route" for the HTTP instrumentation.
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		mux.Handle(pattern, handler)
	}
	handleFunc("GET /unko", func(w http.ResponseWriter, r *http.Request) {
		_ = r.Context()
		slog.Info("GET /unko")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("unko\n"))
	})
	srv := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: otelhttp.NewHandler(mux, "/"),
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

package common

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

type TraceHandler struct {
	slog.Handler
}

func (h *TraceHandler) Handle(ctx context.Context, r slog.Record) error {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		r.AddAttrs(
			slog.String("trace_id", span.SpanContext().TraceID().String()),
			slog.String("span_id", span.SpanContext().SpanID().String()),
		)
	}
	return h.Handler.Handle(ctx, r)
}

func NewLogger(serviceName string) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		AddSource: true,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	traceHandler := &TraceHandler{Handler: handler}
	
	logger := slog.New(traceHandler).With(slog.String("service", serviceName))
	slog.SetDefault(logger)
	
	return logger
}

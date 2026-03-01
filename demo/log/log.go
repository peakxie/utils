package log

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	otelzap "github.com/uptrace/opentelemetry-go-extra/otelzap"
)

// Init initializes the global otelzap logger.
// Call this once at application startup after the TracerProvider is set.
func Init() {
	logger := otelzap.New(zap.Must(zap.NewProduction()))
	otelzap.ReplaceGlobals(logger)
}

// Ctx returns a SugaredLogger with trace context.
// If the context carries an active span with a trace ID, it is added as a field.
func Ctx(ctx context.Context) otelzap.SugaredLoggerWithCtx {
	s := otelzap.S()
	if span := trace.SpanFromContext(ctx); span.SpanContext().HasTraceID() {
		s = s.With(zap.String("traceid", span.SpanContext().TraceID().String()))
	}
	return s.Ctx(ctx)
}

// CtxInfof logs an info message with traceid from context.
func CtxInfof(ctx context.Context, format string, args ...any) {
	if span := trace.SpanFromContext(ctx); span.SpanContext().HasTraceID() {
		format = fmt.Sprintf("[traceid:%s] %s", span.SpanContext().TraceID().String(), format)
	}
	otelzap.S().Ctx(ctx).Infof(format, args...)
}

// CtxErrorf logs an error message with traceid from context.
func CtxErrorf(ctx context.Context, format string, args ...any) {
	if span := trace.SpanFromContext(ctx); span.SpanContext().HasTraceID() {
		format = fmt.Sprintf("[traceid:%s] %s", span.SpanContext().TraceID().String(), format)
	}
	otelzap.S().Ctx(ctx).Errorf(format, args...)
}

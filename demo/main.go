package main

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/peakxie/utils/ginmw"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	oteltrace "go.opentelemetry.io/otel/trace"

	dlog "github.com/peakxie/utils/demo/log"
)

func init() {
	ginmw.SetLogFunc(func(c *gin.Context, format string, args ...any) {
		dlog.CtxInfof(c.Request.Context(), format, args...)
	})
}

// initTracer sets up the OpenTelemetry TracerProvider with a stdout exporter.
func initTracer() func() {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		panic(err)
	}
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	otel.SetTracerProvider(tp)
	return func() { tp.Shutdown(context.Background()) }
}

// Trace is a gin middleware that starts an OpenTelemetry span for each request.
func Trace() gin.HandlerFunc {
	tracer := otel.Tracer("demo")
	return func(c *gin.Context) {
		ctx, span := tracer.Start(c.Request.Context(), c.Request.Method+" "+c.FullPath(),
			oteltrace.WithAttributes(
				semconv.HTTPRequestMethodKey.String(c.Request.Method),
				semconv.URLPath(c.Request.URL.Path),
			),
		)
		defer span.End()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// RequestID is a gin middleware that sets a request_id in the context
// and injects it into the JSON response body.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader("X-Request-Id")
		if rid == "" {
			rid = uuid.New().String()
		}
		c.Set("request_id", rid)

		// Intercept response body to inject request_id.
		bw := &bodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = bw
		c.Next()

		// Inject request_id into JSON response.
		raw := bw.body.Bytes()
		var m map[string]any
		if json.Unmarshal(raw, &m) == nil {
			m["request_id"] = rid
			raw, _ = json.Marshal(m)
		}
		bw.ResponseWriter.Write(raw)
	}
}

// bodyWriter buffers the response body so we can modify it before sending.
type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

// HelloReq is the request for the /hello endpoint.
type HelloReq struct {
	Name string `json:"name" binding:"required"`
}

// HelloRsp is the response for the /hello endpoint.
type HelloRsp struct {
	Greeting string `json:"greeting"`
}

func hello(_ *gin.Context, req *HelloReq) (*HelloRsp, error) {
	return &HelloRsp{Greeting: "hello " + req.Name}, nil
}

func main() {
	shutdown := initTracer()
	defer shutdown()

	dlog.Init()

	r := gin.Default()
	r.Use(Trace())
	r.Use(RequestID())
	r.POST("/hello", ginmw.WrapperH(hello))
	r.Run(":8080")
}

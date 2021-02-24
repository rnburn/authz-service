package auth

import (
  "fmt"
  "log"
  "net/http"
  "bytes"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
)

type responseCaptureServer struct {
	tracer trace.Tracer
}

func NewResponseCaptureServer() *responseCaptureServer {
	return &responseCaptureServer{
		tracer: otel.Tracer("response-capture"),
	}
}

func (server *responseCaptureServer) ServeHTTP(responseWriter http.ResponseWriter, request * http.Request) {
	propagator := otel.GetTextMapPropagator()
  ctx := propagator.Extract(request.Context(), request.Header)
  ctx, span := server.tracer.Start(ctx, "response capture")
  buffer := new(bytes.Buffer)
  buffer.ReadFrom(request.Body)
  span.SetAttributes(
			label.String("http.response.body", buffer.String()))
  fmt.Printf("incoming http request")
  fmt.Fprintf(responseWriter, "nod\n")
  defer span.End()
}


func (server *responseCaptureServer) Run() {
  port := 8080
  fmt.Printf("Listening for responses on: %d\n", port)
  http.Handle("/", server)
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

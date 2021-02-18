package auth

import (
  "fmt"
  "log"
  "context"
  "net/http"

	"go.opentelemetry.io/otel"
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
  ctx := context.Background()   
  ctx, span := server.tracer.Start(ctx, "response capture")
  defer span.End()
}


func (server *responseCaptureServer) Run() {
  port := 8080
  fmt.Printf("Listening for responses on: %d\n", port)
  http.Handle("/response", server)
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

package auth

import (
	"context"
	"fmt"
  "time"

	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
  "github.com/golang/protobuf/ptypes"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
)

type server struct {
	tracer trace.Tracer
}

var _ envoy_service_auth_v3.AuthorizationServer = &server{}

// New creates a new authorization server.
func New() envoy_service_auth_v3.AuthorizationServer {
	return &server{
		tracer: otel.Tracer("authz-service"),
	}
}

func setSpanAttributes(span trace.Span,
	req *envoy_service_auth_v3.AttributeContext_HttpRequest) {
  span.SetAttributes(label.String("http.url", req.Method))
	for key, value := range req.Headers {
		span.SetAttributes(
			label.String(fmt.Sprintf("http.request.header.%s", key), value))
	}
}

// Check implements authorization's Check interface which performs authorization check based on the
// attributes associated with the incoming request.
func (s *server) Check(
	ctx context.Context,
	req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
  timestamp, err := ptypes.Timestamp(req.Attributes.Request.Time)
  if err != nil {
    timestamp = time.Now()
  }
	http := req.Attributes.Request.Http
	ctx, span := s.tracer.Start(ctx, http.Method, trace.WithTimestamp(timestamp))
	setSpanAttributes(span, http)
	defer span.End()
	return &envoy_service_auth_v3.CheckResponse{
		Status: &status.Status{
			Code: int32(code.Code_OK),
		},
	}, nil
}

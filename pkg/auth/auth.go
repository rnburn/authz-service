package auth

import (
	"context"

	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"

	"go.opentelemetry.io/otel"
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

// Check implements authorization's Check interface which performs authorization check based on the
// attributes associated with the incoming request.
func (s *server) Check(
	ctx context.Context,
	req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
	ctx, span := s.tracer.Start(ctx, "authz-request")
	defer span.End()
	return &envoy_service_auth_v3.CheckResponse{
		Status: &status.Status{
			Code: int32(code.Code_OK),
		},
	}, nil
}

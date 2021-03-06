package auth

import (
	"context"

	envoy_config_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

type serverV3 struct {
	tracer trace.Tracer
}

var _ envoy_service_auth_v3.AuthorizationServer = &serverV3{}

// New creates a new authorization serverV3.
func NewServerV3() envoy_service_auth_v3.AuthorizationServer {
	return &serverV3{
		tracer: otel.Tracer("authz-service"),
	}
}

func setRequestBodyV3(span trace.Span, req *envoy_service_auth_v3.AttributeContext_HttpRequest) {
	if len(req.Body) == 0 && len(req.RawBody) == 0 {
		return
	}
	if !shouldRecordBody(req.Headers) {
		return
	}
	if len(req.Body) > 0 {
		span.SetAttributes(label.String("http.request.body", req.Body))
	} else {
		span.SetAttributes(label.String("http.request.body", string(req.RawBody)))
	}
}

func setSpanAttributesV3(span trace.Span,
	req *envoy_service_auth_v3.AttributeContext_HttpRequest) {
	setRequestBodyV3(span, req)
	span.SetAttributes(label.String("http.url", req.Path))
	setHeaderAnnotations(span, req.Headers)
}

func setPortAttributeV3(span trace.Span, key label.Key, address *envoy_config_v3.SocketAddress) {
	switch address.PortSpecifier.(type) {
	case *envoy_config_v3.SocketAddress_PortValue:
		span.SetAttributes(key.Int(int(address.GetPortValue())))
	case *envoy_config_v3.SocketAddress_NamedPort:
		span.SetAttributes(key.String(address.GetNamedPort()))
	}
}

func setSourcePeerV3(span trace.Span,
	source *envoy_service_auth_v3.AttributeContext_Peer) {
	switch address := source.Address.Address.(type) {
	case *envoy_config_v3.Address_SocketAddress:
		span.SetAttributes(semconv.NetPeerIPKey.String(address.SocketAddress.Address))
		setPortAttributeV3(span, semconv.NetPeerPortKey, address.SocketAddress)
	case *envoy_config_v3.Address_Pipe:
		return
	}
}

// Check implements authorization's Check interface which performs authorization check based on the
// attributes associated with the incoming request.
func (s *serverV3) Check(
	ctx context.Context,
	req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
	http := req.Attributes.Request.Http
	ctx, span := startSpan(ctx, s.tracer, req.Attributes.Request.Time, http.Method)
	setSourcePeerV3(span, req.Attributes.Source)
	setSpanAttributesV3(span, http)
	defer span.End()
	return &envoy_service_auth_v3.CheckResponse{
		Status: &status.Status{
			Code: int32(code.Code_OK),
		},
	}, nil
}

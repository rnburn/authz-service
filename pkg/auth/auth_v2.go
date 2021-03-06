package auth

import (
	"context"

	envoy_config_v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoy_service_auth_v2 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

type serverV2 struct {
	tracer trace.Tracer
}

var _ envoy_service_auth_v2.AuthorizationServer = &serverV2{}

// New creates a new authorization serverV2.
func NewServerV2() envoy_service_auth_v2.AuthorizationServer {
	return &serverV2{
		tracer: otel.Tracer("authz-service"),
	}
}

func setRequestBodyV2(span trace.Span, req *envoy_service_auth_v2.AttributeContext_HttpRequest) {
	if len(req.Body) == 0 {
		return
	}
	if !shouldRecordBody(req.Headers) {
		return
	}
	span.SetAttributes(label.String("http.request.body", req.Body))
}

func setSpanAttributesV2(span trace.Span,
	req *envoy_service_auth_v2.AttributeContext_HttpRequest) {
	setRequestBodyV2(span, req)
	span.SetAttributes(label.String("http.url", req.Path))
	setHeaderAnnotations(span, req.Headers)
}

func setPortAttributeV2(span trace.Span, key label.Key, address *envoy_config_v2.SocketAddress) {
	switch address.PortSpecifier.(type) {
	case *envoy_config_v2.SocketAddress_PortValue:
		span.SetAttributes(key.Int(int(address.GetPortValue())))
	case *envoy_config_v2.SocketAddress_NamedPort:
		span.SetAttributes(key.String(address.GetNamedPort()))
	}
}

func setSourcePeerV2(span trace.Span,
	source *envoy_service_auth_v2.AttributeContext_Peer) {
	switch address := source.Address.Address.(type) {
	case *envoy_config_v2.Address_SocketAddress:
		span.SetAttributes(semconv.NetPeerIPKey.String(address.SocketAddress.Address))
		setPortAttributeV2(span, semconv.NetPeerPortKey, address.SocketAddress)
	case *envoy_config_v2.Address_Pipe:
		return
	}
}

// Check implements authorization's Check interface which performs authorization check based on the
// attributes associated with the incoming request.
func (s *serverV2) Check(
	ctx context.Context,
	req *envoy_service_auth_v2.CheckRequest) (*envoy_service_auth_v2.CheckResponse, error) {
	http := req.Attributes.Request.Http
	ctx, span := startSpan(ctx, s.tracer, req.Attributes.Request.Time, http.Method)
	setSourcePeerV2(span, req.Attributes.Source)
	setSpanAttributesV2(span, http)
	defer span.End()
	return &envoy_service_auth_v2.CheckResponse{
		Status: &status.Status{
			Code: int32(code.Code_OK),
		},
	}, nil
}

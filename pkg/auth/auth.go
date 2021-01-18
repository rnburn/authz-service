package auth

import (
  "os"
	"context"
	"fmt"
  "time"
  "strings"

	envoy_config_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
  "github.com/golang/protobuf/ptypes"
  "google.golang.org/grpc/metadata"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
  "go.opentelemetry.io/otel/semconv"
)

// contentTypeAllowList is the list of allowed content types in lowercase
var contentTypeAllowListLowerCase = []string{
	"application/json",
	"application/x-www-form-urlencoded",
}

type textMapCarrier struct {
  headers map[string][]string
}

func (carrier *textMapCarrier) Get(key string) string {
  values := carrier.headers[key]
  if len(values) > 0 {
    fmt.Printf("Get key: %s\t%s\n", key, values[0])
    return values[0]
  } else {
    return ""
  }
}

func (carrier *textMapCarrier) Set(key string, value string) {
}

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

func shouldRecordBody(content_type string) bool {
  for _, recordableContentType := range contentTypeAllowListLowerCase {
    if strings.Contains(content_type, recordableContentType) {
      return true
    }
  }
  return false
}

func setRequestBody(span trace.Span, req* envoy_service_auth_v3.AttributeContext_HttpRequest) {
  if len(req.Body) == 0 && len(req.RawBody) == 0 {
    return
  }
  content_type := strings.ToLower(req.Headers["Content-Type"])
  if !shouldRecordBody(content_type) {
    return
  }
  if len(req.Body) > 0 {
    span.SetAttributes(label.String("http.request.body", req.Body))
  } else {
    span.SetAttributes(label.String("http.request.body", string(req.RawBody)))
  }
}

func setSpanAttributes(span trace.Span,
	req *envoy_service_auth_v3.AttributeContext_HttpRequest) {
  span.SetAttributes(label.String("http.url", req.Path))
	for key, value := range req.Headers {
		span.SetAttributes(
			label.String(fmt.Sprintf("http.request.header.%s", key), value))
	}
}

func setPortAttribute(span trace.Span, key label.Key, address *envoy_config_v3.SocketAddress) {
  switch address.PortSpecifier.(type) {
    case *envoy_config_v3.SocketAddress_PortValue:
      span.SetAttributes(key.Int(int(address.GetPortValue())))
    case *envoy_config_v3.SocketAddress_NamedPort:
      span.SetAttributes(key.String(address.GetNamedPort()))
  }
}

func setSourcePeer(span trace.Span,
  source *envoy_service_auth_v3.AttributeContext_Peer) {
    switch address := source.Address.Address.(type) {
    case *envoy_config_v3.Address_SocketAddress:
      span.SetAttributes(semconv.NetPeerIPKey.String(address.SocketAddress.Address))
      setPortAttribute(span, semconv.NetPeerPortKey, address.SocketAddress)
    case *envoy_config_v3.Address_Pipe:
      return
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
  propagator := otel.GetTextMapPropagator()
  md, ok := metadata.FromIncomingContext(ctx)
  if ok {
    carrier := textMapCarrier{md}
    ctx = propagator.Extract(ctx, &carrier)
  }
  pctx := trace.RemoteSpanContextFromContext(ctx)
  fmt.Fprintf(os.Stderr, "parent TraceID %s\n", pctx.TraceID.String())
  _ = pctx
	ctx, span := s.tracer.Start(ctx, http.Method, trace.WithTimestamp(timestamp))
  fmt.Fprintf(os.Stderr, "span TraceID %s\n", span.SpanContext().TraceID.String())
  setSourcePeer(span, req.Attributes.Source)
	setSpanAttributes(span, http)
	defer span.End()
	return &envoy_service_auth_v3.CheckResponse{
		Status: &status.Status{
			Code: int32(code.Code_OK),
		},
	}, nil
}

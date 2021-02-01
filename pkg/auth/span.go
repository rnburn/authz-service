package auth

import (
	"context"
  "fmt"
  "time"

	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel"
  "github.com/golang/protobuf/ptypes"
  "google.golang.org/grpc/metadata"
  timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	"go.opentelemetry.io/otel/label"
)

func setHeaderAnnotations(span trace.Span, headers map[string]string) {
	for key, value := range headers {
		span.SetAttributes(
			label.String(fmt.Sprintf("http.request.header.%s", key), value))
	}
}

func startSpan(ctx context.Context, tracer trace.Tracer, requestTimestamp *timestamppb.Timestamp, method string) (context.Context, trace.Span) {
  timestamp, err := ptypes.Timestamp(requestTimestamp)
  if err != nil {
    timestamp = time.Now()
  }
  propagator := otel.GetTextMapPropagator()
  md, ok := metadata.FromIncomingContext(ctx)
  if ok {
    carrier := textMapCarrier{md}
    ctx = propagator.Extract(ctx, &carrier)
  }
	return tracer.Start(ctx, method, trace.WithTimestamp(timestamp))
}

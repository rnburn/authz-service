package auth

import (
  "fmt"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/label"
)

func setHeaderAnnotations(span trace.Span, headers map[string]string) {
	for key, value := range headers {
		span.SetAttributes(
			label.String(fmt.Sprintf("http.request.header.%s", key), value))
	}
}

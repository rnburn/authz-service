module github.com/rnburn/authz-service

go 1.15

require (
	github.com/envoyproxy/envoy/examples/ext_authz/auth/grpc-service v0.0.0-20210106224028-4cb14ea2da6e
	github.com/envoyproxy/go-control-plane v0.9.8
	github.com/golang/protobuf v1.4.3
	github.com/hypertrace/goagent v0.0.0-20201216150242-e980621edb2e
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/contrib/propagators v0.15.0
	go.opentelemetry.io/otel v0.15.0
	google.golang.org/genproto v0.0.0-20200806141610-86f49bd18e98
	google.golang.org/grpc v1.34.0
	google.golang.org/protobuf v1.25.0
)

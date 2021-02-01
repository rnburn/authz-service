package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	envoy_service_auth_v2 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/grpc"

	auth "github.com/rnburn/authz-service/pkg/auth"

	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/instrumentation/hypertrace"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func main() {
	hypertraceConfig := config.Load()
	hypertraceConfig.ServiceName = config.String("authz-service")
	shutdown := hypertrace.Init(hypertraceConfig)
	defer shutdown()

	port := flag.Int("port", 9001, "gRPC port")

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen to %d: %v", *port, err)
	}

	authzConfig := auth.LoadConfig()
	if authzConfig.PropagationMode == auth.B3PropagationMode {
		otel.SetTextMapPropagator(b3.B3{})
	} else if authzConfig.PropagationMode == auth.TraceContextPropagationMode {
		otel.SetTextMapPropagator(propagation.TraceContext{})
	}

	gs := grpc.NewServer()

	envoy_service_auth_v3.RegisterAuthorizationServer(gs, auth.NewServerV3())
	envoy_service_auth_v2.RegisterAuthorizationServer(gs, auth.NewServerV2())

	log.Printf("starting gRPC server on: %d\n", authzConfig.Port)

	gs.Serve(lis)
}

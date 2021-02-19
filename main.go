package main

import (
	"fmt"
	"log"
	"net"
  // "os"
  // "os/signal"
  // "syscall"

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

  // Ambassador only works with b3 propagation, so we support configuring it.
  //    See https://datawire-oss.slack.com/archives/CAULN7S76/p1611811639150700
  // This might be removed in the future if configurable propagation is supported directly
  // in the hypertrace agents.
	authzConfig := auth.LoadConfig()
	if authzConfig.PropagationMode == auth.B3PropagationMode {
		otel.SetTextMapPropagator(b3.B3{})
	} else if authzConfig.PropagationMode == auth.TraceContextPropagationMode {
		otel.SetTextMapPropagator(propagation.TraceContext{})
	}


	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", authzConfig.Port))
	if err != nil {
		log.Fatalf("failed to listen to %d: %v", authzConfig.Port, err)
	}


  // Support both Envoy's v2 and v3 protcols
	gs := grpc.NewServer()
	envoy_service_auth_v3.RegisterAuthorizationServer(gs, auth.NewServerV3())
	envoy_service_auth_v2.RegisterAuthorizationServer(gs, auth.NewServerV2())

	log.Printf("starting gRPC server on: %d\n", authzConfig.Port)
	gs.Serve(lis)

  /*
  // Run a server for capturing responses
  responseServer := auth.NewResponseCaptureServer()
  go responseServer.Run()

  // from https://www.alexsears.com/2019/10/fun-with-concurrency-in-golang/
  signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
  */
}

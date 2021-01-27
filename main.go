package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoy_service_auth_v2 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	"google.golang.org/grpc"

	auth "github.com/rnburn/authz-service/pkg/auth"

	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/instrumentation/hypertrace"
)

func main() {
	cfg := config.Load()
	cfg.ServiceName = config.String("authz-service")
	shutdown := hypertrace.Init(cfg)
	defer shutdown()

	port := flag.Int("port", 9001, "gRPC port")

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen to %d: %v", *port, err)
	}

	gs := grpc.NewServer()

  if true {
	  envoy_service_auth_v3.RegisterAuthorizationServer(gs, auth.NewServerV3())
  } else {
    envoy_service_auth_v2.RegisterAuthorizationServer(gs, auth.NewServerV2())
  }

	log.Printf("starting gRPC server on: %d\n", *port)

	gs.Serve(lis)
}

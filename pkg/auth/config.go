package auth

import (
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	TraceContextPropagationMode = 1
	B3PropagationMode           = 2
	DefaultAuthzPropagationMode = TraceContextPropagationMode
	DefaultAuthzPort            = 9001
)

func getAuthzPort() int {
	value := os.Getenv("HT_AUTHZ_PORT")
	if value == "" {
		return DefaultAuthzPort
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Invalid port %s", value)
	}
	return intValue
}

func getAuthzPropagationMode() int {
	value := os.Getenv("HT_AUTHZ_PROPAGATION_MODE")
	if value == "" {
		return DefaultAuthzPropagationMode
	}
	value = strings.ToLower(value)
	if value == "trace-context" {
		return TraceContextPropagationMode
	}
	if value != "b3" {
		log.Fatalf("Invalid propagation mode %s", value)
	}
	return B3PropagationMode
}

type Config struct {
	Port            int
	PropagationMode int
}

func LoadConfig() Config {
	return Config{
		Port:            getAuthzPort(),
		PropagationMode: getAuthzPropagationMode(),
	}
}

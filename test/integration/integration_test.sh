#!/bin/bash

ENVOY=third_party/envoy
SERVICE=./authz-service_/authz-service
CONFIG=test/integration/config.yaml

$ENVOY --help
$SERVICE
$ENVOY -c $CONFIG --service-cluster front-proxy

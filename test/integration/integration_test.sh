#!/bin/bash

CONFIG=test/integration/config.yaml

envoy --help
envoy -c $CONFIG --service-cluster front-proxy
# $SERVICE

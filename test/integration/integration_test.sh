#!/bin/bash

TEST_DIR=test/integration

# Start upstream server
python3 $TEST_DIR/upstream/server.py &
UPSTREAM_PID=$!

# Start authz service
./authz-service &
AUTHZ_PID=$!

# Start envoy
envoy -c $TEST_DIR/config.yaml --service-cluster front-proxy &
ENVOY_PID=$!

# Shutdown
sleep 5
kill $ENVOY_PID
kill $UPSTREAM_PID
kill $AUTHZ_PID

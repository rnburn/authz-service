#!/bin/bash

ENVOY=third_party/envoy
SERVICE=./authz-service_/authz-service

$ENVOY --help
$SERVICE

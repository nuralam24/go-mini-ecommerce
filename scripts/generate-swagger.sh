#!/bin/bash

set -e

# Generate Swagger documentation into Go package `swagger`
# Requires: go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/server/main.go -o swagger

echo "Swagger documentation generated in swagger/ package"

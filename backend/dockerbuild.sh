#!/bin/bash

set -e

# Variables
BINARY_NAME="app"
DOCKER_IMAGE_NAME="send-blue-backend"
DOCKER_IMAGE_TAG="latest"

# Compile Go binary
echo "Compiling Go code..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $BINARY_NAME .

# Build Docker image
echo "Building Docker image..."
docker build -t $DOCKER_IMAGE_NAME:$DOCKER_IMAGE_TAG .

echo "Build complete! Image: $DOCKER_IMAGE_NAME:$DOCKER_IMAGE_TAG"
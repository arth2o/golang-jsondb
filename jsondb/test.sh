#!/bin/bash

# Kill any existing server instances
pkill -f "jsondb/bin/server"

# Clean test cache
go clean -testcache

# Run all tests
go test -v ./...

# Run specific package tests
go test -v ./internal/config/...
go test -v ./internal/server/...
go test -v ./internal/engine/...
go test -v ./internal/encryption/...
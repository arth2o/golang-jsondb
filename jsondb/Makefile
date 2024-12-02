.PHONY: test build clean run test-coverage test-verbose

build:
	go build -o bin/server cmd/server/main.go

test:
	go clean -testcache
	go mod tidy
	go test -v ./internal/config/...
	go test -v ./internal/server/...
	go test -v ./internal/storage/...
	go test -v ./cmd/server/...

clean:
	rm -rf bin/
	go clean -testcache

run:
	go run cmd/server/main.go

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-verbose:
	go test -v -race ./...
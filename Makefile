all: build test

build:
	@echo "Building..."
	
	@go build -o proxy cmd/proxy/proxy.go

run:
	@go run cmd/proxy/proxy.go

test:
	@echo "Testing..."
	
	@go test ./... -v
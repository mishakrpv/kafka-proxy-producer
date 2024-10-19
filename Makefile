all: build test

build:
	@echo "Building..."
	
	@go build -o main cmd/proxy/main.go

run:
	@go run cmd/proxy/main.go

test:
	@echo "Testing..."
	@go test ./... -v
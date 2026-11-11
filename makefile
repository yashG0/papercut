.PHONY: run build test tidy lint clean

# Default port
PORT ?= 8080

# Run the development server
run:
	go run ./cmd/server/main.go

# Compile binary to bin/
build:
	mkdir -p bin
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/papercut ./cmd/server/main.go

# Run tests
test:
	go test -v -race ./...

# Tidy dependencies
tidy:
	go mod tidy

# Clean up build artifacts
clean:
	rm -rf bin/

.PHONY: build test clean

build:
	@echo "Building devdock (dev version)..."
	go build -ldflags "-X main.Version=dev" -o bin/devdock ./cmd/devdock

test:
	go test ./...

clean:
	rm -rf bin/

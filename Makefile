BINARY_NAME=cortex
VERSION?=0.5.0

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_UNIX=$(BINARY_NAME)_unix

# Build directories
BUILD_DIR=./bin
CMD_DIR=./cmd/cortex

.PHONY: all build clean test coverage deps lint run dev install-tools install

install:
	go install ./cmd/cortex

# Rest of the Makefile remains unchanged...
build: $(BUILD_DIR)
	CGO_ENABLED=1 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(CMD_DIR)

test:
	$(GOTEST) -v ./...

coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

deps:
	$(GOMOD) download
	$(GOMOD) verify

lint:
	golangci-lint run

dev:
	CGO_ENABLED=1 air -c $(PWD)/.air.toml

install-tools:
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

build-linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_UNIX) -v $(CMD_DIR)

docker-build:
	docker build -t $(BINARY_NAME):$(VERSION) .
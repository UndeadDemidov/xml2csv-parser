# ==============================================================================
# Arguments passing to Makefile commands
GO_INSTALLED := $(shell which go)
MG_INSTALLED := $(shell which mockgen 2> /dev/null)
SS_INSTALLED := $(shell which staticcheck 2> /dev/null)
GL_INSTALLED := $(shell which golint 2> /dev/null)

PROJECT_NAME=$(notdir $(shell pwd))

# ==============================================================================
# Install commands
install-tools:
	@echo Checking tools are installed...
ifndef SS_INSTALLED
	@echo Installing staticcheck...
	@go install honnef.co/go/tools/cmd/staticcheck@latest
endif
ifndef GL_INSTALLED
	@echo Installing golint...
	@go install golang.org/x/lint/golint@latest
endif

# ==============================================================================
# Modules support

tidy:
	@echo Running go mod tidy...
	@go mod tidy

# ==============================================================================
# Build commands

build:install-tools
	@echo Building...
	@go build -v ./...

# ==============================================================================
# Test commands

lint: build
	@echo Running lints...
	@go vet ./...
	@staticcheck ./...
	@golint ./...
	@golangci-lint run

test:
	@echo Running tests...
	@go test -v -race -vet=off $$(go list ./... | grep -v /gen_pb/ | grep -v /googleapis/ | grep -v /proto/)

win: tidy
	@echo Building for windows...
	@GOOS=windows GOARCH=386 go build -o $(PROJECT_NAME).exe ./

mac: tidy
	@echo Building for mac...
	@GOOS=darwin GOARCH=amd64 go build -o $(PROJECT_NAME) ./

linux: tidy
	@echo Building for linux...
	@GOOS=linux GOARCH=amd64 go build -o $(PROJECT_NAME) ./
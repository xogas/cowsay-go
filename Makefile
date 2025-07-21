.PHONY: tidy build test

AppVersion := $(shell git describe --always 2>/dev/null || echo "--")
GitCommit := $(shell git rev-parse HEAD 2>/dev/null || echo "--")
BuildTime := $(shell date +%Y-%m-%dT%H:%M:%S%z)

LDFLAGS="-s -w -X github.com/xogas/cowsay-go/internal/version.AppVersion=$(AppVersion) \
	-X github.com/xogas/cowsay-go/internal/version.GitCommit=$(GitCommit) \
	-X github.com/xogas/cowsay-go/internal/version.BuildTime=$(BuildTime)"

# go mod tidy
tidy:
	go mod tidy

# build executable binary
build: tidy
	CGO_ENABLE=0 go build -ldflags ${LDFLAGS} -o cowsay-go ./main.go

# run unit test
test: tidy
	go test -v ./... -race -coverprofile=coverage.out

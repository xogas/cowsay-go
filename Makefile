.PHONY: all tidy build test fmt lint install-tools clean run


# Build metadata (evaluated at make time)
AppVersion := $(shell git describe --always 2>/dev/null || echo "--")
GitCommit := $(shell git rev-parse HEAD 2>/dev/null || echo "--")
BuildTime := $(shell date +%Y-%m-%dT%H:%M:%S%z)

# Local bins and tool versions
LocalBin := $(shell pwd)/bin
GolangciLintVersion := latest

# LDFLAGS for embedding build info
LDFLAGS := -s -w \
	-X github.com/xogas/cowsay-go/buildinfo.AppVersion=$(AppVersion) \
	-X github.com/xogas/cowsay-go/buildinfo.GitCommit=$(GitCommit) \
	-X github.com/xogas/cowsay-go/buildinfo.BuildTime=$(BuildTime)

# Default:
all: tidy lint vet test build

# Ensure local bin exists (for tools / installation)
$(LocalBin):
	mkdir -p $(LocalBin)

# install lint tool
.PHONY: install-tools
install-tools: $(LocalBin)
	GOBIN=$(LocalBin) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GolangciLintVersion)

# fmt: format code
fmt:
	go fmt ./...

# tidy: go mod tidy
tidy:
	go mod tidy

# vet
vet:
	go vet ./...

# lint
lint: install-tools
	$(LocalBin)/golangci-lint run

# build: produce a static binary
build: tidy
	CGO_ENABLE=0 go build -ldflags "${LDFLAGS}" -o cowsay-go ./main.go

# run unit test
test: tidy
	go test -v ./... -race -coverprofile=coverage.out

# clean: remove LocalBin and coverage
clean:
	rm -rf $(LocalBin) coverage.out cowsay-go

# Makefile for viewkit
# Usage:
#   make build              # build ./build/viewkit (CGO disabled)
#   make install            # install to $(GOBIN) or GOPATH/bin
#   make run ARGS="--help"  # run directly
#   make test               # run tests
#   make cross              # build for darwin/linux/windows amd64+arm64 into ./build
#   make clean              # remove build artifacts
#
# NOTE: If your main package lives at ./ (root) instead of ./cmd/viewkit,
# set MAIN_PKG=./ below.

BINARY      ?= viewkit
BUILD_DIR   ?= build
MAIN_PKG    ?= ./cmd/viewkit

# Metadata injected at build-time. Adjust variable names in -X if your main package uses different ones.
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT      ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE        ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS     ?= -X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.date=$(DATE)'
GOFLAGS     ?=
CGO_ENABLED ?= 0

.PHONY: all build install run test tidy clean cross env

all: build

build:
	@echo ">> building $(BINARY) → $(BUILD_DIR)/$(BINARY)"
	@mkdir -p $(BUILD_DIR)
	@GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY) $(MAIN_PKG)
	@echo ">> done"

install:
	@echo ">> installing $(BINARY) to GOBIN/GOPATH/bin"
	@GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) go install $(GOFLAGS) -ldflags "$(LDFLAGS)" $(MAIN_PKG)
	@echo ">> installed"

run:
	@echo ">> running $(BINARY)"
	@GO111MODULE=on go run $(GOFLAGS) -ldflags "$(LDFLAGS)" $(MAIN_PKG) $(ARGS)

test:
	@echo ">> testing"
	@GO111MODULE=on go test -count=1 ./...

tidy:
	@echo ">> go mod tidy"
	@GO111MODULE=on go mod tidy

clean:
	@echo ">> cleaning"
	@rm -rf $(BUILD_DIR)

cross:
	@echo ">> cross-compiling"
	@mkdir -p $(BUILD_DIR)
	@set -e; \
	for os in darwin linux windows; do \
	  for arch in amd64 arm64; do \
	    out="$(BUILD_DIR)/$(BINARY)-$$os-$$arch"; \
	    ext=""; if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
	    echo "   GOOS=$$os GOARCH=$$arch -> $$out$$ext"; \
	    GOOS=$$os GOARCH=$$arch CGO_ENABLED=$(CGO_ENABLED) go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o "$$out$$ext" $(MAIN_PKG); \
	  done; \
	done
	@echo ">> cross builds complete"

# Print a one-liner you can add to your shell profile to include ./build on PATH
env:
	@echo 'export PATH="$$PATH:$(PWD)/$(BUILD_DIR)"'

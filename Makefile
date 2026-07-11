.PHONY: all build build-desktop build-mobile clean test lint fmt help

BINARY    := omniconfig
VERSION   ?= 0.1.0-dev
BUILD_DIR ?= build

all: fmt lint test build-desktop

## build: Build for the current platform
build:
	@echo "==> Building for current platform..."
	go build -ldflags="-X 'github.com/omnicofig/cli/cmd.version=$(VERSION)'" -o $(BUILD_DIR)/$(BINARY) .

## build-desktop: Cross-compile for all desktop platforms (Linux/macOS/Windows)
build-desktop:
	@echo "==> Building for desktop platforms..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
		go build -ldflags="-X 'github.com/omnicofig/cli/cmd.version=$(VERSION)'" -o $(BUILD_DIR)/$(BINARY)-linux-amd64 .
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
		go build -ldflags="-X 'github.com/omnicofig/cli/cmd.version=$(VERSION)'" -o $(BUILD_DIR)/$(BINARY)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 \
		go build -ldflags="-X 'github.com/omnicofig/cli/cmd.version=$(VERSION)'" -o $(BUILD_DIR)/$(BINARY)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 \
		go build -ldflags="-X 'github.com/omnicofig/cli/cmd.version=$(VERSION)'" -o $(BUILD_DIR)/$(BINARY)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
		go build -ldflags="-X 'github.com/omnicofig/cli/cmd.version=$(VERSION)'" -o $(BUILD_DIR)/$(BINARY)-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 \
		go build -ldflags="-X 'github.com/omnicofig/cli/cmd.version=$(VERSION)'" -o $(BUILD_DIR)/$(BINARY)-windows-arm64.exe .
	@echo "Desktop builds complete."

## build-mobile: Cross-compile for mobile platforms (iOS/Android)
build-mobile:
	@echo "==> Building for mobile platforms..."
	GOOS=ios GOARCH=arm64 CGO_ENABLED=0 \
		go build -tags ios -ldflags="-X 'github.com/omnicofig/cli/cmd.version=$(VERSION)'" \
		-o $(BUILD_DIR)/$(BINARY)-ios-arm64 . 2>/dev/null || \
		echo "  ⚠ iOS build skipped (requires Xcode)"
	GOOS=android GOARCH=arm64 CGO_ENABLED=0 \
		go build -ldflags="-X 'github.com/omnicofig/cli/cmd.version=$(VERSION)'" \
		-o $(BUILD_DIR)/$(BINARY)-android-arm64 . 2>/dev/null || \
		echo "  ⚠ Android arm64 build skipped (requires NDK)"
	GOOS=android GOARCH=amd64 CGO_ENABLED=0 \
		go build -ldflags="-X 'github.com/omnicofig/cli/cmd.version=$(VERSION)'" \
		-o $(BUILD_DIR)/$(BINARY)-android-amd64 . 2>/dev/null || \
		echo "  ⚠ Android amd64 build skipped (requires NDK)"
	@echo "Mobile builds complete."

## build-all: Build everything (desktop + mobile)
build-all: build-desktop build-mobile

## clean: Remove build artifacts
clean:
	rm -rf $(BUILD_DIR)/
	@echo "Cleaned."

## test: Run all tests
test:
	@echo "==> Running tests..."
	go test -v -race -count=1 ./...

## lint: Run linters
lint:
	@echo "==> Running linters..."
	go vet ./...

## fmt: Format Go code
fmt:
	@echo "==> Formatting code..."
	go fmt ./...

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':'
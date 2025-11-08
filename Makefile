.PHONY: build install run clean test

# Get version from git tags
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -s -w -X github.com/rshdhere/vibecheck/cmd.version=$(VERSION)

# Build the binary
build:
	@echo "Building vibecheck $(VERSION)..."
	@go build -ldflags "$(LDFLAGS)" -o vibecheck .

# Install the binary to $GOPATH/bin
install:
	@echo "Installing vibecheck $(VERSION)..."
	@go install -ldflags "$(LDFLAGS)"

# Run without building
run:
	@go run -ldflags "$(LDFLAGS)" . $(ARGS)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f vibecheck
	@rm -rf dist/

# Run tests
test:
	@go test -v ./...

# Show current version
version:
	@echo $(VERSION)


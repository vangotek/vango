# VanGo Static Site Generator Makefile

.PHONY: build run serve clean test help install dev

# Default target
all: build

# Build the site
build: tailwind
	@echo "Building site..."
	@go run main.go

# Run Tailwind CSS compiler
tailwind:
	@echo "Running Tailwind CSS..."
	@./tailwindcss.exe -i assets/css/style.css -o static/style.css

# Run development server
serve:
	@echo "Starting development server..."
	@go run main.go -mode serve

# Run development server on custom port
serve-port:
	@echo "Starting development server on port 8080..."
	@go run main.go -mode serve -port 8080

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	@rm -rf public/
	@echo "Clean complete."

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Build binary
binary:
	@echo "Building binary..."
	@go build -o vango main.go
	@echo "Binary built: ./vango"

# Build binary for multiple platforms
build-all:
	@echo "Building for all platforms..."
	@GOOS=linux GOARCH=amd64 go build -o dist/vango-linux-amd64 main.go
	@GOOS=darwin GOARCH=amd64 go build -o dist/vango-darwin-amd64 main.go
	@GOOS=windows GOARCH=amd64 go build -o dist/vango-windows-amd64.exe main.go
	@echo "Binaries built in dist/"

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	@go vet ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@golangci-lint run

# Development workflow
dev: clean deps build
	@echo "Development setup complete!"

# Quick development server
quick:
	@go run main.go -mode serve

# Production build
prod: clean
	@echo "Building for production..."
	@go run main.go
	@echo "Production build complete in public/"

# Watch and rebuild (requires entr or similar)
watch:
	@echo "Watching for changes... (requires 'entr' to be installed)"
	@find . -name '*.go' -o -name '*.md' -o -name '*.html' -o -name '*.toml' | entr -r make build

# Install entr for file watching (macOS)
install-entr-mac:
	@brew install entr

# Install entr for file watching (Ubuntu/Debian)
install-entr-linux:
	@sudo apt-get install entr

# Show help
help:
	@echo "VanGo Static Site Generator - Available commands:"
	@echo ""
	@echo "  build         Build the static site"
	@echo "  serve         Start development server (port 1313)"
	@echo "  serve-port    Start development server on port 8080"
	@echo "  clean         Remove generated files"
	@echo "  deps          Install Go dependencies"
	@echo "  binary        Build VanGo binary"
	@echo "  build-all     Build binaries for all platforms"
	@echo "  test          Run tests"
	@echo "  fmt           Format Go code"
	@echo "  vet           Vet Go code"
	@echo "  lint          Lint code (requires golangci-lint)"
	@echo "  dev           Full development setup"
	@echo "  quick         Quick development server start"
	@echo "  prod          Production build"
	@echo "  watch         Watch files and rebuild (requires entr)"
	@echo "  help          Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build              # Build the site"
	@echo "  make serve              # Start dev server"
	@echo "  make clean build        # Clean and build"
	@echo "  make dev                # Full development setup"

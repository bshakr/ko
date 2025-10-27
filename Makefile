.PHONY: build install clean test help release

# Binary name
BINARY_NAME=ko
INSTALL_PATH=/usr/local/bin

# Version can be overridden at build time: make build VERSION=v1.0.0
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X github.com/bshakr/ko/cmd.Version=$(VERSION)"

# Build the application
build:
	go build $(LDFLAGS) -o $(BINARY_NAME)

# Install to system PATH
install: build
	cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✅ Installed $(BINARY_NAME) to $(INSTALL_PATH)"

# Uninstall from system PATH
uninstall:
	rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✅ Uninstalled $(BINARY_NAME)"

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	go clean

# Run tests
test:
	go test ./...

# Tidy dependencies
tidy:
	go mod tidy

# Release a new version
# Usage: make release VERSION=0.2.0
# This will:
#   1. Validate the version format
#   2. Ensure we're on main branch with clean working directory
#   3. Run tests
#   4. Create and push a git tag
#   5. Create a GitHub release
release:
	@# Check if VERSION was explicitly provided (not from command line if it matches default)
	@if [ "$(origin VERSION)" != "command line" ]; then \
		echo "❌ Error: VERSION is required. Usage: make release VERSION=0.2.0"; \
		echo "   Do not include 'v' prefix - it will be added automatically"; \
		exit 1; \
	fi
	@echo "🚀 Starting release process for v$(VERSION)..."
	@echo ""
	@# Check if gh CLI is installed
	@if ! command -v gh >/dev/null 2>&1; then \
		echo "❌ Error: GitHub CLI (gh) is required but not installed."; \
		echo "   Install it with: brew install gh"; \
		echo "   Then authenticate: gh auth login"; \
		exit 1; \
	fi
	@# Check if gh is authenticated
	@if ! gh auth status >/dev/null 2>&1; then \
		echo "❌ Error: GitHub CLI is not authenticated."; \
		echo "   Run: gh auth login"; \
		exit 1; \
	fi
	@# Ensure we're on main branch
	@current_branch=$$(git branch --show-current); \
	if [ "$$current_branch" != "main" ]; then \
		echo "❌ Error: Not on main branch (currently on $$current_branch)"; \
		echo "   Run: git checkout main"; \
		exit 1; \
	fi
	@echo "✓ On main branch"
	@# Check for clean working directory
	@if ! git diff-index --quiet HEAD --; then \
		echo "❌ Error: Working directory has uncommitted changes"; \
		echo "   Commit or stash changes first"; \
		exit 1; \
	fi
	@echo "✓ Working directory is clean"
	@# Check if tag already exists
	@if git rev-parse "v$(VERSION)" >/dev/null 2>&1; then \
		echo "❌ Error: Tag v$(VERSION) already exists"; \
		exit 1; \
	fi
	@echo "✓ Tag v$(VERSION) is available"
	@# Pull latest changes
	@echo ""
	@echo "📥 Pulling latest changes..."
	@git pull origin main
	@# Run tests
	@echo ""
	@echo "🧪 Running tests..."
	@$(MAKE) test
	@echo "✓ All tests passed"
	@# Build to verify
	@echo ""
	@echo "🔨 Building..."
	@$(MAKE) build VERSION=v$(VERSION)
	@echo "✓ Build successful"
	@# Create and push tag
	@echo ""
	@echo "🏷️  Creating git tag v$(VERSION)..."
	@git tag -a "v$(VERSION)" -m "Release v$(VERSION)"
	@git push origin "v$(VERSION)"
	@echo "✓ Tag created and pushed"
	@# Create GitHub release
	@echo ""
	@echo "📝 Creating GitHub release..."
	@echo "   Opening editor for release notes..."
	@gh release create "v$(VERSION)" \
		--title "v$(VERSION)" \
		--notes "" \
		--draft=false \
		--latest || { \
		echo "❌ Failed to create GitHub release"; \
		echo "   You can create it manually at: https://github.com/bshakr/ko/releases/new?tag=v$(VERSION)"; \
		exit 1; \
	}
	@echo ""
	@echo "✅ Release v$(VERSION) completed successfully!"
	@echo ""
	@echo "📦 Next steps:"
	@echo "   • GitHub Actions will automatically update the Homebrew formula"
	@echo "   • Users can install/update with: brew upgrade ko"
	@echo "   • View release at: https://github.com/bshakr/ko/releases/tag/v$(VERSION)"

# Help
help:
	@echo "Available targets:"
	@echo "  make build              - Build the application"
	@echo "  make install            - Install to $(INSTALL_PATH)"
	@echo "  make uninstall          - Remove from $(INSTALL_PATH)"
	@echo "  make clean              - Remove build artifacts"
	@echo "  make test               - Run tests"
	@echo "  make tidy               - Tidy Go modules"
	@echo "  make release VERSION=X  - Create and publish a new release"

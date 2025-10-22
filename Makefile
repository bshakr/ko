.PHONY: build install clean test help

# Binary name
BINARY_NAME=ko
INSTALL_PATH=/usr/local/bin

# Build the application
build:
	go build -o $(BINARY_NAME)

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

# Help
help:
	@echo "Available targets:"
	@echo "  make build     - Build the application"
	@echo "  make install   - Install to $(INSTALL_PATH)"
	@echo "  make uninstall - Remove from $(INSTALL_PATH)"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make test      - Run tests"
	@echo "  make tidy      - Tidy Go modules"

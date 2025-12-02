## Variables
HOSTNAME            := local
NAMESPACE           := growthbook
TYPE                := growthbook
VERSION             := 0.1.0

GOOS                := $(shell go env GOOS)
GOARCH              := $(shell go env GOARCH)

BINARY              := terraform-provider-$(TYPE)_v$(VERSION)
PLUGIN_DIR          := ~/.terraform.d/plugins/$(HOSTNAME)/$(NAMESPACE)/$(TYPE)/$(VERSION)/$(GOOS)_$(GOARCH)

## Default target
default: build

## Build the provider binary locally
build:
	@echo "🚧 Building provider..."
	go build -o $(BINARY) ./cmd/main.go
	@echo "✔️  Build complete: $(BINARY)"

## Install the provider into the Terraform/OpenTofu plugin directory
install: build
	@echo "📁 Installing provider into Terraform/OpenTofu plugins directory..."
	mkdir -p $(PLUGIN_DIR)
	cp -f $(BINARY) $(PLUGIN_DIR)
	chmod +x $(PLUGIN_DIR)/$(BINARY)
	@echo "✔️  Installed to $(PLUGIN_DIR)"

## Remove build artifacts
clean:
	@echo "🧹 Cleaning generated files..."
	rm -f $(BINARY)
	@echo "✔️  Clean complete"

fclean:
	@echo "🧹 Cleaning generated files..."
	rm -f $(BINARY)
	rm -rf $(PLUGIN_DIR)/$(BINARY)
	@echo "✔️  Clean complete"

info:
	@echo "HOSTNAME:   $(HOSTNAME)"
	@echo "NAMESPACE:  $(NAMESPACE)"
	@echo "TYPE:       $(TYPE)"
	@echo "VERSION:    $(VERSION)"
	@echo "OS:         $(GOOS)"
	@echo "ARCH:       $(GOARCH)"
	@echo "Binary:     $(BINARY)"
	@echo "Plugin dir: $(PLUGIN_DIR)"

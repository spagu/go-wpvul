APP_NAME=wpvul
BUILD_DIR=build

# Build matrix platforms and arches
PLATFORMS := linux darwin freebsd
ARCHITECTURES := amd64 arm64

.PHONY: all clean build brew $(PLATFORMS)

all: clean build

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

build:
	@echo "Building default..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) main.go

linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64 main.go

darwin:
	@echo "Building for Mac OS (Darwin)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 main.go

freebsd:
	@echo "Building for FreeBSD..."
	@mkdir -p $(BUILD_DIR)
	GOOS=freebsd GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-freebsd-amd64 main.go
	GOOS=freebsd GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-freebsd-arm64 main.go

compile-all: linux darwin freebsd
	@echo "Multi-arch build completed! Check the '$(BUILD_DIR)' directory."

brew: darwin
	@echo "Setting up Homebrew-like installation locally..."
	@mkdir -p /usr/local/bin
	@cp $(BUILD_DIR)/$(APP_NAME)-darwin-$$(uname -m | sed s/x86_64/amd64/ | sed s/aarch64/arm64/ | sed s/arm64/arm64/) /usr/local/bin/$(APP_NAME)
	@echo "$(APP_NAME) has been installed via make brew to /usr/local/bin/"

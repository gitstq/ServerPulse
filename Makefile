.PHONY: build install clean test lint run

# Binary name / 二进制名称
BINARY=serverpulse
# Build directory / 构建目录
BUILD_DIR=./bin
# Main package path / 主包路径
MAIN_PKG=./cmd/serverpulse

# Build the project / 构建项目
build:
	@echo "Building ServerPulse..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 go build -o $(BUILD_DIR)/$(BINARY) $(MAIN_PKG)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY)"

# Install to $GOPATH/bin / 安装到 GOPATH/bin
install:
	@echo "Installing ServerPulse..."
	CGO_ENABLED=1 go install $(MAIN_PKG)
	@echo "Install complete."

# Clean build artifacts / 清理构建产物
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY)
	@echo "Clean complete."

# Run tests / 运行测试
test:
	@echo "Running tests..."
	go test -v ./...

# Run linter / 运行代码检查
lint:
	@echo "Linting..."
	golangci-lint run ./...

# Run the binary / 运行二进制
run: build
	./$(BUILD_DIR)/$(BINARY) serve

# Build with optimizations / 优化构建
release:
	@echo "Building release..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY) $(MAIN_PKG)
	@echo "Release build complete: $(BUILD_DIR)/$(BINARY)"

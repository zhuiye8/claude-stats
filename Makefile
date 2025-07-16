# claude-stats Makefile
# 构建跨平台的Claude使用统计工具

BINARY_NAME=claude-stats
VERSION?=1.0.0
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go构建参数
GOFLAGS=-ldflags="-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# 平台定义
PLATFORMS=windows/amd64 linux/amd64 darwin/amd64 darwin/arm64

# 默认构建当前平台
.PHONY: build
build:
	@echo "🔨 构建 $(BINARY_NAME)..."
	go build $(GOFLAGS) -o $(BINARY_NAME) .

# 清理构建文件
.PHONY: clean
clean:
	@echo "🧹 清理构建文件..."
	rm -f $(BINARY_NAME)
	rm -rf dist/
	rm -f *.exe

# 安装依赖
.PHONY: deps
deps:
	@echo "📦 安装依赖..."
	go mod tidy
	go mod download

# 格式化代码
.PHONY: fmt
fmt:
	@echo "🎨 格式化代码..."
	go fmt ./...

# 代码检查
.PHONY: lint
lint:
	@echo "🔍 代码检查..."
	golangci-lint run --timeout=5m

# 运行测试
.PHONY: test
test:
	@echo "🧪 运行测试..."
	go test -v ./...

# 测试覆盖率
.PHONY: test-coverage
test-coverage:
	@echo "📊 测试覆盖率..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 构建所有平台的二进制文件
.PHONY: build-all
build-all: clean
	@echo "🚀 构建所有平台的二进制文件..."
	@mkdir -p dist
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		echo "构建 $$os/$$arch..."; \
		if [ "$$os" = "windows" ]; then \
			GOOS=$$os GOARCH=$$arch go build $(GOFLAGS) -o dist/$(BINARY_NAME)-$$os-$$arch.exe .; \
		else \
			GOOS=$$os GOARCH=$$arch go build $(GOFLAGS) -o dist/$(BINARY_NAME)-$$os-$$arch .; \
		fi; \
	done
	@echo "✅ 所有平台构建完成！文件位于 dist/ 目录"

# 创建发布包
.PHONY: release
release: build-all
	@echo "📦 创建发布包..."
	@cd dist && for file in $(BINARY_NAME)-*; do \
		echo "压缩 $$file..."; \
		if [[ "$$file" == *".exe" ]]; then \
			zip "$${file%.exe}.zip" "$$file" ../README.md ../LICENSE; \
		else \
			tar -czf "$$file.tar.gz" "$$file" ../README.md ../LICENSE; \
		fi; \
	done
	@echo "✅ 发布包创建完成！"

# 运行示例
.PHONY: demo
demo: build
	@echo "🎯 运行演示..."
	@echo "分析当前目录的示例数据..."
	./$(BINARY_NAME) analyze --verbose

# 安装到系统
.PHONY: install
install: build
	@echo "📥 安装到系统..."
	sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "✅ 安装完成！现在可以直接使用 '$(BINARY_NAME)' 命令"

# 卸载
.PHONY: uninstall
uninstall:
	@echo "🗑️  卸载..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✅ 卸载完成！"

# Docker构建
.PHONY: docker-build
docker-build:
	@echo "🐳 构建Docker镜像..."
	docker build -t claude-stats:$(VERSION) .
	docker tag claude-stats:$(VERSION) claude-stats:latest

# 开发环境设置
.PHONY: dev-setup
dev-setup:
	@echo "🛠️  设置开发环境..."
	@echo "安装开发工具..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "✅ 开发环境设置完成！"

# 快速检查（格式化+检查+测试）
.PHONY: check
check: fmt lint test
	@echo "✅ 所有检查通过！"

# 生成文档
.PHONY: docs
docs:
	@echo "📚 生成文档..."
	@mkdir -p docs
	./$(BINARY_NAME) --help > docs/help.txt
	./$(BINARY_NAME) analyze --help > docs/analyze-help.txt
	@echo "✅ 文档生成完成！"

# 显示帮助
.PHONY: help
help:
	@echo "claude-stats 构建工具"
	@echo ""
	@echo "可用命令："
	@echo "  build        - 构建当前平台的二进制文件"
	@echo "  build-all    - 构建所有平台的二进制文件"
	@echo "  release      - 创建发布包"
	@echo "  clean        - 清理构建文件"
	@echo "  deps         - 安装依赖"
	@echo "  fmt          - 格式化代码"
	@echo "  lint         - 代码检查"
	@echo "  test         - 运行测试"
	@echo "  test-coverage- 测试覆盖率"
	@echo "  demo         - 运行演示"
	@echo "  install      - 安装到系统"
	@echo "  uninstall    - 从系统卸载"
	@echo "  docker-build - 构建Docker镜像"
	@echo "  dev-setup    - 设置开发环境"
	@echo "  check        - 快速检查（格式化+检查+测试）"
	@echo "  docs         - 生成文档"
	@echo "  help         - 显示此帮助信息"

# 默认目标
.DEFAULT_GOAL := help 
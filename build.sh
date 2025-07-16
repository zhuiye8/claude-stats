#!/bin/bash
# claude-stats Unix构建脚本 (Linux/macOS/WSL)

set -e

BINARY_NAME="claude-stats"
VERSION="1.0.0"
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo "🔨 构建 $BINARY_NAME..."

# 构建当前平台
go build -ldflags="-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT" -o $BINARY_NAME .

echo "✅ 构建成功！"
echo "📍 二进制文件: $BINARY_NAME"
echo ""
echo "💡 使用方法:"
echo "   ./$BINARY_NAME analyze"
echo "   ./$BINARY_NAME analyze --help"
echo ""
echo "🔧 安装到系统 (可选):"
echo "   sudo cp $BINARY_NAME /usr/local/bin/"
echo "   # 或者"
echo "   ./install.sh" 
#!/bin/bash

# Claude Stats 本地构建脚本
# 用于在本地快速构建所有平台版本

set -e

echo "🚀 开始构建 Claude Stats..."

# 获取版本信息
VERSION=${1:-"v1.0.1"}
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo "📦 版本: $VERSION"
echo "⏰ 构建时间: $BUILD_TIME"
echo "🔗 Git提交: $GIT_COMMIT"

# 创建构建目录
BUILD_DIR="build"
mkdir -p "$BUILD_DIR"
cd "$BUILD_DIR"

echo ""
echo "🔨 开始构建各平台二进制文件..."

# 构建函数
build_platform() {
    local goos=$1
    local goarch=$2
    local extension=$3
    
    echo "  构建 $goos/$goarch..."
    
    if [ "$goos" = "windows" ]; then
        binary_name="claude-stats-${goos}-${goarch}.exe"
    else
        binary_name="claude-stats-${goos}-${goarch}"
    fi
    
    GOOS=$goos GOARCH=$goarch go build \
        -ldflags="-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT" \
        -o "$binary_name" \
        ../
    
    # 创建压缩包
    if [ "$goos" = "windows" ]; then
        zip "${binary_name%.exe}.zip" "$binary_name" ../README.md ../LICENSE
        echo "    ✅ 已创建: ${binary_name%.exe}.zip"
    else
        tar -czf "${binary_name}.tar.gz" "$binary_name" ../README.md ../LICENSE
        echo "    ✅ 已创建: ${binary_name}.tar.gz"
    fi
}

# 构建主要平台
build_platform "linux" "amd64"
build_platform "windows" "amd64"
build_platform "darwin" "amd64"

echo ""
echo "🎉 构建完成！"
echo ""
echo "📂 构建产物位于 build/ 目录："
ls -la

echo ""
echo "🧪 快速测试（Linux版本）："
echo "  ./claude-stats-linux-amd64 --version"
echo "  ./claude-stats-linux-amd64 analyze --help"

echo ""
echo "💡 使用说明："
echo "  1. 解压对应平台的压缩包"
echo "  2. 运行对应的二进制文件"
echo "  3. 享受强大的Claude使用统计功能！" 
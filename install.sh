#!/bin/bash
# claude-stats 安装脚本

set -e

BINARY_NAME="claude-stats"
INSTALL_DIR="/usr/local/bin"

echo "📦 安装 $BINARY_NAME 到系统..."

# 检查是否已构建
if [ ! -f "$BINARY_NAME" ]; then
    echo "⚠️  未找到 $BINARY_NAME 二进制文件"
    echo "请先运行构建脚本:"
    echo "  ./build.sh  (或 make build)"
    exit 1
fi

# 检查权限
if [ ! -w "$INSTALL_DIR" ]; then
    echo "📋 需要管理员权限安装到 $INSTALL_DIR"
    echo "正在使用 sudo..."
    sudo cp "$BINARY_NAME" "$INSTALL_DIR/"
    sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
else
    cp "$BINARY_NAME" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
fi

echo "✅ 安装完成！"
echo ""
echo "🎉 现在您可以在任何地方使用："
echo "   $BINARY_NAME analyze"
echo "   $BINARY_NAME analyze --help"
echo ""
echo "📍 安装位置: $INSTALL_DIR/$BINARY_NAME"
echo ""
echo "🗑️  卸载方法:"
echo "   sudo rm $INSTALL_DIR/$BINARY_NAME" 
#!/bin/bash

# Claude Stats 全局安装脚本
# 类似于 npx 的便捷安装方式

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 检测系统架构
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    case "$arch" in
        x86_64)
            arch="amd64"
            ;;
        aarch64|arm64)
            arch="arm64"
            ;;
        armv7l)
            arch="arm"
            ;;
        *)
            echo -e "${RED}❌ 不支持的架构: $arch${NC}"
            exit 1
            ;;
    esac
    
    case "$os" in
        linux)
            echo "linux-$arch"
            ;;
        darwin)
            echo "darwin-$arch"
            ;;
        mingw*|msys*|cygwin*)
            echo "windows-$arch"
            ;;
        *)
            echo -e "${RED}❌ 不支持的操作系统: $os${NC}"
            exit 1
            ;;
    esac
}

# 检查是否已安装Go
check_go() {
    if command -v go >/dev/null 2>&1; then
        local go_version=$(go version | cut -d' ' -f3 | cut -d'o' -f2)
        echo -e "${GREEN}✅ 检测到Go版本: $go_version${NC}"
        return 0
    else
        echo -e "${YELLOW}⚠️  未检测到Go环境${NC}"
        return 1
    fi
}

# Go安装方式（推荐）
install_with_go() {
    echo -e "${CYAN}🚀 使用Go直接安装...${NC}"
    
    # 检查Go环境
    if ! check_go; then
        echo -e "${RED}❌ 需要Go 1.21+环境${NC}"
        echo -e "${BLUE}📖 安装Go: https://golang.org/dl/${NC}"
        exit 1
    fi
    
    echo -e "${BLUE}📦 正在安装claude-stats...${NC}"
    
    # 使用go install直接安装
    if go install github.com/zhuiye8/claude-stats@latest; then
        echo -e "${GREEN}✅ 安装成功!${NC}"
        echo ""
        echo -e "${PURPLE}📋 快速使用:${NC}"
        echo -e "  ${CYAN}claude-stats${NC}                    # 查看每日使用情况"
        echo -e "  ${CYAN}claude-stats --breakdown${NC}        # 显示详细模型分解"
        echo -e "  ${CYAN}claude-stats blocks --live${NC}      # 实时监控5小时窗口"
        echo -e "  ${CYAN}claude-stats monthly${NC}            # 查看月度统计"
        echo ""
        echo -e "${YELLOW}💡 提示: 确保 \$GOPATH/bin 或 \$GOBIN 在您的 PATH 中${NC}"
        
        # 检查命令是否可用
        if command -v claude-stats >/dev/null 2>&1; then
            echo -e "${GREEN}🎉 claude-stats 命令已可用!${NC}"
            claude-stats --version
        else
            echo -e "${YELLOW}⚠️  请将Go bin目录添加到PATH:${NC}"
            echo -e "  ${CYAN}export PATH=\$PATH:\$(go env GOPATH)/bin${NC}"
        fi
    else
        echo -e "${RED}❌ 安装失败${NC}"
        exit 1
    fi
}

# 二进制安装方式
install_with_binary() {
    echo -e "${CYAN}🚀 使用预编译二进制安装...${NC}"
    
    local platform=$(detect_platform)
    local version="v2.0.0"
    local binary_name="claude-stats"
    
    if [[ "$platform" == *"windows"* ]]; then
        binary_name="claude-stats.exe"
    fi
    
    local download_url="https://github.com/zhuiye8/claude-stats/releases/download/${version}/claude-stats-${platform}"
    local install_dir="/usr/local/bin"
    
    # 检查权限
    if [[ ! -w "$install_dir" ]]; then
        echo -e "${YELLOW}⚠️  需要sudo权限安装到 $install_dir${NC}"
        install_dir="$HOME/.local/bin"
        mkdir -p "$install_dir"
        echo -e "${BLUE}📁 安装到用户目录: $install_dir${NC}"
    fi
    
    echo -e "${BLUE}📥 下载: $download_url${NC}"
    
    # 下载二进制文件
    if command -v curl >/dev/null 2>&1; then
        curl -L "$download_url" -o "$install_dir/claude-stats"
    elif command -v wget >/dev/null 2>&1; then
        wget "$download_url" -O "$install_dir/claude-stats"
    else
        echo -e "${RED}❌ 需要curl或wget来下载文件${NC}"
        exit 1
    fi
    
    # 设置执行权限
    chmod +x "$install_dir/claude-stats"
    
    echo -e "${GREEN}✅ 安装成功到: $install_dir/claude-stats${NC}"
    
    # 检查PATH
    if [[ ":$PATH:" != *":$install_dir:"* ]]; then
        echo -e "${YELLOW}⚠️  请将 $install_dir 添加到PATH:${NC}"
        echo -e "  ${CYAN}export PATH=\$PATH:$install_dir${NC}"
    else
        echo -e "${GREEN}🎉 claude-stats 命令已可用!${NC}"
        "$install_dir/claude-stats" --version
    fi
}

# 显示使用帮助
show_usage() {
    echo -e "${PURPLE}Claude Stats 全局安装脚本${NC}"
    echo ""
    echo -e "${CYAN}用法:${NC}"
    echo -e "  ${GREEN}bash <(curl -fsSL https://raw.githubusercontent.com/zhuiye8/claude-stats/main/install.sh)${NC}"
    echo ""
    echo -e "${CYAN}安装方式:${NC}"
    echo -e "  ${BLUE}1. Go安装 (推荐)${NC} - 如果您有Go环境"
    echo -e "  ${BLUE}2. 二进制安装${NC}     - 下载预编译的二进制文件"
    echo ""
    echo -e "${CYAN}特点:${NC}"
    echo -e "  ${GREEN}• 类似npx的便捷体验${NC}"
    echo -e "  ${GREEN}• 自动检测系统架构${NC}"
    echo -e "  ${GREEN}• 支持多种安装方式${NC}"
    echo -e "  ${GREEN}• 完整的中文支持${NC}"
}

# 主函数
main() {
    echo -e "${PURPLE}╔══════════════════════════════════════════╗${NC}"
    echo -e "${PURPLE}║                                          ║${NC}"
    echo -e "${PURPLE}║    🎯 Claude Stats 全局安装工具           ║${NC}"
    echo -e "${PURPLE}║    专业的Claude Code使用统计分析工具      ║${NC}"
    echo -e "${PURPLE}║                                          ║${NC}"
    echo -e "${PURPLE}╚══════════════════════════════════════════╝${NC}"
    echo ""
    
    # 检测平台
    local platform=$(detect_platform)
    echo -e "${BLUE}🔍 检测到平台: $platform${NC}"
    
    # 选择安装方式
    if check_go; then
        echo -e "${GREEN}🎯 推荐使用Go安装 (最新版本)${NC}"
        read -p "使用Go安装? [Y/n] " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Nn]$ ]]; then
            install_with_binary
        else
            install_with_go
        fi
    else
        echo -e "${BLUE}📦 使用二进制安装${NC}"
        install_with_binary
    fi
    
    echo ""
    echo -e "${GREEN}🎉 安装完成! 开始分析您的Claude Code使用情况吧!${NC}"
    echo ""
    echo -e "${CYAN}📚 更多命令:${NC}"
    echo -e "  ${CYAN}claude-stats help${NC}               # 查看所有命令"
    echo -e "  ${CYAN}claude-stats daily --help${NC}       # 查看daily命令帮助"
    echo -e "  ${CYAN}claude-stats monthly --breakdown${NC}  # 月度详细分析"
}

# 脚本入口
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi 
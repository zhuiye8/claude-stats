#!/bin/bash
# Claude Stats 全局安装脚本 (Linux/macOS)
# 将 claude-stats 安装为全局命令，就像 Claude Code 一样使用

set -euo pipefail

# 配置
TOOL_NAME="claude-stats"
EXECUTABLE_NAME="claude-stats"
GITHUB_REPO="zhuiye8/claude-stats"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# 检测操作系统和架构
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    case "$arch" in
        x86_64) arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        *) 
            echo -e "${RED}❌ 不支持的架构: $arch${NC}"
            exit 1
            ;;
    esac
    
    case "$os" in
        linux) 
            BINARY_NAME="claude-stats-linux-$arch"
            ;;
        darwin) 
            BINARY_NAME="claude-stats-darwin-$arch"
            ;;
        *) 
            echo -e "${RED}❌ 不支持的操作系统: $os${NC}"
            exit 1
            ;;
    esac
}

# 显示帮助
show_help() {
    cat << EOF
${CYAN}🚀 Claude Stats 全局安装工具${NC}

用法:
  ./install-global.sh              # 安装最新版本
  ./install-global.sh --force      # 强制重新安装
  ./install-global.sh --uninstall  # 卸载工具
  ./install-global.sh --help       # 显示此帮助

安装后，您可以在任何位置使用:
  ${BOLD}claude-stats analyze${NC}
  ${BOLD}claude-stats --version${NC}
  ${BOLD}claude-stats --help${NC}

EOF
}

# 检查依赖
check_dependencies() {
    local missing_deps=()
    
    if ! command -v curl >/dev/null 2>&1 && ! command -v wget >/dev/null 2>&1; then
        missing_deps+=("curl 或 wget")
    fi
    
    if ! command -v tar >/dev/null 2>&1; then
        missing_deps+=("tar")
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        echo -e "${RED}❌ 缺少依赖: ${missing_deps[*]}${NC}"
        echo -e "${YELLOW}💡 请安装缺少的依赖后重试${NC}"
        exit 1
    fi
}

# 获取安装路径
get_install_path() {
    # 检查常见的PATH目录
    local candidate_paths=(
        "$HOME/.local/bin"
        "$HOME/bin"
        "/usr/local/bin"
    )
    
    for path in "${candidate_paths[@]}"; do
        if [[ ":$PATH:" == *":$path:"* ]] || [ "$path" = "/usr/local/bin" ]; then
            # 如果路径不存在，尝试创建
            if [ ! -d "$path" ]; then
                if mkdir -p "$path" 2>/dev/null; then
                    echo -e "${GREEN}📂 创建安装目录: $path${NC}"
                fi
            fi
            
            # 检查是否可写
            if [ -w "$path" ] || [ -w "$(dirname "$path")" ]; then
                echo "$path"
                return 0
            fi
        fi
    done
    
    # 如果没有找到合适的路径，默认使用用户bin目录
    local default_path="$HOME/.local/bin"
    mkdir -p "$default_path"
    echo "$default_path"
}

# 添加到PATH（如果需要）
add_to_path() {
    local install_path="$1"
    
    # 检查PATH是否已包含安装目录
    if [[ ":$PATH:" == *":$install_path:"* ]]; then
        echo -e "${GREEN}✅ PATH已包含安装目录${NC}"
        return 0
    fi
    
    # 检测shell类型并添加到相应的配置文件
    local shell_config=""
    local shell_name=$(basename "$SHELL")
    
    case "$shell_name" in
        bash)
            if [ -f "$HOME/.bashrc" ]; then
                shell_config="$HOME/.bashrc"
            elif [ -f "$HOME/.bash_profile" ]; then
                shell_config="$HOME/.bash_profile"
            elif [ -f "$HOME/.profile" ]; then
                shell_config="$HOME/.profile"
            fi
            ;;
        zsh)
            shell_config="$HOME/.zshrc"
            ;;
        fish)
            shell_config="$HOME/.config/fish/config.fish"
            ;;
        *)
            shell_config="$HOME/.profile"
            ;;
    esac
    
    if [ -n "$shell_config" ]; then
        local path_line="export PATH=\"$install_path:\$PATH\""
        
        # 检查是否已经添加过
        if [ -f "$shell_config" ] && grep -q "$install_path" "$shell_config"; then
            echo -e "${GREEN}✅ PATH配置已存在${NC}"
        else
            echo "$path_line" >> "$shell_config"
            echo -e "${GREEN}✅ 已添加到 $shell_config${NC}"
            echo -e "${YELLOW}💡 请运行: source $shell_config 或重启终端${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️  无法自动添加到PATH，请手动添加: $install_path${NC}"
    fi
}

# 安装函数
install_claude_stats() {
    echo -e "${CYAN}🚀 开始安装 Claude Stats...${NC}"
    
    # 检测平台
    detect_platform
    
    # 检查是否已构建
    local built_executable="./build/$BINARY_NAME"
    if [ ! -f "$built_executable" ]; then
        echo -e "${RED}❌ 未找到构建的可执行文件: $built_executable${NC}"
        echo -e "${YELLOW}💡 请先运行: ./build-local.sh${NC}"
        exit 1
    fi
    
    # 获取安装路径
    local install_path
    install_path=$(get_install_path)
    local target_path="$install_path/$EXECUTABLE_NAME"
    
    # 检查现有安装
    if [ -f "$target_path" ] && [ "$FORCE" != "1" ]; then
        echo -e "${YELLOW}⚠️  Claude Stats 已安装在: $target_path${NC}"
        read -p "是否覆盖安装? (y/N): " choice
        case "$choice" in
            y|Y)
                echo -e "${CYAN}继续安装...${NC}"
                ;;
            *)
                echo -e "${RED}❌ 安装已取消${NC}"
                exit 1
                ;;
        esac
    fi
    
    # 复制可执行文件
    if cp "$built_executable" "$target_path"; then
        chmod +x "$target_path"
        echo -e "${GREEN}✅ 已安装到: $target_path${NC}"
    else
        echo -e "${RED}❌ 安装失败: 无法复制文件${NC}"
        exit 1
    fi
    
    # 添加到PATH
    add_to_path "$install_path"
    
    # 测试安装
    echo -e "\n${CYAN}🧪 测试安装...${NC}"
    
    # 临时添加到当前PATH
    export PATH="$install_path:$PATH"
    
    if command -v "$EXECUTABLE_NAME" >/dev/null 2>&1; then
        local version
        version=$("$EXECUTABLE_NAME" --version 2>/dev/null || echo "未知版本")
        echo -e "${GREEN}✅ 安装成功!${NC}"
        echo -e "${BLUE}📊 版本信息: $version${NC}"
        
        echo -e "\n${GREEN}🎉 安装完成! 现在您可以在任何位置使用:${NC}"
        echo -e "   ${BOLD}claude-stats analyze${NC}"
        echo -e "   ${BOLD}claude-stats --help${NC}"
        echo -e "   ${BOLD}claude-stats --version${NC}"
        
        if [[ ":$PATH:" != *":$install_path:"* ]]; then
            echo -e "\n${YELLOW}💡 注意: 请重启终端或运行 'source ~/.bashrc' (或相应的shell配置文件) 来刷新PATH${NC}"
        fi
    else
        echo -e "${RED}❌ 安装验证失败${NC}"
        echo -e "${YELLOW}💡 请检查文件权限或手动运行: $target_path --version${NC}"
    fi
}

# 卸载函数
uninstall_claude_stats() {
    echo -e "${CYAN}🗑️  开始卸载 Claude Stats...${NC}"
    
    local install_paths=(
        "$HOME/.local/bin/$EXECUTABLE_NAME"
        "$HOME/bin/$EXECUTABLE_NAME"
        "/usr/local/bin/$EXECUTABLE_NAME"
    )
    
    local found=false
    
    for target_path in "${install_paths[@]}"; do
        if [ -f "$target_path" ]; then
            if rm "$target_path"; then
                echo -e "${GREEN}✅ 已删除: $target_path${NC}"
                found=true
            else
                echo -e "${RED}❌ 删除失败: $target_path${NC}"
            fi
        fi
    done
    
    if [ "$found" = false ]; then
        echo -e "${YELLOW}⚠️  未找到已安装的 Claude Stats${NC}"
    else
        echo -e "${GREEN}✅ 卸载完成!${NC}"
        echo -e "${YELLOW}💡 PATH中的条目需要手动清理（如有需要）${NC}"
    fi
}

# 解析命令行参数
FORCE=0
UNINSTALL=0

while [[ $# -gt 0 ]]; do
    case $1 in
        --force)
            FORCE=1
            shift
            ;;
        --uninstall)
            UNINSTALL=1
            shift
            ;;
        --help|-h)
            show_help
            exit 0
            ;;
        *)
            echo -e "${RED}❌ 未知参数: $1${NC}"
            show_help
            exit 1
            ;;
    esac
done

# 主逻辑
if [ "$UNINSTALL" = "1" ]; then
    uninstall_claude_stats
    exit 0
fi

# 检查依赖
check_dependencies

# 执行安装
install_claude_stats 
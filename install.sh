#!/bin/bash
# Claude Stats 一键安装脚本 (Linux/macOS)
# 自动构建并安装为全局命令

set -euo pipefail

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# 参数解析
FORCE=0
HELP=0

while [[ $# -gt 0 ]]; do
    case $1 in
        --force)
            FORCE=1
            shift
            ;;
        --help|-h)
            HELP=1
            shift
            ;;
        *)
            echo -e "${RED}❌ 未知参数: $1${NC}"
            HELP=1
            shift
            ;;
    esac
done

show_help() {
    cat << EOF
${CYAN}🚀 Claude Stats 一键安装工具${NC}

用法:
  ./install.sh         # 构建并安装
  ./install.sh --force # 强制重新安装
  ./install.sh --help  # 显示此帮助

此脚本将：
  1. 🔨 构建当前平台版本
  2. 🌍 安装为全局命令
  3. ✅ 测试安装结果

安装后可在任何位置使用:
  ${BOLD}claude-stats analyze${NC}
  ${BOLD}claude-stats --version${NC}

EOF
}

if [ "$HELP" = "1" ]; then
    show_help
    exit 0
fi

echo -e "${GREEN}🚀 Claude Stats 一键安装开始...${NC}"
echo ""

# 步骤1: 构建
echo -e "${CYAN}🔨 步骤 1/3: 构建当前平台版本...${NC}"
if [ -f "./build-local.sh" ]; then
    chmod +x ./build-local.sh
    ./build-local.sh
else
    echo -e "${RED}❌ 未找到构建脚本 build-local.sh${NC}"
    exit 1
fi

echo -e "${GREEN}✅ 构建完成！${NC}"
echo ""

# 步骤2: 全局安装
echo -e "${CYAN}🌍 步骤 2/3: 安装为全局命令...${NC}"
install_args=()
if [ "$FORCE" = "1" ]; then
    install_args+=("--force")
fi

if [ -f "./install-global.sh" ]; then
    chmod +x ./install-global.sh
    ./install-global.sh "${install_args[@]}"
else
    echo -e "${RED}❌ 未找到安装脚本 install-global.sh${NC}"
    exit 1
fi

echo ""

# 步骤3: 测试
echo -e "${CYAN}🧪 步骤 3/3: 测试安装...${NC}"

# 刷新当前会话的PATH
if [ -d "$HOME/.local/bin" ]; then
    export PATH="$HOME/.local/bin:$PATH"
fi
if [ -d "$HOME/bin" ]; then
    export PATH="$HOME/bin:$PATH"
fi

if command -v claude-stats >/dev/null 2>&1; then
    version=$(claude-stats --version 2>/dev/null || echo "未知版本")
    echo -e "${GREEN}✅ 测试成功！${NC}"
    echo -e "${BLUE}📊 版本信息: $version${NC}"
else
    echo -e "${YELLOW}⚠️  命令测试失败，可能需要重启终端${NC}"
    echo -e "${YELLOW}💡 请尝试重启终端或运行 'source ~/.bashrc' 后测试${NC}"
fi

echo ""
echo -e "${GREEN}🎉 一键安装完成！${NC}"
echo ""
echo -e "${NC}现在您可以在任何位置使用：${NC}"
echo -e "  ${YELLOW}claude-stats analyze${NC}              # 分析Claude使用情况"
echo -e "  ${YELLOW}claude-stats analyze --verbose${NC}    # 详细分析模式"
echo -e "  ${YELLOW}claude-stats analyze --details${NC}    # 显示详细统计"
echo -e "  ${YELLOW}claude-stats --help${NC}               # 查看帮助"
echo -e "  ${YELLOW}claude-stats --version${NC}            # 查看版本"
echo ""
echo -e "${CYAN}💡 如果命令不可用，请重启终端或运行 'source ~/.bashrc'${NC}" 
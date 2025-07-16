# 🚀 Claude Stats

**完美的Claude Code使用统计工具** - 智能分析token使用、成本统计、订阅模式优化

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)

> 专为Claude Code用户设计的终极统计工具，解决现有工具的所有痛点

**⚡ [5分钟快速开始](QUICKSTART.md) | 📖 [完整文档](#使用指南) | 🐛 [问题反馈](https://github.com/zhuiye8/claude-stats/issues)**

## ✨ 核心优势

### 🎯 解决现有工具痛点
| 功能特性 | Claude Stats | ccusage | claude-code-log | claude-token-monitor |
|---------|-------------|---------|----------------|---------------------|
| 智能模式检测 | ✅ 自动识别 | ❌ 仅API | ❌ 仅基础 | ❌ 手动 |
| 订阅模式支持 | ✅ 等价成本分析 | ❌ 无 | ❌ 无 | ❌ 无 |
| 5小时窗口分析 | ✅ 内建支持 | ❌ 无 | ❌ 无 | ❌ 无 |
| 美化终端输出 | ✅ 渐变+图标 | ❌ 基础 | ❌ 基础 | ❌ 基础 |
| 成本可视化 | ✅ 进度条+色彩 | ❌ 无 | ❌ 无 | ❌ 无 |
| 跨平台支持 | ✅ 全平台 | ❌ 有限 | ❌ 有限 | ❌ 有限 |
| 缓存token分析 | ✅ 完整支持 | ❌ 无 | ❌ 无 | ❌ 无 |
| 订阅计划建议 | ✅ 智能推荐 | ❌ 无 | ❌ 无 | ❌ 无 |

### 🌟 独有功能
- **订阅模式"等价API成本"** - 当/cost命令失效时的完美替代
- **5小时重置窗口分析** - 专为Claude Code订阅模式设计
- **智能订阅计划建议** - 根据使用模式推荐最优计划
- **完整缓存token统计** - 精确的缓存创建和读取分析

## 🚀 快速开始

### 本地构建（推荐）

#### Windows用户
```powershell
# 1. 克隆项目
git clone https://github.com/zhuiye8/claude-stats.git
cd claude-stats

# 2. 构建
.\build-local.ps1

# 3. 运行
.\build\claude-stats-windows-amd64.exe analyze
```

#### Linux/macOS用户
```bash
# 1. 克隆项目
git clone https://github.com/zhuiye8/claude-stats.git
cd claude-stats

# 2. 构建
./build-local.sh

# 3. 运行 (Linux)
./build/claude-stats-linux-amd64 analyze

# 3. 运行 (macOS)
./build/claude-stats-darwin-amd64 analyze
```

#### 快速单平台构建
```bash
# 当前平台快速构建
go build -o claude-stats .

# 直接运行
./claude-stats analyze
```

### 系统要求
- **Go 1.21+** (用于构建)
- **Windows 10+** / **macOS 10.15+** / **Linux** (任意发行版)
- **WSL支持** (Windows Subsystem for Linux)

## 📖 使用指南

### 基础使用

```bash
# 自动分析默认Claude目录
./claude-stats analyze

# 分析指定目录
./claude-stats analyze ~/claude-logs

# 查看详细信息
./claude-stats analyze --details

# 导出JSON报告
./claude-stats analyze --format json --output report.json

# 按日期范围过滤
./claude-stats analyze --start 2025-07-01 --end 2025-07-16

# 按模型过滤
./claude-stats analyze --model sonnet
```

### 命令行选项

#### `analyze` - 分析使用情况
```bash
claude-stats analyze [目录] [选项]
```

**选项:**
- `--start, -s`: 开始日期 (YYYY-MM-DD)
- `--end, -e`: 结束日期 (YYYY-MM-DD)  
- `--model, -m`: 按模型过滤 (sonnet, haiku, opus, claude-4)
- `--format, -f`: 输出格式 (table, json, csv)
- `--output, -o`: 输出文件路径
- `--details, -d`: 显示详细信息
- `--no-color`: 禁用颜色输出
- `--verbose, -v`: 详细输出

**示例:**
```bash
# 基础分析
claude-stats analyze

# 高级过滤和美化输出
claude-stats analyze --start 2025-07-01 --model sonnet --details

# 禁用颜色输出 (适用于日志文件)
claude-stats analyze --no-color

# 导出报告
claude-stats analyze --format csv --output monthly-report.csv
```

## 🎨 支持的Claude模型和定价

| 模型 | 输入($/MTok) | 输出($/MTok) | 缓存($/MTok) |
|------|-------------|-------------|-------------|
| Claude 4 Sonnet | $15.00 | $75.00 | $1.875 |
| Claude 4 Opus | $60.00 | $300.00 | $7.50 |
| Claude 3.5 Sonnet | $3.00 | $15.00 | $0.375 |
| Claude 3.5 Haiku | $1.00 | $5.00 | $0.125 |

> 💡 **提示**: 定价会自动更新，支持向后兼容的模型名称识别

## 🔍 高级功能

### 5小时窗口分析
专为Claude Code订阅模式设计，分析每个5小时重置窗口的使用情况：

```bash
claude-stats analyze --details
```

输出包含:
- 每个窗口的请求数和token使用量
- 等价API成本
- 使用效率建议

### 智能模式检测
自动检测你的使用模式:
- **API模式**: 显示实际成本
- **订阅模式**: 显示等价成本和计划建议

### 缓存Token优化建议
分析缓存token使用效率，提供优化建议:
- 缓存命中率分析
- 成本节省计算
- 使用模式优化建议

## 🌍 跨平台支持

### Windows
- 支持 PowerShell 和 CMD
- 自动检测 `%APPDATA%\claude\projects`
- WSL环境自动适配

### macOS
- 支持 Terminal 和 iTerm2
- 自动检测 `~/Library/Application Support/claude/projects`
- Apple Silicon 原生支持

### Linux
- 支持所有主流发行版
- 自动检测 `~/.config/claude/projects`
- 单一二进制文件，无依赖

### WSL (Windows Subsystem for Linux)
- 完全兼容WSL 1和WSL 2
- 自动检测Windows和Linux路径
- 跨文件系统支持

## 🚀 性能优化

- **并发处理**: 多JSONL文件并行解析
- **内存优化**: 流式处理大文件，最低内存占用
- **缓存机制**: 重复分析时复用已解析数据
- **增量分析**: 只分析新增文件

## 🔧 配置

创建配置文件 `~/.claude-stats.yaml`:

```yaml
# 默认输出格式
default_format: table

# 默认显示详细信息
show_details: false

# 自定义Claude目录路径
claude_directory: "/path/to/your/claude/projects"

# 自定义模型定价 (覆盖默认定价)
custom_pricing:
  "claude-4-sonnet":
    input_price_per_mtoken: 15.0
    output_price_per_mtoken: 75.0
    cache_price_per_mtoken: 1.875

# 报告模板
report_template: |
  ## Claude使用报告
  **分析时间**: {{.AnalysisPeriod.Duration}}
  **总成本**: ${{.EstimatedCost.TotalCost}}
```

## 🔒 安全和隐私

- **本地处理**: 所有数据在本地处理，不上传到任何服务器
- **无网络请求**: 除版本检查外，无任何网络连接
- **数据保护**: 不存储敏感信息，只分析token使用统计

## 🛠️ 开发和构建

### 开发环境设置
```bash
git clone https://github.com/zhuiye8/claude-stats.git
cd claude-stats
go mod tidy
go run main.go analyze --help
```

### 构建所有平台
```bash
# 使用脚本构建
./build-local.sh v1.0.2           # Linux/macOS
.\build-local.ps1 -Version v1.0.2 # Windows

# 使用Makefile (如果有make)
make build-all

# 手动构建特定平台
GOOS=windows GOARCH=amd64 go build -o claude-stats-windows.exe .
GOOS=linux GOARCH=amd64 go build -o claude-stats-linux .
GOOS=darwin GOARCH=amd64 go build -o claude-stats-macos .
```

### 测试
```bash
# 运行测试
go test ./...

# 测试覆盖率
go test -cover ./...

# 功能测试
./claude-stats --version
./claude-stats analyze --help
```

## 🤝 贡献

我们欢迎各种形式的贡献！

### 贡献方式
1. **报告问题**: 在[Issues](https://github.com/zhuiye8/claude-stats/issues)中报告bug
2. **功能建议**: 提出新功能需求
3. **代码贡献**: 提交Pull Request
4. **文档改进**: 改善文档和示例

## 📝 更新日志

### v1.0.2 (2025-07-16)
- 🔧 简化为本地构建模式
- 📦 优化构建脚本和跨平台支持
- 📖 更新文档，专注本地使用
- 🧹 清理不必要的CI/CD配置

### v1.0.1 (2025-07-16)
- 🐛 修复GitHub Actions构建问题
- 📚 添加详细的故障排除指南

### v1.0.0 (2025-07-16)
- 🎉 首个正式版本发布
- ✅ 支持API和订阅模式智能检测
- ✅ 完整的token统计和成本分析
- ✅ 跨平台二进制文件支持
- ✅ 美化的终端输出和彩色支持
- ✅ 5小时窗口分析和订阅计划建议

## 📄 许可证

本项目采用 [MIT License](LICENSE) 开源协议。

## 💡 为什么选择Claude Stats？

- ✅ **专为Claude Code设计** - 深度理解使用模式
- ✅ **订阅模式优化** - 其他工具都不支持的功能
- ✅ **本地构建** - 完全掌控，无依赖外部服务
- ✅ **美化输出** - 专业级终端界面
- ✅ **跨平台** - 一次构建，到处运行
- ✅ **开源免费** - MIT协议，永久免费

**立即体验这个最强大的Claude使用统计工具！** 
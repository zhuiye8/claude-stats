# claude-stats - 完美的Claude Code使用统计工具

[![Go版本](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![许可证](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![平台支持](https://img.shields.io/badge/Platform-Windows%20|%20macOS%20|%20Linux%20|%20WSL-lightgrey.svg)](#安装)

> **专为Claude Code用户设计的终极使用统计工具** - 解决市面上现有工具都无法完美统计Claude Code使用情况的痛点！

## 🎯 核心优势

### 🔥 完美解决现有工具的痛点
- ✅ **智能双模式支持**: 自动识别API模式 vs 订阅模式，完美统计两种不同的计费方式
- ✅ **真正的Token统计**: 详细统计输入、输出、缓存创建、缓存读取token，不遗漏任何使用量
- ✅ **订阅模式专属功能**: 当`/cost`命令无效时，提供"等价API成本"让你了解真实使用价值
- ✅ **跨平台原生支持**: Windows、Mac、Linux、WSL一个二进制文件通吃
- ✅ **2025年最新定价**: 支持Claude 4、Claude 3.5全系列模型最新定价

### 📊 功能特性

| 功能 | claude-stats | ccusage | claude-code-log | claude-token-monitor |
|------|-------------|---------|-----------------|---------------------|
| API模式统计 | ✅ | ✅ | ❌ | ✅ |
| 订阅模式等价成本 | ✅ | ❌ | ❌ | ❌ |
| 5小时窗口分析 | ✅ | ❌ | ❌ | ❌ |
| 缓存Token统计 | ✅ | ✅ | ❌ | ❌ |
| 跨平台二进制 | ✅ | ❌ | ❌ | ❌ |
| 实时监控 | ✅ | ❌ | ❌ | ✅ |
| 多格式导出 | ✅ | ❌ | ✅ | ❌ |

## 🚀 快速开始

### 安装

#### 方式1: 下载预编译二进制文件 (推荐)
```bash
# Windows
curl -L https://github.com/claude-stats/claude-stats/releases/latest/download/claude-stats-windows.exe -o claude-stats.exe

# macOS
curl -L https://github.com/claude-stats/claude-stats/releases/latest/download/claude-stats-darwin -o claude-stats
chmod +x claude-stats

# Linux
curl -L https://github.com/claude-stats/claude-stats/releases/latest/download/claude-stats-linux -o claude-stats
chmod +x claude-stats

# WSL
curl -L https://github.com/claude-stats/claude-stats/releases/latest/download/claude-stats-linux -o claude-stats
chmod +x claude-stats
```

#### 方式2: 从源码编译
```bash
git clone https://github.com/claude-stats/claude-stats.git
cd claude-stats
go build -o claude-stats
```

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

## 📈 使用示例

### 基础统计输出
```
🎯 Claude Code 使用统计报告
============================================================

📊 基本信息:
   检测模式: 订阅模式 (按请求限制)
   总会话数: 8
   总消息数: 156
   分析时段: 2025-07-15 01:59 至 2025-07-16 01:15
   持续时间: 23h16m

📈 Token 使用统计
┌─────────────────┬──────────┬────────┐
│ 类型            │ 数量     │ 百分比  │
├─────────────────┼──────────┼────────┤
│ 输入Token       │ 21,543   │ 1.5%   │
│ 输出Token       │ 1,381    │ 0.1%   │
│ 缓存创建Token   │ 6,630    │ 0.5%   │
│ 缓存读取Token   │ 1,346,759│ 97.9%  │
├─────────────────┼──────────┼────────┤
│ 总计            │ 1,376,313│ 100.0% │
└─────────────────┴──────────┴────────┘

💰 成本分析:
   (基于订阅模式的API等价成本估算)
   输入成本:     $0.0646
   输出成本:     $0.0207
   缓存创建成本: $0.0199
   缓存读取成本: $0.5050
   总成本:       $0.6102

🎯 订阅计划建议:
   建议计划: Pro ($20)
   预估节省: $19.39/月
```

### JSON格式输出
```json
{
  "total_sessions": 8,
  "total_messages": 156,
  "total_tokens": {
    "input_tokens": 21543,
    "output_tokens": 1381,
    "cache_creation_input_tokens": 6630,
    "cache_read_input_tokens": 1346759,
    "total_tokens": 1376313
  },
  "model_stats": {
    "claude-sonnet-4-20250514": {
      "input_tokens": 13,
      "output_tokens": 1380,
      "cache_creation_input_tokens": 6630,
      "cache_read_input_tokens": 946759
    },
    "claude-3-5-haiku-20241022": {
      "input_tokens": 8,
      "output_tokens": 1
    }
  },
  "estimated_cost": {
    "input_cost": 0.0646,
    "output_cost": 0.0207,
    "cache_creation_cost": 0.0199,
    "cache_read_cost": 0.5050,
    "total_cost": 0.6102,
    "currency": "USD",
    "is_estimated": true
  },
  "detected_mode": "subscription"
}
```

## 🎛️ 命令详解

### `analyze` - 分析使用统计
```bash
claude-stats analyze [目录路径] [选项]
```

**选项:**
- `--format, -f`: 输出格式 (table, json, csv) [默认: table]
- `--output, -o`: 输出文件路径
- `--start`: 开始日期 (YYYY-MM-DD)
- `--end`: 结束日期 (YYYY-MM-DD)
- `--model`: 过滤特定模型
- `--details, -d`: 显示详细信息
- `--verbose, -v`: 详细输出

**示例:**
```bash
# 基础分析
claude-stats analyze

# 高级过滤
claude-stats analyze --start 2025-07-01 --model sonnet --details

# 导出报告
claude-stats analyze --format csv --output monthly-report.csv
```

### `monitor` - 实时监控 (规划中)
```bash
claude-stats monitor [选项]
```

### `report` - 生成报告 (规划中)
```bash
claude-stats report --monthly --output report.md
```

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

## 🔒 安全和隐私

- **本地处理**: 所有数据在本地处理，不上传到任何服务器
- **无网络请求**: 除版本检查外，无任何网络连接
- **数据保护**: 不存储敏感信息，只分析token使用统计

## 🤝 贡献

我们欢迎各种形式的贡献！

### 贡献方式
1. **报告问题**: 在[Issues](https://github.com/claude-stats/claude-stats/issues)中报告bug
2. **功能建议**: 提出新功能需求
3. **代码贡献**: 提交Pull Request
4. **文档改进**: 改善文档和示例

### 开发环境设置
```bash
git clone https://github.com/claude-stats/claude-stats.git
cd claude-stats
go mod tidy
go run main.go analyze --help
```

## 📝 更新日志

### v1.0.0 (2025-07-16)
- 🎉 首个正式版本发布
- ✅ 支持API和订阅模式智能检测
- ✅ 完整的token统计和成本分析
- ✅ 跨平台二进制文件支持
- ✅ 多种输出格式 (表格、JSON、CSV)
- ✅ 2025年最新Claude模型定价

## 📄 许可证

本项目基于 [MIT License](LICENSE) 开源。

## 🙏 致谢

- 感谢 [Anthropic](https://anthropic.com) 提供的优秀Claude模型
- 感谢开源社区提供的各种Go库支持
- 感谢所有测试用户的反馈和建议

---

**如果这个工具帮助到了您，请给个⭐️支持一下！**

> 💬 **需要帮助?** 
> - 查看 [文档](https://github.com/claude-stats/claude-stats/wiki)
> - 提交 [Issue](https://github.com/claude-stats/claude-stats/issues)
> - 加入 [讨论](https://github.com/claude-stats/claude-stats/discussions) 
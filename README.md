# Claude Stats - Claude Code 使用统计分析工具

[![Version](https://img.shields.io/badge/version-1.0.9-blue.svg)](https://github.com/zhuiye8/claude-stats)
[![Go](https://img.shields.io/badge/go-1.21%2B-brightgreen.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## 🎯 功能特色

Claude Stats 是一个专为 Claude Code 用户设计的使用统计分析工具，帮助您：

### 📊 详细的使用分析
- **Token 使用统计** - 输入、输出、缓存创建、缓存读取的详细分析
- **成本估算** - 基于官方API定价的成本计算
- **百分比分析** - 使用科学的"基础Token"方法计算百分比（避免缓存Token干扰）
- **模型使用统计** - 按模型类型分析使用情况

### 🔍 订阅模式支持
- **智能计划检测** - 自动识别Pro/Max5x/Max20x计划
- **限额状态监控** - 实时显示当前窗口使用情况
- **重置时间预测** - 考虑用户时区的重置时间计算
- **模型切换提示** - 显示当前使用的模型类型（Opus/Sonnet）

### ⚠️ 数据准确性说明
- **估算数据警告** - 明确标示估算值的局限性
- **官方命令建议** - 推荐使用Claude Code的 `/status` 命令获取准确信息
- **时区自动检测** - 显示用户当前时区和UTC转换
- **谨慎的用词** - 使用"推测"、"估算"等词汇避免误导

### 📈 项目和会话分析
- **项目统计** - 按项目分类的使用情况
- **每日使用趋势** - 按日期分析使用模式
- **会话详情** - 最近会话的详细信息

### 🎨 美观的界面
- **ASCII艺术标题** - 彩色渐变效果的品牌标识
- **作者标识** - 右下角显示"作者: zhuiye"
- **丰富的颜色** - 使用emoji和颜色增强可读性
- **表格格式** - 清晰的数据展示

## 🚀 快速开始

### 安装

1. **从源码构建**（推荐）
```bash
git clone https://github.com/zhuiye8/claude-stats.git
cd claude-stats
go build -o claude-stats main.go
```

2. **运行分析**
```bash
# 分析当前目录的Claude Code数据
./claude-stats analyze

# 显示详细信息
./claude-stats analyze --details

# 分析特定目录
./claude-stats analyze /path/to/claude/data --details
```

### 系统要求

- Go 1.21+ （仅构建时需要）
- Claude Code 的JSONL数据文件
- 支持Windows/macOS/Linux

## 📋 使用说明

### 基本命令

```bash
# 基础分析（简洁模式）
claude-stats analyze

# 详细分析（推荐）
claude-stats analyze --details

# 指定数据目录
claude-stats analyze /path/to/data --details

# 时间范围过滤
claude-stats analyze --start 2025-07-01 --end 2025-07-16

# 特定模型过滤
claude-stats analyze --model claude-sonnet-4

# 导出为JSON
claude-stats analyze --format json --output stats.json
```

### 数据位置

Claude Code数据通常位于：
- **Windows**: `%USERPROFILE%\.claude\history\`
- **macOS**: `~/.claude/history/`
- **Linux**: `~/.claude/history/`

## 🔧 配置说明

### 环境变量
- `CLAUDE_DATA_DIR` - 指定Claude数据目录
- `NO_COLOR` - 设置为任意值以禁用颜色输出

### 配置文件
支持YAML配置文件（可选）：
```yaml
# ~/.claude-stats.yaml
data_dir: "/path/to/claude/data"
default_format: "table"
show_details: true
```

## ⚠️ 重要提醒

### 数据准确性
本工具提供的订阅限额数据为**估算值**，包括：
- 当前使用量估算
- 剩余消息数推测  
- 重置时间预测
- 计划类型推测

**获取准确信息**：请在Claude Code中运行 `/status` 命令查看官方数据。

### 时区考虑
- 重置时间可能基于UTC时区
- 显示包含您的本地时区信息
- 如果重置时间不准确，请反馈给开发者

## 📊 输出示例

```
🎯 订阅限额状态
   ⚠️  注意：以下数据为估算值，实际限额请使用Claude Code的 /status 命令查看
   💡 在Claude Code中运行 /status 可获取准确的当前窗口使用情况

   🥇 推测计划: Max20x ($200/月)
   🕐 限额机制: 5小时窗口 (每天4个窗口，可能基于UTC时区)
   🟡 估算使用:  280 / 900 消息
   📊 估算进度: [██████░░░░░░░░░░░░░░] 31.1%
   ✨ 估算剩余: 620
   ⚡ 推测模型: Claude 4 Sonnet (标准模型)
   ⏳ 预测重置: 18:00 (4小时26分钟后) [UTC+8]

   💡 获取准确信息：在Claude Code中运行 /status 命令
   🔧 如果重置时间不准确，请反馈给开发者
```

## 🛠 开发者信息

### 技术实现
- **语言**: Go 1.21+
- **架构**: 模块化设计，支持扩展
- **数据格式**: JSONL解析，兼容Claude Code格式
- **时区处理**: 自动检测用户时区，UTC转换

### 项目结构
```
claude-stats/
├── cmd/           # 命令行接口
├── pkg/
│   ├── models/    # 数据模型
│   ├── parser/    # JSONL解析器
│   └── formatter/ # 输出格式化
├── main.go        # 程序入口
└── README.md
```

## 📝 更新日志

### v1.0.9 (2025-07-16)
- **🔧 修复时区处理** - 正确显示用户时区和UTC转换
- **⚠️ 添加数据准确性警告** - 明确标示估算值的局限性
- **💡 优化重置时间计算** - 基于整点重置机制改进算法
- **📝 改进用词** - 使用"推测"、"估算"等谨慎表述
- **🎯 增强使用建议** - 推荐使用官方`/status`命令

### v1.0.8 (2025-07-16)
- **🎨 添加ASCII艺术标题** - 彩色渐变效果
- **👤 添加作者标识** - 显示"作者: zhuiye"
- **🔢 修复百分比计算** - 使用基础Token方法
- **📊 增强订阅支持** - 智能计划检测和限额估算
- **🌈 美化界面** - emoji和颜色优化

## 🤝 贡献指南

欢迎提交Issue和Pull Request！

### 报告问题
- 时区显示不准确
- 重置时间计算错误
- 数据解析问题
- 功能建议

### 开发贡献
1. Fork项目
2. 创建特性分支
3. 提交更改
4. 发起Pull Request

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

## 🙏 致谢

- Claude Code团队提供的优秀开发工具
- Go社区的开源库支持
- 用户反馈和建议

---

**作者**: zhuiye  
**项目**: https://github.com/zhuiye8/claude-stats  
**问题反馈**: https://github.com/zhuiye8/claude-stats/issues

> 💡 **提示**: 本工具为非官方工具，仅供参考。准确的使用情况请以Claude Code官方显示为准。 
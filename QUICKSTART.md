# 🚀 Claude Stats 快速开始

## ⚡ 5分钟快速上手

### 1. 下载项目
```bash
git clone https://github.com/zhuiye8/claude-stats.git
cd claude-stats
```

### 2. 一键构建

**Windows:**
```powershell
.\build-local.ps1
```

**Linux/macOS:**
```bash
./build-local.sh
```

### 3. 运行分析
```bash
# Windows
.\build\claude-stats-windows-amd64.exe analyze

# Linux  
./build/claude-stats-linux-amd64 analyze

# macOS
./build/claude-stats-darwin-amd64 analyze
```

## 🎯 常用命令

```bash
# 基础分析
./claude-stats analyze

# 查看详细信息
./claude-stats analyze --details

# 导出JSON报告
./claude-stats analyze --format json --output report.json

# 按日期过滤
./claude-stats analyze --start 2025-07-01 --end 2025-07-16
```

## 🔧 系统要求

- **Go 1.21+** (用于构建)
- **Windows 10+** / **macOS 10.15+** / **Linux**

## 💡 提示

- 构建后的二进制文件在 `build/` 目录
- 支持所有主流Claude模型（Claude 4, 3.5 Sonnet, Haiku等）
- 自动检测订阅模式和API模式
- 完整的中文界面和帮助信息

**就是这么简单！** 🎉 
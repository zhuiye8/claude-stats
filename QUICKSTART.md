# 🚀 Claude Stats 快速开始指南

## 📋 平台兼容性

| 平台 | 构建方式 | 安装方式 | 命令行支持 |
|------|---------|---------|-----------|
| **Windows** | `build.bat` 或 直接go build | `install.ps1` | ✅ PowerShell/CMD |
| **WSL** | `./build.sh` 或 `make build` | `./install.sh` | ✅ 完全兼容 |
| **Linux** | `./build.sh` 或 `make build` | `./install.sh` | ✅ 原生支持 |
| **macOS** | `./build.sh` 或 `make build` | `./install.sh` | ✅ 原生支持 |

## 🔧 快速安装 (推荐方式)

### Windows (PowerShell)
```powershell
# 1. 构建
.\build.bat

# 2. 安装到系统 (可选)
.\install.ps1

# 3. 使用
claude-stats analyze
```

### WSL/Linux/macOS
```bash
# 1. 构建
./build.sh

# 2. 安装到系统 (可选)
./install.sh

# 3. 使用
claude-stats analyze
```

## 🛠️ 无Make环境解决方案

### Windows (不使用make)
```powershell
# 方法1: 使用批处理脚本
.\build.bat                    # 构建单平台
.\build-all.bat               # 构建所有平台

# 方法2: 直接go命令
go build -o claude-stats.exe .
```

### 其他平台 (不使用make)
```bash
# 方法1: 使用shell脚本
./build.sh                    # 构建单平台

# 方法2: 直接go命令
go build -o claude-stats .
```

## 📦 全局安装 (像Claude Code一样)

### 🎯 方法1: 使用安装脚本 (推荐)

**Windows:**
```powershell
.\install.ps1  # 自动添加到PATH
```

**Unix系统:**
```bash
./install.sh  # 安装到 /usr/local/bin
```

### 🎯 方法2: 手动安装

**Windows:**
```powershell
# 复制到用户bin目录
mkdir $env:USERPROFILE\bin -Force
copy claude-stats.exe $env:USERPROFILE\bin\

# 添加到PATH (一次性设置)
$oldPath = [Environment]::GetEnvironmentVariable("PATH", "User")
$newPath = "$oldPath;$env:USERPROFILE\bin"
[Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
```

**Unix系统:**
```bash
# 复制到系统bin目录
sudo cp claude-stats /usr/local/bin/
sudo chmod +x /usr/local/bin/claude-stats
```

## ⚡ 立即使用

安装完成后，在任何目录都可以直接使用：

```bash
# 基础分析
claude-stats analyze

# 分析指定目录
claude-stats analyze ~/claude-projects

# 查看帮助
claude-stats --help
claude-stats analyze --help

# 高级用法
claude-stats analyze --details --model sonnet
claude-stats analyze --start 2025-07-01 --format json
```

## 🔍 常见问题

### Q: Windows下没有make命令怎么办？
**A:** 使用提供的批处理脚本：
- `build.bat` - 构建
- `build-all.bat` - 全平台构建
- `install.ps1` - 安装

### Q: 如何卸载？
**A:** 
```bash
# Unix系统
sudo rm /usr/local/bin/claude-stats

# Windows (PowerShell)
Remove-Item "$env:USERPROFILE\bin\claude-stats.exe"
```

### Q: 安装后找不到命令？
**A:** 
1. 重启命令行/终端
2. 检查PATH设置
3. 使用完整路径测试

### Q: WSL中能否使用Windows的Claude日志？
**A:** 可以！WSL可以访问Windows文件系统：
```bash
claude-stats analyze /mnt/c/Users/YourName/AppData/Roaming/claude/projects
```

## 🎉 完成！

现在您就拥有了一个完美的Claude Code使用统计工具，可以在任何地方使用命令行直接分析您的Claude使用情况！ 
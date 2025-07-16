#!/usr/bin/env pwsh
# Claude Stats 一键安装脚本 (Windows)
# 自动构建并安装为全局命令

param(
    [switch]$Force,     # 强制重新安装
    [switch]$Help       # 显示帮助
)

function Show-Help {
    Write-Host @"
🚀 Claude Stats 一键安装工具

用法:
  .\install.ps1         # 构建并安装
  .\install.ps1 -Force  # 强制重新安装
  .\install.ps1 -Help   # 显示此帮助

此脚本将：
  1. 🔨 构建 Windows 版本
  2. 🌍 安装为全局命令
  3. ✅ 测试安装结果

安装后可在任何位置使用:
  claude-stats analyze
  claude-stats --version

"@ -ForegroundColor Cyan
}

if ($Help) {
    Show-Help
    exit 0
}

Write-Host "🚀 Claude Stats 一键安装开始..." -ForegroundColor Green
Write-Host ""

# 步骤1: 构建
Write-Host "🔨 步骤 1/3: 构建 Windows 版本..." -ForegroundColor Cyan
$buildArgs = @()
if (Test-Path ".\build-local.ps1") {
    & .\build-local.ps1
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ 构建失败！" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "❌ 未找到构建脚本 build-local.ps1" -ForegroundColor Red
    exit 1
}

Write-Host "✅ 构建完成！" -ForegroundColor Green
Write-Host ""

# 步骤2: 全局安装
Write-Host "🌍 步骤 2/3: 安装为全局命令..." -ForegroundColor Cyan
$installArgs = @()
if ($Force) {
    $installArgs += "-Force"
}

if (Test-Path ".\install-global.ps1") {
    & .\install-global.ps1 @installArgs
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ 安装失败！" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "❌ 未找到安装脚本 install-global.ps1" -ForegroundColor Red
    exit 1
}

Write-Host ""

# 步骤3: 测试
Write-Host "🧪 步骤 3/3: 测试安装..." -ForegroundColor Cyan

# 刷新PATH（仅当前会话）
$userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
$machinePath = [Environment]::GetEnvironmentVariable("PATH", "Machine")
$env:PATH = "$userPath;$machinePath"

try {
    $version = claude-stats --version 2>&1
    Write-Host "✅ 测试成功！" -ForegroundColor Green
    Write-Host "📊 版本信息: $version" -ForegroundColor Blue
} catch {
    Write-Host "⚠️  命令测试失败，可能需要重启终端" -ForegroundColor Yellow
    Write-Host "💡 请尝试重启终端后运行: claude-stats --version" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "🎉 一键安装完成！" -ForegroundColor Green
Write-Host ""
Write-Host "现在您可以在任何位置使用：" -ForegroundColor White
Write-Host "  claude-stats analyze              # 分析Claude使用情况" -ForegroundColor Yellow
Write-Host "  claude-stats analyze --verbose    # 详细分析模式" -ForegroundColor Yellow
Write-Host "  claude-stats analyze --details    # 显示详细统计" -ForegroundColor Yellow
Write-Host "  claude-stats --help               # 查看帮助" -ForegroundColor Yellow
Write-Host "  claude-stats --version            # 查看版本" -ForegroundColor Yellow
Write-Host ""
Write-Host "💡 如果命令不可用，请重启终端或运行 'refreshenv'" -ForegroundColor Cyan 
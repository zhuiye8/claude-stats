#!/usr/bin/env pwsh
# Claude Stats 全局安装脚本 (Windows)
# 将 claude-stats 安装为全局命令，就像 Claude Code 一样使用

param(
    [switch]$Force,        # 强制覆盖已存在的安装
    [switch]$Uninstall,    # 卸载命令
    [switch]$Help          # 显示帮助
)

# 配置
$TOOL_NAME = "claude-stats"
$EXECUTABLE_NAME = "claude-stats.exe"
$GITHUB_REPO = "zhuiye8/claude-stats"

# 检测安装路径
$INSTALL_PATHS = @(
    "$env:USERPROFILE\.local\bin",           # 用户级安装 (推荐)
    "$env:PROGRAMFILES\$TOOL_NAME",          # 系统级安装
    "$env:LOCALAPPDATA\Programs\$TOOL_NAME"  # 应用级安装
)

function Show-Help {
    Write-Host @"
🚀 Claude Stats 全局安装工具

用法:
  .\install-global.ps1              # 安装最新版本
  .\install-global.ps1 -Force       # 强制重新安装
  .\install-global.ps1 -Uninstall   # 卸载工具
  .\install-global.ps1 -Help        # 显示此帮助

安装后，您可以在任何位置使用:
  claude-stats analyze
  claude-stats --version
  claude-stats --help

"@ -ForegroundColor Cyan
}

function Test-AdminPrivileges {
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

function Get-InstallPath {
    # 优先选择用户级安装路径
    $userPath = $INSTALL_PATHS[0]
    
    # 如果用户路径不存在，创建它
    if (-not (Test-Path $userPath)) {
        New-Item -ItemType Directory -Path $userPath -Force | Out-Null
        Write-Host "📂 创建安装目录: $userPath" -ForegroundColor Green
    }
    
    return $userPath
}

function Add-ToPath {
    param($InstallPath)
    
    # 检查是否已在PATH中
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($currentPath -split ';' -contains $InstallPath) {
        Write-Host "✅ PATH已包含安装目录" -ForegroundColor Green
        return
    }
    
    # 添加到用户PATH
    try {
        $newPath = "$currentPath;$InstallPath"
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        Write-Host "✅ 已添加到用户PATH: $InstallPath" -ForegroundColor Green
        Write-Host "💡 请重启终端或运行: refreshenv" -ForegroundColor Yellow
    } catch {
        Write-Host "❌ 添加PATH失败: $($_.Exception.Message)" -ForegroundColor Red
        Write-Host "💡 请手动添加 '$InstallPath' 到系统PATH" -ForegroundColor Yellow
    }
}

function Install-ClaudeStats {
    Write-Host "🚀 开始安装 Claude Stats..." -ForegroundColor Cyan
    
    # 检查是否已构建
    $builtExecutable = ".\build\claude-stats-windows-amd64.exe"
    if (-not (Test-Path $builtExecutable)) {
        Write-Host "❌ 未找到构建的可执行文件: $builtExecutable" -ForegroundColor Red
        Write-Host "💡 请先运行: .\build-local.ps1" -ForegroundColor Yellow
        exit 1
    }
    
    # 获取安装路径
    $installPath = Get-InstallPath
    $targetPath = Join-Path $installPath $EXECUTABLE_NAME
    
    # 检查现有安装
    if ((Test-Path $targetPath) -and -not $Force) {
        Write-Host "⚠️  Claude Stats 已安装在: $targetPath" -ForegroundColor Yellow
        $choice = Read-Host "是否覆盖安装? (y/N)"
        if ($choice -ne 'y' -and $choice -ne 'Y') {
            Write-Host "❌ 安装已取消" -ForegroundColor Red
            exit 1
        }
    }
    
    # 复制可执行文件
    try {
        Copy-Item $builtExecutable $targetPath -Force
        Write-Host "✅ 已安装到: $targetPath" -ForegroundColor Green
    } catch {
        Write-Host "❌ 安装失败: $($_.Exception.Message)" -ForegroundColor Red
        exit 1
    }
    
    # 添加到PATH
    Add-ToPath $installPath
    
    # 测试安装
    Write-Host "`n🧪 测试安装..." -ForegroundColor Cyan
    
    # 刷新当前会话的PATH
    $env:PATH = "$env:PATH;$installPath"
    
    try {
        $version = & $targetPath --version 2>&1
        Write-Host "✅ 安装成功!" -ForegroundColor Green
        Write-Host "📊 版本信息: $version" -ForegroundColor Blue
        
        Write-Host "`n🎉 安装完成! 现在您可以在任何位置使用:" -ForegroundColor Green
        Write-Host "   claude-stats analyze" -ForegroundColor White
        Write-Host "   claude-stats --help" -ForegroundColor White
        Write-Host "   claude-stats --version" -ForegroundColor White
        
        if ($currentPath -notlike "*$InstallPath*") {
            Write-Host "`n💡 注意: 请重启终端或运行 'refreshenv' 来刷新PATH" -ForegroundColor Yellow
        }
        
    } catch {
        Write-Host "❌ 安装验证失败: $($_.Exception.Message)" -ForegroundColor Red
        Write-Host "💡 请检查文件权限或手动运行: $targetPath --version" -ForegroundColor Yellow
    }
}

function Uninstall-ClaudeStats {
    Write-Host "🗑️  开始卸载 Claude Stats..." -ForegroundColor Cyan
    
    $found = $false
    
    foreach ($installPath in $INSTALL_PATHS) {
        $targetPath = Join-Path $installPath $EXECUTABLE_NAME
        if (Test-Path $targetPath) {
            try {
                Remove-Item $targetPath -Force
                Write-Host "✅ 已删除: $targetPath" -ForegroundColor Green
                $found = $true
            } catch {
                Write-Host "❌ 删除失败: $targetPath - $($_.Exception.Message)" -ForegroundColor Red
            }
        }
    }
    
    if (-not $found) {
        Write-Host "⚠️  未找到已安装的 Claude Stats" -ForegroundColor Yellow
    } else {
        Write-Host "✅ 卸载完成!" -ForegroundColor Green
        Write-Host "💡 PATH中的条目需要手动清理（如有需要）" -ForegroundColor Yellow
    }
}

# 主逻辑
if ($Help) {
    Show-Help
    exit 0
}

if ($Uninstall) {
    Uninstall-ClaudeStats
    exit 0
}

# 检查PowerShell版本
if ($PSVersionTable.PSVersion.Major -lt 5) {
    Write-Host "❌ 需要 PowerShell 5.0 或更高版本" -ForegroundColor Red
    exit 1
}

# 执行安装
Install-ClaudeStats 
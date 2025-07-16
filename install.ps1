# claude-stats Windows安装脚本
param(
    [string]$InstallPath = "$env:USERPROFILE\bin"
)

$BinaryName = "claude-stats.exe"

Write-Host "📦 安装 claude-stats 到 Windows..." -ForegroundColor Green

# 检查是否已构建
if (-not (Test-Path $BinaryName)) {
    Write-Host "⚠️  未找到 $BinaryName" -ForegroundColor Yellow
    Write-Host "请先运行构建脚本:"
    Write-Host "  .\build.bat"
    exit 1
}

# 创建安装目录
if (-not (Test-Path $InstallPath)) {
    Write-Host "📁 创建目录: $InstallPath" -ForegroundColor Blue
    New-Item -ItemType Directory -Path $InstallPath -Force | Out-Null
}

# 复制文件
Write-Host "📋 复制文件到: $InstallPath" -ForegroundColor Blue
Copy-Item $BinaryName -Destination $InstallPath -Force

# 添加到PATH (如果不存在)
$currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($currentPath -notlike "*$InstallPath*") {
    Write-Host "🔧 添加到用户PATH..." -ForegroundColor Blue
    $newPath = "$currentPath;$InstallPath"
    [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
    Write-Host "⚠️  请重启命令行以使PATH生效" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "✅ 安装完成！" -ForegroundColor Green
Write-Host ""
Write-Host "🎉 现在您可以在任何地方使用:" -ForegroundColor Green
Write-Host "   claude-stats analyze"
Write-Host "   claude-stats analyze --help"
Write-Host ""
Write-Host "📍 安装位置: $InstallPath\$BinaryName" -ForegroundColor Blue
Write-Host ""
Write-Host "🗑️  卸载方法:" -ForegroundColor Red
Write-Host "   Remove-Item '$InstallPath\$BinaryName'" 
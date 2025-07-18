# Claude Stats 本地构建脚本 (Windows PowerShell)
# 用于在本地快速构建所有平台版本

param(
    [string]$Version = "v1.0.1"
)

Write-Host "🚀 开始构建 Claude Stats..." -ForegroundColor Green

# 获取版本信息
$BuildTime = (Get-Date).ToUniversalTime().ToString("yyyy-MM-dd_HH:mm:ss")
try {
    $GitCommit = (git rev-parse --short HEAD 2>$null)
} catch {
    $GitCommit = "unknown"
}

Write-Host "📦 版本: $Version" -ForegroundColor Cyan
Write-Host "⏰ 构建时间: $BuildTime" -ForegroundColor Cyan
Write-Host "🔗 Git提交: $GitCommit" -ForegroundColor Cyan

# 创建构建目录
$BuildDir = "build"
if (!(Test-Path $BuildDir)) {
    New-Item -ItemType Directory -Path $BuildDir
}
Set-Location $BuildDir

Write-Host ""
Write-Host "🔨 开始构建各平台二进制文件..." -ForegroundColor Yellow

# 构建函数
function Build-Platform {
    param(
        [string]$Goos,
        [string]$Goarch
    )
    
    Write-Host "  构建 $Goos/$Goarch..." -ForegroundColor White
    
    if ($Goos -eq "windows") {
        $BinaryName = "claude-stats-$Goos-$Goarch.exe"
    } else {
        $BinaryName = "claude-stats-$Goos-$Goarch"
    }
    
    $env:GOOS = $Goos
    $env:GOARCH = $Goarch
    
    go build -ldflags="-X main.Version=$Version -X main.BuildTime=$BuildTime -X main.GitCommit=$GitCommit" -o $BinaryName ../
    
    # 检查构建是否成功
    if (Test-Path $BinaryName) {
        Write-Host "    ✅ 构建成功: $BinaryName" -ForegroundColor Green
        
        # 创建压缩包
        if ($Goos -eq "windows") {
            $ZipName = "$($BinaryName.Replace('.exe', '')).zip"
            Compress-Archive -Path $BinaryName, ../README.md, ../LICENSE -DestinationPath $ZipName -Force
            Write-Host "    📦 已创建: $ZipName" -ForegroundColor Cyan
        } else {
            # Windows上创建tar.gz需要额外工具，这里简化为zip
            $ZipName = "$BinaryName.zip"
            Compress-Archive -Path $BinaryName, ../README.md, ../LICENSE -DestinationPath $ZipName -Force
            Write-Host "    📦 已创建: $ZipName" -ForegroundColor Cyan
        }
    } else {
        Write-Host "    ❌ 构建失败: $BinaryName" -ForegroundColor Red
        Write-Host "    请检查Go环境和依赖是否正确安装" -ForegroundColor Yellow
    }
}

# 构建主要平台
Build-Platform "linux" "amd64"
Build-Platform "windows" "amd64"
Build-Platform "darwin" "amd64"

Write-Host ""
Write-Host "🎉 构建完成！" -ForegroundColor Green
Write-Host ""
Write-Host "📂 构建产物位于 build/ 目录：" -ForegroundColor Cyan
Get-ChildItem

Write-Host ""
Write-Host "🧪 快速测试（Windows版本）：" -ForegroundColor Yellow
Write-Host "  .\claude-stats-windows-amd64.exe --version"
Write-Host "  .\claude-stats-windows-amd64.exe analyze --help"

Write-Host ""
Write-Host "💡 使用说明：" -ForegroundColor Magenta
Write-Host "  1. 解压对应平台的压缩包"
Write-Host "  2. 运行对应的二进制文件"
Write-Host "  3. 享受强大的Claude使用统计功能！"

# 提供全局安装选项
Write-Host ""
Write-Host "🌍 想要全局安装吗？" -ForegroundColor Green
Write-Host "运行以下命令将 claude-stats 安装为全局命令：" -ForegroundColor White
Write-Host "  ..\install-global.ps1" -ForegroundColor Cyan
Write-Host ""
Write-Host "安装后可在任何位置使用：" -ForegroundColor White
Write-Host "  claude-stats analyze" -ForegroundColor Yellow
Write-Host "  claude-stats --version" -ForegroundColor Yellow

# 清理环境变量
Remove-Item Env:GOOS -ErrorAction SilentlyContinue
Remove-Item Env:GOARCH -ErrorAction SilentlyContinue 
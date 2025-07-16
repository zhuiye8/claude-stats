@echo off
REM claude-stats Windows构建脚本

echo 🔨 构建claude-stats (Windows)...

REM 设置构建参数
set BINARY_NAME=claude-stats
set VERSION=1.0.0

REM 构建当前平台
echo 构建Windows版本...
go build -ldflags="-X main.Version=%VERSION%" -o %BINARY_NAME%.exe .

if %ERRORLEVEL% EQU 0 (
    echo ✅ 构建成功！
    echo 📍 二进制文件: %BINARY_NAME%.exe
    echo.
    echo 💡 使用方法:
    echo    %BINARY_NAME%.exe analyze
    echo    %BINARY_NAME%.exe analyze --help
) else (
    echo ❌ 构建失败！
    exit /b 1
)

echo.
echo 📦 想要构建所有平台版本？运行: build-all.bat
pause 
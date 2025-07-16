@echo off
REM claude-stats 全平台构建脚本 (Windows)

echo 🚀 构建所有平台的claude-stats...

set BINARY_NAME=claude-stats
set VERSION=1.0.0

REM 创建dist目录
if not exist dist mkdir dist

REM 构建Windows版本
echo 构建 Windows amd64...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-X main.Version=%VERSION%" -o dist/%BINARY_NAME%-windows-amd64.exe .

REM 构建Linux版本
echo 构建 Linux amd64...
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-X main.Version=%VERSION%" -o dist/%BINARY_NAME%-linux-amd64 .

REM 构建macOS版本
echo 构建 macOS amd64...
set GOOS=darwin
set GOARCH=amd64
go build -ldflags="-X main.Version=%VERSION%" -o dist/%BINARY_NAME%-darwin-amd64 .

REM 构建macOS ARM版本
echo 构建 macOS arm64...
set GOOS=darwin
set GOARCH=arm64
go build -ldflags="-X main.Version=%VERSION%" -o dist/%BINARY_NAME%-darwin-arm64 .

REM 重置环境变量
set GOOS=
set GOARCH=

echo.
echo ✅ 所有平台构建完成！
echo 📁 文件位于 dist/ 目录:
dir dist\%BINARY_NAME%-*

pause 
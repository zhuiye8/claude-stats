# 🚀 GitHub Release 发布指南

## 📋 一键发布流程

### 步骤1: 确保代码已推送
```bash
# 确保所有更改都已提交和推送
git add .
git commit -m "准备发布 v1.0.0"
git push origin main
```

### 步骤2: 创建并推送版本标签
```bash
# 创建版本标签 (重要：必须以 v 开头)
git tag v1.0.0

# 推送标签到GitHub (这会触发自动构建)
git push origin v1.0.0
```

### 步骤3: 等待自动构建完成
- 推送标签后，GitHub Actions会自动开始构建
- 访问 `https://github.com/zhuiye8/claude-stats/actions` 查看构建进度
- 构建完成后会自动创建Release并上传所有平台的二进制文件

## 🔄 自动构建的文件列表

构建完成后会生成以下文件：

### 📦 二进制文件
- `claude-stats-windows-amd64.exe` - Windows 64位
- `claude-stats-linux-amd64` - Linux 64位  
- `claude-stats-linux-arm64` - Linux ARM64
- `claude-stats-darwin-amd64` - macOS Intel
- `claude-stats-darwin-arm64` - macOS Apple Silicon

### 📁 压缩包
- `claude-stats-windows-amd64.zip` - Windows版本+文档
- `claude-stats-linux-amd64.tar.gz` - Linux版本+文档
- `claude-stats-linux-arm64.tar.gz` - Linux ARM版本+文档
- `claude-stats-darwin-amd64.tar.gz` - macOS Intel版本+文档
- `claude-stats-darwin-arm64.tar.gz` - macOS ARM版本+文档

## 🎯 首次发布步骤

### 1. 设置GitHub仓库
确保您的GitHub仓库有以下文件：
- `.github/workflows/release.yml` ✅ (已创建)
- `README.md` ✅ (已更新)
- `LICENSE` ✅ (已创建)

### 2. 执行发布命令
```bash
# 在项目根目录执行
git tag v1.0.0
git push origin v1.0.0
```

### 3. 检查构建状态
1. 访问 https://github.com/zhuiye8/claude-stats/actions
2. 查看"发布新版本"工作流是否成功
3. 构建成功后，访问 https://github.com/zhuiye8/claude-stats/releases

### 4. 验证下载链接
构建完成后，以下链接应该可以正常工作：

**Windows:**
```
https://github.com/zhuiye8/claude-stats/releases/download/v1.0.0/claude-stats-windows-amd64.exe
```

**macOS (Intel):**
```
https://github.com/zhuiye8/claude-stats/releases/download/v1.0.0/claude-stats-darwin-amd64
```

**macOS (Apple Silicon):**
```
https://github.com/zhuiye8/claude-stats/releases/download/v1.0.0/claude-stats-darwin-arm64
```

**Linux:**
```
https://github.com/zhuiye8/claude-stats/releases/download/v1.0.0/claude-stats-linux-amd64
```

## 🔧 发布新版本

### 修正版本 (v1.0.1)
```bash
git tag v1.0.1
git push origin v1.0.1
```

### 新功能版本 (v1.1.0)
```bash
git tag v1.1.0
git push origin v1.1.0
```

### 重大版本 (v2.0.0)
```bash
git tag v2.0.0
git push origin v2.0.0
```

## 🛠️ 手动构建 (备用方案)

如果自动构建失败，可以手动构建：

```bash
# 安装依赖
go mod tidy

# 构建所有平台 (使用Makefile)
make build-all

# 或者使用脚本
./build-all.bat    # Windows
make build-all     # Unix系统
```

## ✅ 验证发布成功

发布成功后，用户应该能够：

1. **直接下载**: 使用README中的curl命令下载
2. **查看Release页面**: 访问 https://github.com/zhuiye8/claude-stats/releases
3. **运行工具**: 下载后直接运行 `./claude-stats analyze`

## 🎉 发布完成！

一旦推送了标签，整个过程就会自动完成。大约5-10分钟后，用户就可以从GitHub Releases页面下载预编译的二进制文件了！

## 🔍 故障排除

### 构建失败？
1. 检查 `.github/workflows/release.yml` 文件是否正确
2. 查看Actions页面的错误日志
3. 确保go.mod文件正确

### 下载链接404？
1. 确认Release已创建：https://github.com/zhuiye8/claude-stats/releases
2. 检查文件名是否正确
3. 等待几分钟让CDN更新

### 权限问题？
确保您有仓库的管理员权限，能够创建Release。 
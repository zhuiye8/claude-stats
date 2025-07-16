# GitHub Actions 构建问题解决方案

## 🚨 常见构建错误

### 1. 账户计费问题
```
The job was not started because recent account payments have failed or your spending limit needs to be increased.
```

**问题原因：**
- GitHub免费账户有GitHub Actions使用限制（每月2000分钟）
- 构建多平台二进制文件消耗较多时间
- 可能超出免费账户的使用配额

**解决方案：**

#### 方案A：检查GitHub账户设置
1. 登录GitHub，访问：Settings → Billing and plans
2. 检查Actions的使用情况和限制
3. 如果超出限制，可以考虑升级到GitHub Pro

#### 方案B：使用优化后的构建配置
我已经将构建从6个平台同时构建改为分步构建，大幅减少资源消耗：
- Linux AMD64（最常用）
- Windows AMD64
- macOS AMD64（Intel）

#### 方案C：本地构建（推荐）
如果GitHub Actions持续出现问题，可以使用本地构建：

##### 使用一键构建脚本（推荐）

**Linux/macOS：**
```bash
# 1. 克隆项目
git clone https://github.com/zhuiye8/claude-stats.git
cd claude-stats

# 2. 运行一键构建脚本
./build-local.sh

# 或指定版本
./build-local.sh v1.0.2
```

**Windows：**
```powershell
# 1. 克隆项目
git clone https://github.com/zhuiye8/claude-stats.git
cd claude-stats

# 2. 运行PowerShell构建脚本
.\build-local.ps1

# 或指定版本
.\build-local.ps1 -Version "v1.0.2"
```

##### 手动构建单个平台
```bash
# 构建当前平台版本
go build -o claude-stats .

# 手动构建所有平台版本
# Windows
GOOS=windows GOARCH=amd64 go build -o claude-stats-windows-amd64.exe .

# Linux
GOOS=linux GOARCH=amd64 go build -o claude-stats-linux-amd64 .

# macOS
GOOS=darwin GOARCH=amd64 go build -o claude-stats-darwin-amd64 .
```

构建完成后，所有二进制文件和压缩包将在`build/`目录中。

### 2. 构建矩阵失败
```
The strategy configuration was canceled because "build.linux_amd64" failed
```

**问题原因：**
- 原始配置同时构建太多平台组合
- 资源竞争导致某个构建失败

**解决方案：**
已优化为分步构建，每个平台独立构建，避免相互影响。

## 🔄 重新尝试发布流程

### 使用优化后的GitHub Actions

1. **提交更新的配置：**
```bash
git add .github/workflows/release.yml
git commit -m "优化GitHub Actions配置，减少资源消耗"
git push origin main
```

2. **创建新版本标签：**
```bash
git tag v1.0.1
git push origin v1.0.1
```

3. **等待构建完成**（约3-5分钟）

### 本地构建发布包

如果GitHub Actions仍有问题，使用本地构建：

```bash
# 创建发布目录
mkdir release
cd release

# 构建所有平台
GOOS=windows GOARCH=amd64 go build -ldflags="-X main.Version=v1.0.1" -o claude-stats-windows-amd64.exe ../
GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Version=v1.0.1" -o claude-stats-linux-amd64 ../
GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.Version=v1.0.1" -o claude-stats-darwin-amd64 ../

# 创建压缩包
zip claude-stats-windows-amd64.zip claude-stats-windows-amd64.exe ../README.md ../LICENSE
tar -czf claude-stats-linux-amd64.tar.gz claude-stats-linux-amd64 ../README.md ../LICENSE
tar -czf claude-stats-darwin-amd64.tar.gz claude-stats-darwin-amd64 ../README.md ../LICENSE
```

然后手动上传到GitHub Release页面。

## 💡 长期解决方案

### 1. GitHub Pro账户（推荐）
- 每月3000分钟的Actions时间
- 私有仓库无限制
- 更好的支持和稳定性

### 2. 使用其他CI/CD服务
- **GitLab CI/CD**：免费账户有400分钟/月
- **Azure DevOps**：免费账户有1800分钟/月
- **Gitea Actions**：自托管，无限制

### 3. 设置本地发布脚本
创建自动化的本地发布脚本，避免依赖云端构建。

## 🆘 如果仍有问题

1. **检查GitHub Status**：https://www.githubstatus.com/
2. **联系GitHub Support**：如果是账户相关问题
3. **使用本地构建**：最可靠的方案
4. **降级到单平台构建**：只构建最需要的平台

## 📝 快速测试

验证本地构建是否正常：

```bash
# 构建当前平台版本
go build -o test-claude-stats .

# 测试基本功能
./test-claude-stats --version
./test-claude-stats analyze --help

# 清理测试文件
rm test-claude-stats
```

---

**记住：本地构建始终是最可靠的选择，不依赖任何外部服务！** 
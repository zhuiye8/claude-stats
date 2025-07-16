# 🐧 WSL 环境设置指南

## 🚨 问题解决

### 问题1: "go: cannot find main module"

**原因:** Go模块依赖未正确初始化

**解决方案:**
```bash
# 在claude-stats目录中运行
cd ~/work/claude-stats

# 重新初始化Go模块依赖
go mod tidy

# 然后重新构建
./build-local.sh
```

### 问题2: "Permission denied" 

**解决方案:**
```bash
# 给构建脚本添加执行权限
chmod +x build-local.sh

# 然后运行
./build-local.sh
```

## 📋 完整的WSL设置流程

### 1. 确保Go环境正确
```bash
# 检查Go版本
go version
# 应该显示 go1.21+ 或更高版本

# 如果Go版本过低，更新Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
```

### 2. 克隆并设置项目
```bash
# 克隆项目
git clone https://github.com/zhuiye8/claude-stats.git
cd claude-stats

# 确保有执行权限
chmod +x build-local.sh

# 初始化Go依赖
go mod tidy
```

### 3. 构建项目
```bash
# 运行构建脚本
./build-local.sh

# 或者手动构建当前平台
go build -o claude-stats .
```

### 4. 测试运行
```bash
# 查看帮助
./build/claude-stats-linux-amd64 --help

# 运行分析（如果有Claude日志）
./build/claude-stats-linux-amd64 analyze
```

## 🐛 常见问题排查

### 问题: 构建失败，提示"package not found"
```bash
# 清理模块缓存
go clean -modcache

# 重新下载依赖
go mod download

# 重新构建
./build-local.sh
```

### 问题: 权限错误
```bash
# 检查文件权限
ls -la build-local.sh

# 如果没有执行权限，添加权限
chmod +x build-local.sh
```

### 问题: Git相关错误
```bash
# 确保在正确的目录
pwd
# 应该显示 /root/work/claude-stats 或类似路径

# 检查Git状态
git status

# 如果需要，重新克隆
cd ..
rm -rf claude-stats
git clone https://github.com/zhuiye8/claude-stats.git
cd claude-stats
```

## ✅ 验证安装

运行以下命令验证一切正常：

```bash
# 检查Go环境
go version

# 检查项目结构
ls -la

# 尝试构建
go build -o test-claude-stats .

# 运行测试
./test-claude-stats --version

# 清理测试文件
rm test-claude-stats
```

## 🎯 WSL特定优化

### 1. 设置Go环境变量（如果需要）
```bash
# 添加到 ~/.bashrc
echo 'export PATH=/usr/local/go/bin:$PATH' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export GOROOT=/usr/local/go' >> ~/.bashrc

# 重新加载
source ~/.bashrc
```

### 2. 性能优化
```bash
# 如果构建较慢，可以启用Go模块代理
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=sum.golang.google.cn

# 或者添加到 ~/.bashrc 永久生效
echo 'export GOPROXY=https://goproxy.cn,direct' >> ~/.bashrc
echo 'export GOSUMDB=sum.golang.google.cn' >> ~/.bashrc
```

### 3. 跨文件系统路径
WSL可以访问Windows文件系统，Claude日志可能在：
```bash
# Windows路径映射
/mnt/c/Users/[用户名]/AppData/Roaming/claude/projects

# 使用示例
./claude-stats analyze /mnt/c/Users/[用户名]/AppData/Roaming/claude/projects
```

## 🆘 如果仍有问题

1. **查看详细错误信息:**
```bash
go build -v .
```

2. **检查Go模块状态:**
```bash
go mod verify
go list -m all
```

3. **重新开始:**
```bash
cd ..
rm -rf claude-stats
git clone https://github.com/zhuiye8/claude-stats.git
cd claude-stats
go mod tidy
chmod +x build-local.sh
./build-local.sh
```

---

**记住：WSL是Linux环境，所以使用Linux相关的命令和路径！** 🐧 
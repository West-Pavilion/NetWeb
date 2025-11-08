# NetWeb - Network Testing Tool

[![CI](https://github.com/YOUR_USERNAME/NetWeb/workflows/CI/badge.svg)](https://github.com/West-Pavilion/NetWeb/actions)
[![Release](https://github.com/West-Pavilion/NetWeb/workflows/Release%20to%20GitHub%20Packages/badge.svg)](https://github.com/YOUR_USERNAME/NetWeb/actions)

一个功能强大的网络连通性测试工具，使用 Go 语言构建后端 API，React 构建前端界面。

## 功能特性

- **cURL 测试**: 发送 HTTP/HTTPS 请求并查看完整响应
- **Ping 测试**: 测试主机可达性和网络延迟
- **Traceroute 测试**: 追踪数据包到达目标的路径
- **自定义命令**: 支持运行自定义网络诊断命令
- **实时结果**: 显示详细的连接信息和命令输出
- **跨平台**: 自动适配 Windows 和 Linux 系统命令

## 技术栈

### 后端
- Go 1.25
- Gorilla Mux (路由)
- 标准库 (net/http, os/exec)

### 前端
- React 18
- CSS3 (渐变背景和现代 UI)
- Fetch API

## 项目结构

```
NetWeb/
├── main.go                 # Go 后端服务器
├── go.mod                  # Go 依赖管理
├── go.sum                  # Go 依赖锁定
├── README.md               # 项目文档
└── frontend/               # React 前端
    ├── package.json        # Node 依赖
    ├── public/
    │   └── index.html      # HTML 模板
    └── src/
        ├── index.js        # React 入口
        ├── index.css       # 全局样式
        ├── App.js          # 主应用组件
        └── App.css         # 应用样式
```

## 快速开始

### 使用 Docker（推荐）

最快的方式是使用 Docker 运行应用：

#### 从 GitHub Packages 拉取镜像

```bash
# 拉取最新版本
docker pull ghcr.io/west-pavilion/netweb:latest

# 运行容器
docker run -d -p 8080:8080 --name netweb ghcr.io/west-pavilion/netweb:latest
```

访问 `http://localhost:8080` 即可使用！

#### 使用 Docker Compose

```bash
# 克隆仓库
git clone https://github.com/West-Pavilion/NetWeb.git
cd NetWeb

# 使用 Docker Compose 启动
docker-compose up -d
```

#### 本地构建 Docker 镜像

```bash
# 构建镜像
docker build -t netweb .

# 运行容器
docker run -d -p 8080:8080 --name netweb netweb
```

### 从源代码运行

#### 前置要求

- Go 1.25 或更高版本
- Node.js 14+ 和 npm
- 确保系统已安装以下命令工具:
  - `curl`
  - `ping`
  - `tracert` (Windows) 或 `traceroute` (Linux/Mac)

### 1. 安装 Go 依赖

```bash
go mod tidy
```

### 2. 安装前端依赖

```bash
cd frontend
npm install
```

### 3. 开发模式运行

#### 启动后端服务器 (终端 1)

```bash
go run main.go
```

服务器将在 `http://localhost:8080` 启动

#### 启动前端开发服务器 (终端 2)

```bash
cd frontend
npm start
```

前端开发服务器将在 `http://localhost:3000` 启动，并自动代理 API 请求到后端。

### 4. 生产模式部署

#### 构建前端

```bash
cd frontend
npm run build
```

#### 运行服务器

```bash
go run main.go
```

访问 `http://localhost:8080` 即可使用完整应用（后端会自动服务前端静态文件）。

## API 端点

### POST /api/test

执行网络测试命令

**请求体:**

```json
{
  "command": "curl|ping|tracert|custom",
  "url": "https://example.com",
  "custom": "nslookup {url}"  // 仅当 command 为 "custom" 时需要
}
```

**响应:**

```json
{
  "success": true,
  "command": "curl",
  "output": "命令输出内容...",
  "error": "",
  "duration": "1.234s",
  "connection": {
    "target": "https://example.com",
    "timestamp": "2025-11-09T00:00:00Z",
    "os": "windows"
  },
  "metadata": {
    "command_type": "curl",
    "execution_time": "1.234s"
  }
}
```

### GET /api/health

健康检查端点

**响应:**

```json
{
  "status": "ok",
  "time": "2025-11-09T00:00:00Z"
}
```

## 使用示例

### cURL 测试

1. 选择 "cURL - HTTP Request"
2. 输入 URL: `https://www.google.com`
3. 点击 "Run Test"
4. 查看完整的 HTTP 响应头和响应体

### Ping 测试

1. 选择 "Ping - ICMP Echo"
2. 输入主机: `google.com` 或 `8.8.8.8`
3. 点击 "Run Test"
4. 查看往返时间和丢包率

### Traceroute 测试

1. 选择 "Traceroute - Path Tracing"
2. 输入主机: `google.com`
3. 点击 "Run Test"
4. 查看数据包经过的每一跳路由

### 自定义命令

1. 选择 "Custom Command"
2. 输入自定义命令: `nslookup {url}` 或 `dig {url}`
3. 输入目标: `google.com`
4. 点击 "Run Test"
5. 查看命令执行结果

注意: 在自定义命令中使用 `{url}` 作为占位符，会被实际输入的 URL 替换。

## 安全注意事项

- 自定义命令功能应谨慎使用
- 建议在受控环境中运行此工具
- 生产环境部署时建议添加身份验证
- 考虑添加命令白名单以限制可执行的命令

## 超时设置

- cURL: 30 秒
- Ping: 15 秒
- Traceroute: 60 秒
- Custom: 30 秒

## CI/CD 和自动部署

本项目已配置完整的 GitHub Actions CI/CD 流程。

### 持续集成 (CI)

当代码推送到 `main` 或 `develop` 分支，或创建 Pull Request 时，会自动运行：

1. **后端构建**: 编译 Go 代码并运行 `go vet`
2. **前端构建**: 构建 React 应用
3. **集成测试**: 启动服务器并运行端到端测试

查看 CI 状态：[Actions 页面](https://github.com/West-Pavilion/NetWeb/actions)

### 发布流程

#### 创建新版本

1. 创建并推送标签：

```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

2. GitHub Actions 会自动：
   - 构建多平台 Docker 镜像（amd64, arm64）
   - 推送到 GitHub Container Registry
   - 为 Linux、Windows、macOS 构建二进制文件
   - 创建 GitHub Release 并上传构建产物

#### 使用发布的版本

**Docker 镜像：**

```bash
# 拉取特定版本
docker pull ghcr.io/west-pavilion/netweb:v1.0.1

# 运行
docker run -d -p 8080:8080 ghcr.io/west-pavilion/netweb:v1.0.1
```

**二进制文件：**

从 [Releases 页面](https://github.com/West-Pavilion/NetWeb/releases) 下载适合你系统的预编译二进制文件。

### 手动触发发布

也可以在 Actions 页面手动触发发布工作流：

1. 访问 [Actions](https://github.com/West-Pavilion/NetWeb/actions)
2. 选择 "Release to GitHub Packages"
3. 点击 "Run workflow"

## 故障排除

### 端口已被占用

如果 8080 端口已被占用，修改 `main.go` 中的端口号:

```go
port := ":8080"  // 改为其他端口，如 ":3001"
```

### 命令未找到

确保系统已安装相应的网络工具:

- Windows: curl, ping, tracert 应该已预装
- Linux: 使用 `apt install curl iputils-ping traceroute` (Debian/Ubuntu)
- Mac: 使用 `brew install curl` (curl 和 ping 已预装)

### CORS 错误

开发模式下，确保 `frontend/package.json` 中配置了正确的代理:

```json
"proxy": "http://localhost:8080"
```

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

## 作者

Built with Go and React
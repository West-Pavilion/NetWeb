# GitHub Actions Workflows

本目录包含 NetWeb 项目的所有自动化工作流。

## 工作流概览

### 1. CI Workflow (`ci.yml`)

**触发条件：**
- 推送到 `main` 或 `develop` 分支
- 创建针对 `main` 或 `develop` 的 Pull Request

**工作流程：**

1. **build-backend** - 构建 Go 后端
   - 设置 Go 1.21 环境
   - 缓存 Go 模块
   - 下载并验证依赖
   - 编译二进制文件
   - 运行 `go vet` 检查
   - 上传构建产物

2. **build-frontend** - 构建 React 前端
   - 设置 Node.js 18 环境
   - 缓存 npm 依赖
   - 安装依赖 (`npm ci`)
   - 构建生产版本
   - 上传构建产物

3. **integration-test** - 集成测试
   - 下载后端和前端构建产物
   - 启动服务器
   - 测试健康检查端点
   - 测试 API 功能

### 2. Release Workflow (`release.yml`)

**触发条件：**
- 推送符合 `v*.*.*` 格式的标签 (例如: `v1.0.0`)
- 手动触发 (workflow_dispatch)

**工作流程：**

1. **build-and-push** - 构建并推送 Docker 镜像
   - 设置 Docker Buildx (多平台构建)
   - 登录到 GitHub Container Registry (ghcr.io)
   - 提取元数据和标签
   - 构建 Linux amd64 和 arm64 镜像
   - 推送到 ghcr.io/west-pavilion/netweb
   - 创建 GitHub Release 和发布说明

2. **build-binaries** - 构建多平台二进制文件
   - 为以下平台构建：
     - Linux (amd64, arm64)
     - Windows (amd64)
     - macOS (amd64, arm64)
   - 创建压缩包 (.tar.gz / .zip)
   - 上传到 GitHub Release

## 使用指南

### 开发流程

1. 创建功能分支：
   ```bash
   git checkout -b feature/my-feature
   ```

2. 提交更改并推送：
   ```bash
   git add .
   git commit -m "Add new feature"
   git push origin feature/my-feature
   ```

3. 创建 Pull Request 到 `develop` 分支
   - CI 工作流会自动运行
   - 确保所有检查通过后再合并

### 发布新版本

1. 确保 `main` 分支是最新的：
   ```bash
   git checkout main
   git pull origin main
   ```

2. 创建版本标签：
   ```bash
   git tag -a v1.0.0 -m "Release version 1.0.0"
   ```

3. 推送标签：
   ```bash
   git push origin v1.0.0
   ```

4. Release 工作流会自动：
   - 构建 Docker 镜像
   - 发布到 GitHub Packages
   - 构建二进制文件
   - 创建 GitHub Release

### 手动触发发布

1. 访问 [Actions 页面](https://github.com/West-Pavilion/NetWeb/actions)
2. 选择 "Release to GitHub Packages"
3. 点击 "Run workflow"
4. 选择分支并确认

## 所需权限

工作流需要以下权限（已在 workflow 文件中配置）：

- `contents: read` - 读取仓库内容
- `packages: write` - 推送到 GitHub Packages

## Secrets 配置

工作流使用以下 secrets（由 GitHub 自动提供）：

- `GITHUB_TOKEN` - 用于认证 GitHub API 和 Container Registry

## Docker 镜像标签策略

发布工作流会自动创建以下标签：

- `latest` - 最新版本
- `v1.0.0` - 完整版本号
- `v1.0` - 主版本.次版本
- `v1` - 主版本
- `sha-abc123` - Git commit SHA

## 故障排除

### CI 失败

**后端构建失败：**
- 检查 Go 代码语法错误
- 运行 `go vet ./...` 本地检查
- 确保 `go.mod` 和 `go.sum` 是最新的

**前端构建失败：**
- 检查 React 代码语法错误
- 运行 `npm run build` 本地测试
- 确保 `package-lock.json` 已提交

**集成测试失败：**
- 检查 API 端点是否正常工作
- 确保服务器启动成功
- 查看工作流日志了解详细错误

### Release 失败

**Docker 构建失败：**
- 检查 Dockerfile 语法
- 本地运行 `docker build .` 测试
- 确保所有依赖文件都已包含

**二进制构建失败：**
- 检查跨平台编译兼容性
- 确保前端构建正常
- 验证文件路径是否正确

## 优化建议

### 加速构建

1. **缓存依赖**：工作流已配置缓存
   - Go 模块缓存
   - npm 依赖缓存
   - Docker 层缓存

2. **并行执行**：独立任务并行运行
   - 后端和前端构建同时进行
   - 多平台二进制并行构建

### 节省资源

- 只在必要时运行完整工作流
- Pull Request 可以只运行 CI
- 发布流程只在打标签时触发

## 监控和通知

- 在 GitHub Actions 页面查看工作流状态
- 失败时会在 Actions 标签页显示
- 可以配置 Slack/Email 通知（需额外设置）

## 更多资源

- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [Docker Buildx 文档](https://docs.docker.com/buildx/working-with-buildx/)
- [GitHub Container Registry 文档](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
# 🤝 贡献指南 | Contributing Guide

感谢你对 ServerPulse 的关注！本文档将帮助你了解如何参与本项目。

Thank you for your interest in ServerPulse! This document will help you understand how to participate in this project.

---

## 目录 | Table of Contents

- [行为准则 | Code of Conduct](#行为准则--code-of-conduct)
- [如何提交 Issue | How to Submit an Issue](#如何提交-issue--how-to-submit-an-issue)
- [如何提交 PR | How to Submit a Pull Request](#如何提交-pr--how-to-submit-a-pull-request)
- [代码规范 | Code Standards](#代码规范--code-standards)
- [开发环境搭建 | Development Environment Setup](#开发环境搭建--development-environment-setup)
- [提交信息规范 | Commit Message Convention](#提交信息规范--commit-message-convention)
- [项目结构 | Project Structure](#项目结构--project-structure)

---

## 行为准则 | Code of Conduct

- 尊重所有贡献者，保持友好和专业的沟通态度
- 接受建设性的代码审查意见，以改进代码质量
- 关注项目的整体一致性，而非个人偏好
- 对新手友好，耐心解答问题

- Respect all contributors and maintain a friendly, professional communication tone
- Accept constructive code review feedback to improve code quality
- Focus on overall project consistency rather than personal preferences
- Be beginner-friendly and patient when answering questions

---

## 如何提交 Issue | How to Submit an Issue

### Bug 报告 | Bug Reports

提交 Bug 之前，请先搜索现有 Issue，避免重复提交。提交时请包含以下信息：

Before submitting a bug, please search existing issues to avoid duplicates. When submitting, please include the following information:

1. **环境信息 | Environment Info**
   - 操作系统及版本 (OS and version)
   - Go 版本 (Go version)
   - ServerPulse 版本 (ServerPulse version)

2. **问题描述 | Problem Description**
   - 清晰描述遇到的问题 (Clearly describe the problem)
   - 预期行为 vs 实际行为 (Expected behavior vs actual behavior)

3. **复现步骤 | Steps to Reproduce**
   ```bash
   # 提供可复现问题的最小命令或代码
   # Provide the minimal command or code to reproduce the issue
   serverpulse check --cpu
   ```

4. **日志输出 | Log Output**
   - 附上相关的错误日志或控制台输出 (Attach relevant error logs or console output)

### 功能建议 | Feature Requests

提交功能建议时，请说明：

When submitting a feature request, please describe:

- **需求背景 | Background**: 为什么需要这个功能？(Why is this feature needed?)
- **使用场景 | Use Case**: 具体的使用场景 (Specific use cases)
- **预期行为 | Expected Behavior**: 期望的功能表现 (Expected behavior of the feature)
- **可能的实现方案 | Possible Implementation**: 如果有想法，可以提供实现思路 (If you have ideas, you can suggest an implementation approach)

---

## 如何提交 PR | How to Submit a Pull Request

### PR 提交流程 | PR Submission Workflow

```
1. Fork 本仓库 | Fork this repository
2. 创建特性分支 | Create a feature branch
3. 进行开发并提交 | Develop and commit
4. 推送到你的 Fork | Push to your fork
5. 创建 Pull Request | Create a Pull Request
6. 等待代码审查 | Wait for code review
7. 根据反馈修改 | Make changes based on feedback
8. 合并 | Merge
```

### 分支命名规范 | Branch Naming Convention

| 类型 | 格式 | 示例 |
|------|------|------|
| 新功能 | `feature/xxx` | `feature/web-dashboard` |
| Bug 修复 | `fix/xxx` | `fix/memory-leak` |
| 重构 | `refactor/xxx` | `refactor/collector-optimization` |
| 文档 | `docs/xxx` | `docs/api-documentation` |
| 测试 | `test/xxx` | `test/anomaly-detector` |

### PR 描述模板 | PR Description Template

```markdown
## 变更类型 | Change Type
- [ ] 新功能 | New Feature
- [ ] Bug 修复 | Bug Fix
- [ ] 重构 | Refactor
- [ ] 文档更新 | Documentation
- [ ] 测试 | Testing

## 变更说明 | Description
<!-- 描述你的变更内容 | Describe your changes -->

## 关联 Issue | Related Issue
<!-- 关联的 Issue 编号 | Related issue number -->
Closes #xxx

## 测试 | Testing
<!-- 描述如何测试这些变更 | Describe how to test these changes -->
- [ ] 单元测试通过 | Unit tests pass
- [ ] 集成测试通过 | Integration tests pass
- [ ] 手动测试通过 | Manual testing pass
```

---

## 代码规范 | Code Standards

### Go 代码规范 | Go Code Standards

1. **格式化 | Formatting**: 使用 `gofmt` 或 `goimports` 格式化代码
   ```bash
   gofmt -w .
   goimports -w .
   ```

2. **命名规范 | Naming Conventions**:
   - 包名：小写单词，不使用下划线 (Package names: lowercase, no underscores)
   - 导出标识符：大驼峰 (Exported identifiers: PascalCase)
   - 非导出标识符：小驼峰 (Unexported identifiers: camelCase)
   - 常量：大驼峰或全大写+下划线 (Constants: PascalCase or UPPER_SNAKE_CASE)
   - 接口：单方法接口以 `-er` 结尾 (Single-method interfaces end with `-er`)

3. **注释规范 | Comment Conventions**:
   - 所有导出的类型、函数、常量必须有注释 (All exported types, functions, and constants must have comments)
   - 注释以标识符名称开头 (Comments start with the identifier name)
   ```go
   // Collector 定义系统指标采集器接口。
   // Collector defines the system metrics collector interface.
   type Collector interface {
       // Collect 采集系统指标并返回结果。
       // Collect gathers system metrics and returns the result.
       Collect() (*Metrics, error)
   }
   ```

4. **错误处理 | Error Handling**:
   - 始终检查并处理错误 (Always check and handle errors)
   - 使用 `fmt.Errorf("context: %w", err)` 包装错误 (Wrap errors with `fmt.Errorf`)
   - 避免忽略错误 (Avoid ignoring errors)

5. **并发安全 | Concurrency Safety**:
   - 共享数据必须使用互斥锁保护 (Shared data must be protected with mutexes)
   - 优先使用 channel 进行通信 (Prefer channels for communication)
   - 避免在 goroutine 中直接使用外部变量 (Avoid using external variables directly in goroutines)

### 项目特定规范 | Project-Specific Standards

1. **模块化设计 | Modular Design**: 每个功能模块放在 `internal/` 对应目录下 (Each functional module goes in its corresponding directory under `internal/`)
2. **接口抽象 | Interface Abstraction**: 采集器、探测器等使用接口定义，便于扩展和测试 (Collectors, probers, etc. use interface definitions for easy extension and testing)
3. **配置驱动 | Configuration-Driven**: 可配置参数通过 `config.yaml` 暴露，避免硬编码 (Configurable parameters are exposed through `config.yaml`, avoid hardcoding)
4. **日志规范 | Logging Standards**: 使用结构化日志，包含上下文信息 (Use structured logging with context information)

---

## 开发环境搭建 | Development Environment Setup

### 前置要求 | Prerequisites

- **Go** >= 1.21
- **Git**
- **Make** (可选，推荐 / Optional but recommended)
- **Docker** (可选，用于容器化测试 / Optional, for containerized testing)

### 搭建步骤 | Setup Steps

```bash
# 1. 克隆仓库 | Clone the repository
git clone https://github.com/gitstq/ServerPulse.git
cd ServerPulse

# 2. 安装依赖 | Install dependencies
go mod download

# 3. 编译项目 | Build the project
make build

# 4. 验证安装 | Verify installation
./bin/serverpulse version

# 5. 运行测试 | Run tests
go test ./...

# 6. 运行 linter | Run linter
golangci-lint run
```

### 开发常用命令 | Common Development Commands

```bash
# 编译 | Build
make build

# 运行 | Run
go run ./cmd/serverpulse

# 测试 | Test
go test ./...                    # 运行所有测试 | Run all tests
go test ./internal/collector/... # 运行特定包测试 | Run specific package tests
go test -v -race ./...           # 详细输出 + 竞态检测 | Verbose + race detection

# 代码检查 | Lint
golangci-lint run

# 格式化 | Format
gofmt -w .
goimports -w .
```

---

## 提交信息规范 | Commit Message Convention

本项目采用 **Angular Convention** 提交信息规范。

This project follows the **Angular Convention** for commit messages.

### 格式 | Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type 类型说明 | Type Descriptions

| Type | 说明 | Description |
|------|------|-------------|
| `feat` | 新功能 | New feature |
| `fix` | Bug 修复 | Bug fix |
| `docs` | 文档变更 | Documentation changes |
| `style` | 代码格式（不影响逻辑） | Code formatting (no logic change) |
| `refactor` | 重构（非新功能、非修复） | Refactoring (neither feature nor fix) |
| `perf` | 性能优化 | Performance improvement |
| `test` | 测试相关 | Test-related |
| `build` | 构建系统或外部依赖 | Build system or external dependencies |
| `ci` | CI 配置 | CI configuration |
| `chore` | 其他杂项 | Other chores |
| `revert` | 回滚提交 | Revert a commit |

### Scope 范围说明 | Scope Descriptions

| Scope | 说明 | Description |
|-------|------|-------------|
| `collector` | 系统指标采集 | System metrics collection |
| `probe` | 多协议探测 | Multi-protocol probing |
| `anomaly` | 异常检测 | Anomaly detection |
| `alert` | 告警管理 | Alert management |
| `tui` | TUI 仪表板 | TUI dashboard |
| `server` | Hub/Agent 服务端 | Hub/Agent server |
| `storage` | 数据存储 | Data storage |
| `snapshot` | 快照功能 | Snapshot feature |
| `config` | 配置管理 | Configuration management |
| `cli` | 命令行接口 | Command line interface |

### 示例 | Examples

```
feat(collector): add GPU metrics collection support

feat(probe): add HTTP health check with custom headers support

fix(alert): resolve alert aggregation window not resetting after notification

fix(tui): correct memory display overflow on servers with >100GB RAM

docs(readme): update installation instructions for v1.0.0

refactor(storage): migrate to WAL mode for better concurrent read performance

perf(collector): optimize CPU metrics collection to reduce syscall overhead

test(anomaly): add unit tests for Z-Score detector with edge cases

chore(deps): bump gopsutil to v3.24.0
```

### 提交信息规则 | Commit Message Rules

1. **subject** 使用祈使句，首字母小写，不加句号 (Use imperative mood, lowercase first letter, no period)
2. **subject** 不超过 72 个字符 (Subject must not exceed 72 characters)
3. **body** 详细描述变更原因和影响 (Body describes the reason and impact of changes in detail)
4. **footer** 关联 Issue 编号 (Footer references related issue numbers)
   ```
   Closes #123
   ```
5. 破坏性变更需在 body 和 footer 中注明 (Breaking changes must be noted in body and footer)
   ```
   feat(api)!: change report endpoint request format

   BREAKING CHANGE: The /api/v1/report endpoint now requires
   a JSON body with the new metrics format. See migration guide.
   ```

---

## 项目结构 | Project Structure

```
ServerPulse/
├── cmd/
│   └── serverpulse/       # 程序入口 | Application entry point
│       └── main.go
├── configs/
│   └── serverpulse.yaml   # 默认配置文件 | Default configuration file
├── internal/              # 内部包（不可外部导入）| Internal packages (not importable externally)
│   ├── alert/             # 告警管理 | Alert management
│   ├── anomaly/           # 异常检测 | Anomaly detection
│   ├── collector/         # 系统指标采集 | System metrics collection
│   ├── probe/             # 多协议探测 | Multi-protocol probing
│   ├── server/            # Hub/Agent 服务端 | Hub/Agent server
│   ├── snapshot/          # 快照功能 | Snapshot feature
│   ├── storage/           # 数据存储 | Data storage
│   └── tui/               # TUI 仪表板 | TUI dashboard
├── pkg/
│   └── config/            # 公共配置 | Public configuration
├── assets/                # 静态资源 | Static assets
├── docs/                  # 文档 | Documentation
├── go.mod                 # Go 模块定义 | Go module definition
├── go.sum                 # 依赖校验 | Dependency checksums
├── Makefile               # 构建脚本 | Build scripts
├── LICENSE                # 开源协议 | License
└── README.md              # 项目说明 | Project description
```

---

<div align="center">
  <sub>感谢所有贡献者的付出！| Thank you to all contributors!</sub>
</div>

# 📖 详细使用指南 | Detailed Usage Guide

本文档提供 ServerPulse 所有子命令的详细用法、配置参数说明、告警规则配置示例以及 Hub+Agent 部署指南。

This document provides detailed usage for all ServerPulse subcommands, configuration parameter descriptions, alert rule configuration examples, and the Hub+Agent deployment guide.

---

## 目录 | Table of Contents

- [命令行概览 | CLI Overview](#命令行概览--cli-overview)
- [子命令详解 | Subcommand Details](#子命令详解--subcommand-details)
  - [serverpulse tui — TUI 仪表板](#serverpulse-tui--tui-仪表板--tui-dashboard)
  - [serverpulse check — 系统检查](#serverpulse-check--系統檢查--system-check)
  - [serverpulse serve — Hub 服务端](#serverpulse-serve--hub-服務端--hub-server)
  - [serverpulse agent — Agent 客户端](#serverpulse-agent--agent-客戶端--agent-client)
  - [serverpulse snapshot — 快照管理](#serverpulse-snapshot--快照管理--snapshot-management)
  - [serverpulse probe — 协议探测](#serverpulse-probe--協議探測--protocol-probing)
  - [serverpulse version — 版本信息](#serverpulse-version--版本信息--version-info)
- [配置文件说明 | Configuration File Reference](#配置文件說明--configuration-file-reference)
- [告警规则配置 | Alert Rule Configuration](#告警規則配置--alert-rule-configuration)
- [Hub+Agent 部署指南 | Hub+Agent Deployment Guide](#hubagent-部署指南--hubagent-deployment-guide)
- [REST API 文档 | REST API Documentation](#rest-api-文檔--rest-api-documentation)

---

## 命令行概览 | CLI Overview

```bash
serverpulse [command] [flags]

# 查看全局帮助 | View global help
serverpulse --help

# 查看子命令帮助 | View subcommand help
serverpulse tui --help
serverpulse check --help
serverpulse serve --help
serverpulse agent --help
serverpulse snapshot --help
serverpulse probe --help
```

### 全局参数 | Global Flags

| 参数 | 缩写 | 说明 | 默认值 |
|------|------|------|--------|
| `--config` | `-c` | 指定配置文件路径 | `~/.serverpulse/config.yaml` |
| `--verbose` | `-v` | 启用详细日志输出 | `false` |
| `--quiet` | `-q` | 静默模式，仅输出错误 | `false` |

---

## 子命令详解 | Subcommand Details

### serverpulse tui — TUI 仪表板 | TUI Dashboard

启动基于终端的实时监控仪表板。

Launch the terminal-based real-time monitoring dashboard.

```bash
serverpulse tui [flags]
```

**参数 | Flags:**

| 参数 | 缩写 | 说明 | 默认值 |
|------|------|------|--------|
| `--interval` | `-i` | 数据刷新间隔 | `2s` |
| `--theme` | `-t` | 主题颜色 (dark/light) | `dark` |
| `--no-color` | | 禁用颜色输出 | `false` |

**使用示例 | Examples:**

```bash
# 默认启动 | Default launch
serverpulse tui

# 自定义刷新间隔 | Custom refresh interval
serverpulse tui --interval 5s

# 使用浅色主题 | Use light theme
serverpulse tui --theme light
```

**键盘快捷键 | Keyboard Shortcuts:**

| 按键 | 功能 | Action |
|------|------|--------|
| `Tab` | 切换视图面板 | Switch view panels |
| `↑` / `↓` | 上下浏览 | Scroll up/down |
| `←` / `→` | 左右切换标签 | Switch tabs left/right |
| `Enter` | 查看详情 | View details |
| `Esc` | 返回上一层 | Go back |
| `q` | 退出 | Quit |
| `?` | 帮助 | Help |

**视图说明 | View Descriptions:**

- **Dashboard（仪表板）**: CPU、内存、磁盘、网络的实时使用率图表
- **Alerts（告警面板）**: 当前活跃告警列表，按严重程度排序
- **Servers（服务器视图）**: 多服务器状态一览（Hub 模式下可用）
- **Processes（进程视图）**: 系统进程资源占用排行

---

### serverpulse check — 系统检查 | System Check

快速检查系统各项指标状态。

Quick check of system metric statuses.

```bash
serverpulse check [flags]
```

**参数 | Flags:**

| 参数 | 缩写 | 说明 | 默认值 |
|------|------|------|--------|
| `--cpu` | | 检查 CPU 指标 | `false` |
| `--memory` | | 检查内存指标 | `false` |
| `--disk` | | 检查磁盘指标 | `false` |
| `--network` | | 检查网络指标 | `false` |
| `--docker` | | 检查 Docker 容器状态 | `false` |
| `--gpu` | | 检查 GPU 指标 | `false` |
| `--all` | `-a` | 检查所有指标 | `false` |
| `--json` | `-j` | 以 JSON 格式输出 | `false` |
| `--warn-threshold` | | 告警阈值百分比 | `80` |

**使用示例 | Examples:**

```bash
# 检查所有指标 | Check all metrics
serverpulse check --all

# 仅检查 CPU 和内存 | Check only CPU and memory
serverpulse check --cpu --memory

# JSON 格式输出（适合脚本解析）| JSON output (suitable for script parsing)
serverpulse check --all --json

# 自定义告警阈值 | Custom alert threshold
serverpulse check --all --warn-threshold 90
```

**输出示例 | Output Example:**

```
=== ServerPulse System Check ===

CPU:
  Usage:     23.5%        [OK]
  Cores:     8
  Load Avg:  1.2 / 0.8 / 0.5

Memory:
  Usage:     67.2%        [OK]
  Available: 5.2 GB / 16.0 GB
  Swap:      0.1 GB / 4.0 GB

Disk (/):
  Usage:     45.3%        [OK]
  Total:     500 GB
  Free:      273.5 GB

Network (eth0):
  In:        12.5 MB/s
  Out:       3.2 MB/s
  TCP Conn:  128

Docker:
  Running:   5 containers
  Stopped:   2 containers
  Images:    12

All checks passed. ✓
```

---

### serverpulse serve — Hub 服务端 | Hub Server

启动 Hub 中心服务端，接收 Agent 上报的数据，提供 REST API。

Start the Hub central server to receive data reported by Agents and provide REST API.

```bash
serverpulse serve [flags]
```

**参数 | Flags:**

| 参数 | 缩写 | 说明 | 默认值 |
|------|------|------|--------|
| `--port` | `-p` | HTTP 服务端口 | `8080` |
| `--host` | | 绑定地址 | `0.0.0.0` |
| `--data-dir` | `-d` | 数据存储目录 | `~/.serverpulse/data` |
| `--api-key` | | API 认证密钥 | (空，不启用认证) |
| `--tls-cert` | | TLS 证书路径 | (空，不启用 TLS) |
| `--tls-key` | | TLS 私钥路径 | (空，不启用 TLS) |
| `--max-servers` | | 最大注册服务器数 | `100` |
| `--retention` | | 数据保留天数 | `30` |

**使用示例 | Examples:**

```bash
# 基本启动 | Basic launch
serverpulse serve --port 8080

# 启用 API 认证 | Enable API authentication
serverpulse serve --port 8080 --api-key "your-secret-key"

# 启用 TLS | Enable TLS
serverpulse serve --port 8443 --tls-cert /path/to/cert.pem --tls-key /path/to/key.pem

# 自定义数据目录和保留策略 | Custom data directory and retention policy
serverpulse serve --data-dir /data/serverpulse --retention 60
```

---

### serverpulse agent — Agent 客户端 | Agent Client

启动 Agent 客户端，采集本地系统指标并上报到 Hub。

Start the Agent client to collect local system metrics and report to the Hub.

```bash
serverpulse agent [flags]
```

**参数 | Flags:**

| 参数 | 缩写 | 说明 | 默认值 |
|------|------|------|--------|
| `--hub` | `-H` | Hub 服务地址 | `http://localhost:8080` |
| `--interval` | `-i` | 数据上报间隔 | `10s` |
| `--api-key` | | Hub API 认证密钥 | (空) |
| `--hostname` | | 自定义主机名标识 | (自动获取系统主机名) |
| `--tags` | | 自定义标签（逗号分隔）| (空) |
| `--collectors` | | 启用的采集器（逗号分隔）| `cpu,memory,disk,network` |
| `--no-docker` | | 禁用 Docker 采集 | `false` |
| `--no-gpu` | | 禁用 GPU 采集 | `false` |

**使用示例 | Examples:**

```bash
# 连接到 Hub | Connect to Hub
serverpulse agent --hub http://192.168.1.100:8080

# 自定义上报间隔和标签 | Custom interval and tags
serverpulse agent --hub http://192.168.1.100:8080 --interval 30s --tags "env:prod,team:backend"

# 仅采集 CPU 和内存 | Collect only CPU and memory
serverpulse agent --hub http://192.168.1.100:8080 --collectors cpu,memory

# 使用 API 认证 | Use API authentication
serverpulse agent --hub http://192.168.1.100:8080 --api-key "your-secret-key"
```

---

### serverpulse snapshot — 快照管理 | Snapshot Management

创建和管理系统状态快照，支持快照间对比。

Create and manage system state snapshots with comparison support.

```bash
serverpulse snapshot <subcommand> [flags]
```

**子命令 | Subcommands:**

#### `snapshot create` — 创建快照 | Create Snapshot

```bash
serverpulse snapshot create [flags]
```

| 参数 | 缩写 | 说明 | 默认值 |
|------|------|------|--------|
| `--name` | `-n` | 快照名称 | 自动生成时间戳名称 |
| `--description` | `-d` | 快照描述 | (空) |
| `--include-procs` | | 包含进程列表 | `false` |

```bash
# 创建命名快照 | Create a named snapshot
serverpulse snapshot create --name "baseline"

# 创建带描述的快照 | Create a snapshot with description
serverpulse snapshot create --name "before-deploy" --description "系统部署前基线"

# 创建包含进程信息的快照 | Create a snapshot with process info
serverpulse snapshot create --name "full" --include-procs
```

#### `snapshot list` — 列出快照 | List Snapshots

```bash
serverpulse snapshot list [flags]
```

| 参数 | 缩写 | 说明 | 默认值 |
|------|------|------|--------|
| `--limit` | `-l` | 显示数量 | `20` |
| `--json` | `-j` | JSON 格式输出 | `false` |

```bash
# 列出所有快照 | List all snapshots
serverpulse snapshot list

# JSON 格式输出 | JSON output
serverpulse snapshot list --json
```

#### `snapshot compare` — 对比快照 | Compare Snapshots

```bash
serverpulse snapshot compare <snapshot1> <snapshot2> [flags]
```

| 参数 | 缩写 | 说明 | 默认值 |
|------|------|------|--------|
| `--json` | `-j` | JSON 格式输出 | `false` |

```bash
# 对比两个快照 | Compare two snapshots
serverpulse snapshot compare baseline current

# JSON 格式输出差异 | JSON output for differences
serverpulse snapshot compare baseline current --json
```

#### `snapshot delete` — 删除快照 | Delete Snapshot

```bash
serverpulse snapshot delete <snapshot-name>
```

```bash
serverpulse snapshot delete old-baseline
```

---

### serverpulse probe — 协议探测 | Protocol Probing

对目标主机执行多协议探测。

Perform multi-protocol probing against target hosts.

```bash
serverpulse probe <target> [flags]
```

**参数 | Flags:**

| 参数 | 缩写 | 说明 | 默认值 |
|------|------|------|--------|
| `--type` | `-t` | 探测类型 (ping/tcp/http/ssh) | `ping` |
| `--port` | `-p` | 目标端口（TCP/HTTP/SSH）| `80` |
| `--count` | `-c` | Ping 探测次数 | `4` |
| `--timeout` | | 超时时间 | `5s` |
| `--path` | | HTTP 探测路径 | `/` |
| `--method` | | HTTP 请求方法 | `GET` |
| `--headers` | | HTTP 自定义请求头 | (空) |

**使用示例 | Examples:**

```bash
# Ping 探测 | Ping probe
serverpulse probe example.com --type ping --count 5

# TCP 端口探测 | TCP port probe
serverpulse probe example.com --type tcp --port 443

# HTTP 健康检查 | HTTP health check
serverpulse probe example.com --type http --port 8080 --path /health

# SSH 连通性探测 | SSH connectivity probe
serverpulse probe example.com --type ssh --port 22

# 带自定义请求头的 HTTP 探测 | HTTP probe with custom headers
serverpulse probe example.com --type http --headers "Authorization: Bearer token"
```

---

### serverpulse version — 版本信息 | Version Info

显示当前版本和构建信息。

Display current version and build information.

```bash
serverpulse version

# 输出示例 | Output example:
# ServerPulse v1.0.0
# Go Version: go1.21.0
# OS/Arch: linux/amd64
# Commit: abc1234
# Built: 2024-01-01T00:00:00Z
```

---

## 配置文件说明 | Configuration File Reference

ServerPulse 使用 YAML 格式的配置文件，默认路径为 `~/.serverpulse/config.yaml`。

ServerPulse uses a YAML-format configuration file, with the default path `~/.serverpulse/config.yaml`.

### 完整配置示例 | Full Configuration Example

```yaml
# ServerPulse 配置文件 | ServerPulse Configuration File

# 全局设置 | Global Settings
verbose: false
quiet: false
log-level: "info"           # debug / info / warn / error
log-file: ""                # 日志文件路径（空则输出到标准输出）| Log file path (empty for stdout)

# 采集器配置 | Collector Configuration
collector:
  interval: "10s"           # 采集间隔 | Collection interval
  cpu: true                 # 启用 CPU 采集 | Enable CPU collection
  memory: true              # 启用内存采集 | Enable memory collection
  disk: true                # 启用磁盘采集 | Enable disk collection
  network: true             # 启用网络采集 | Enable network collection
  docker: true              # 启用 Docker 采集 | Enable Docker collection
  gpu: false                # 启用 GPU 采集 | Enable GPU collection
  disk-partitions:          # 监控的磁盘分区 | Monitored disk partitions
    - "/"
    - "/data"
  network-interfaces:       # 监控的网络接口 | Monitored network interfaces
    - "eth0"
    - "lo"

# 异常检测配置 | Anomaly Detection Configuration
anomaly:
  enabled: true             # 启用异常检测 | Enable anomaly detection
  sensitivity: "medium"    # 灵敏度: low / medium / high | Sensitivity
  window-size: 60           # 滑动窗口大小（数据点数）| Sliding window size (data points)
  z-score-threshold: 2.5    # Z-Score 阈值 | Z-Score threshold

# 告警配置 | Alert Configuration
alert:
  enabled: true             # 启用告警 | Enable alerts
  aggregation-window: "60s" # 告警聚合时间窗口 | Alert aggregation time window
  cooldown: "300s"         # 同类告警冷却时间 | Cooldown for same-type alerts
  rules:                    # 告警规则（详见下方）| Alert rules (see below)
    - name: "high-cpu"
      metric: "cpu"
      threshold: 90
      duration: "60s"
      severity: "critical"
    - name: "high-memory"
      metric: "memory"
      threshold: 85
      duration: "120s"
      severity: "warning"
    - name: "disk-full"
      metric: "disk"
      threshold: 95
      duration: "300s"
      severity: "critical"

  # 通知渠道 | Notification Channels
  notifications:
    webhook:
      enabled: false
      url: ""
      method: "POST"
      headers:
        Content-Type: "application/json"

    email:
      enabled: false
      smtp-host: "smtp.gmail.com"
      smtp-port: 587
      username: ""
      password: ""
      from: ""
      to:
        - "admin@example.com"
      tls: true

# Hub 配置 | Hub Configuration
hub:
  host: "0.0.0.0"
  port: 8080
  data-dir: "~/.serverpulse/data"
  api-key: ""               # API 认证密钥（空则不启用）| API auth key (empty to disable)
  max-servers: 100
  retention-days: 30
  tls:
    enabled: false
    cert: ""
    key: ""

# Agent 配置 | Agent Configuration
agent:
  hub: "http://localhost:8080"
  api-key: ""               # Hub API 认证密钥 | Hub API auth key
  interval: "10s"
  hostname: ""              # 自定义主机名（空则自动获取）| Custom hostname (auto-detect if empty)
  tags:                     # 自定义标签 | Custom tags
    env: "production"
    team: "backend"
  collectors:
    - "cpu"
    - "memory"
    - "disk"
    - "network"
  docker: true
  gpu: false

# TUI 配置 | TUI Configuration
tui:
  theme: "dark"             # dark / light
  refresh-interval: "2s"
  no-color: false

# 快照配置 | Snapshot Configuration
snapshot:
  dir: "~/.serverpulse/snapshots"
  max-count: 50             # 最大快照数量 | Maximum snapshot count
  auto-cleanup: true        # 自动清理超限快照 | Auto-cleanup excess snapshots
```

### 配置参数速查表 | Configuration Parameter Quick Reference

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `verbose` | bool | `false` | 详细日志模式 |
| `quiet` | bool | `false` | 静默模式 |
| `log-level` | string | `"info"` | 日志级别 (debug/info/warn/error) |
| `log-file` | string | `""` | 日志文件路径 |
| `collector.interval` | duration | `"10s"` | 采集间隔 |
| `collector.cpu` | bool | `true` | CPU 采集开关 |
| `collector.memory` | bool | `true` | 内存采集开关 |
| `collector.disk` | bool | `true` | 磁盘采集开关 |
| `collector.network` | bool | `true` | 网络采集开关 |
| `collector.docker` | bool | `true` | Docker 采集开关 |
| `collector.gpu` | bool | `false` | GPU 采集开关 |
| `anomaly.enabled` | bool | `true` | 异常检测开关 |
| `anomaly.sensitivity` | string | `"medium"` | 灵敏度 (low/medium/high) |
| `anomaly.window-size` | int | `60` | 滑动窗口大小 |
| `anomaly.z-score-threshold` | float | `2.5` | Z-Score 阈值 |
| `alert.enabled` | bool | `true` | 告警开关 |
| `alert.aggregation-window` | duration | `"60s"` | 告警聚合窗口 |
| `alert.cooldown` | duration | `"300s"` | 告警冷却时间 |
| `hub.port` | int | `8080` | Hub 服务端口 |
| `hub.host` | string | `"0.0.0.0"` | Hub 绑定地址 |
| `hub.api-key` | string | `""` | API 认证密钥 |
| `hub.retention-days` | int | `30` | 数据保留天数 |
| `agent.hub` | string | `"http://localhost:8080"` | Hub 地址 |
| `agent.interval` | duration | `"10s"` | 上报间隔 |
| `tui.theme` | string | `"dark"` | TUI 主题 (dark/light) |
| `tui.refresh-interval` | duration | `"2s"` | TUI 刷新间隔 |

---

## 告警规则配置 | Alert Rule Configuration

### 规则结构 | Rule Structure

每条告警规则包含以下字段：

Each alert rule contains the following fields:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | 规则名称 |
| `metric` | string | 是 | 监控指标 (cpu/memory/disk) |
| `threshold` | float | 是 | 阈值百分比 (0-100) |
| `duration` | duration | 否 | 持续超阈值才触发（默认 `"0s"` 立即触发）|
| `severity` | string | 否 | 严重级别 (info/warning/critical，默认 `"warning"`) |
| `message` | string | 否 | 自定义告警消息模板 |
| `enabled` | bool | 否 | 是否启用（默认 `true`）|

### 告警规则示例 | Alert Rule Examples

```yaml
alert:
  rules:
    # CPU 使用率超过 90% 持续 60 秒，严重告警
    # CPU usage exceeds 90% for 60 seconds, critical alert
    - name: "high-cpu"
      metric: "cpu"
      threshold: 90
      duration: "60s"
      severity: "critical"
      message: "CPU usage is {{.value}}%, exceeding threshold {{.threshold}}%"

    # CPU 使用率超过 80% 持续 120 秒，警告
    # CPU usage exceeds 80% for 120 seconds, warning
    - name: "medium-cpu"
      metric: "cpu"
      threshold: 80
      duration: "120s"
      severity: "warning"

    # 内存使用率超过 85%，警告
    # Memory usage exceeds 85%, warning
    - name: "high-memory"
      metric: "memory"
      threshold: 85
      duration: "120s"
      severity: "warning"

    # 内存使用率超过 95%，严重告警
    # Memory usage exceeds 95%, critical alert
    - name: "critical-memory"
      metric: "memory"
      threshold: 95
      duration: "60s"
      severity: "critical"

    # 磁盘使用率超过 90%，警告
    # Disk usage exceeds 90%, warning
    - name: "disk-warning"
      metric: "disk"
      threshold: 90
      duration: "300s"
      severity: "warning"

    # 磁盘使用率超过 95%，严重告警
    # Disk usage exceeds 95%, critical alert
    - name: "disk-critical"
      metric: "disk"
      threshold: 95
      duration: "60s"
      severity: "critical"

    # 禁用的规则示例
    # Disabled rule example
    - name: "low-cpu"
      metric: "cpu"
      threshold: 10
      severity: "info"
      enabled: false
```

### Webhook 通知配置 | Webhook Notification Configuration

```yaml
alert:
  notifications:
    webhook:
      enabled: true
      url: "https://hooks.example.com/alert"
      method: "POST"
      headers:
        Content-Type: "application/json"
        Authorization: "Bearer your-webhook-token"
```

**Webhook 请求体格式 | Webhook Request Body Format:**

```json
{
  "alert_id": "alert-20240101-001",
  "rule_name": "high-cpu",
  "severity": "critical",
  "metric": "cpu",
  "value": 95.2,
  "threshold": 90,
  "hostname": "server-01",
  "message": "CPU usage is 95.2%, exceeding threshold 90%",
  "timestamp": "2024-01-01T12:00:00Z",
  "aggregated_count": 3
}
```

### 邮件通知配置 | Email Notification Configuration

```yaml
alert:
  notifications:
    email:
      enabled: true
      smtp-host: "smtp.gmail.com"
      smtp-port: 587
      username: "your-email@gmail.com"
      password: "your-app-password"
      from: "serverpulse@gmail.com"
      to:
        - "admin@example.com"
        - "oncall@example.com"
      tls: true
```

> **注意 | Note**: 使用 Gmail 时，建议使用应用专用密码而非账户密码。
>
> When using Gmail, it is recommended to use an app-specific password instead of your account password.

---

## Hub+Agent 部署指南 | Hub+Agent Deployment Guide

### 架构概览 | Architecture Overview

```
                    ┌─────────────────────┐
                    │       Hub           │
                    │  (中心服务端)         │
                    │  Port: 8080         │
                    │                     │
                    │  - 数据聚合          │
                    │  - 异常检测          │
                    │  - 告警管理          │
                    │  - REST API         │
                    │  - 数据存储 (SQLite) │
                    └──────────┬──────────┘
                               │
              ┌────────────────┼────────────────┐
              │                │                │
    ┌─────────┴──────┐ ┌──────┴───────┐ ┌──────┴───────┐
    │   Agent (A)    │ │  Agent (B)   │ │  Agent (C)   │
    │  server-01    │ │  server-02   │ │  server-03   │
    │               │ │              │ │              │
    │ - CPU/Mem/Disk │ │ - CPU/Mem    │ │ - CPU/Mem    │
    │ - Network      │ │ - Disk       │ │ - Network    │
    │ - Docker       │ │ - Docker     │ │              │
    └───────────────┘ └──────────────┘ └──────────────┘
```

### 单机部署 | Standalone Deployment

最简单的部署方式，Hub 和 Agent 运行在同一台机器上。

The simplest deployment where Hub and Agent run on the same machine.

```bash
# 1. 安装 ServerPulse | Install ServerPulse
go install github.com/gitstq/ServerPulse/cmd/serverpulse@latest

# 2. 启动 Hub | Start Hub
serverpulse serve --port 8080 &

# 3. 启动 Agent 连接到本地 Hub | Start Agent connecting to local Hub
serverpulse agent --hub http://localhost:8080 --interval 10s &

# 4. 使用 TUI 查看监控 | Use TUI to view monitoring
serverpulse tui
```

### 多服务器部署 | Multi-Server Deployment

#### 步骤 1：部署 Hub（中心服务器）| Step 1: Deploy Hub (Central Server)

在中心服务器上执行：

Execute on the central server:

```bash
# 安装 | Install
go install github.com/gitstq/ServerPulse/cmd/serverpulse@latest

# 创建配置文件 | Create configuration file
mkdir -p ~/.serverpulse
cat > ~/.serverpulse/config.yaml << 'EOF'
hub:
  port: 8080
  api-key: "your-secure-api-key"
  retention-days: 30

alert:
  enabled: true
  notifications:
    webhook:
      enabled: true
      url: "https://hooks.example.com/alert"
EOF

# 启动 Hub | Start Hub
serverpulse serve --port 8080
```

#### 步骤 2：部署 Agent（被监控服务器）| Step 2: Deploy Agent (Monitored Servers)

在每台被监控的服务器上执行：

Execute on each monitored server:

```bash
# 安装 | Install
go install github.com/gitstq/ServerPulse/cmd/serverpulse@latest

# 创建配置文件 | Create configuration file
mkdir -p ~/.serverpulse
cat > ~/.serverpulse/config.yaml << 'EOF'
agent:
  hub: "http://hub-server:8080"
  api-key: "your-secure-api-key"
  interval: "10s"
  tags:
    env: "production"
    role: "web"
EOF

# 启动 Agent | Start Agent
serverpulse agent
```

#### 步骤 3：使用 systemd 管理（推荐）| Step 3: Manage with systemd (Recommended)

**Hub systemd 服务 | Hub systemd service:**

```ini
# /etc/systemd/system/serverpulse-hub.service
[Unit]
Description=ServerPulse Hub Server
After=network.target

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/serverpulse serve --port 8080 --config /root/.serverpulse/config.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

**Agent systemd 服务 | Agent systemd service:**

```ini
# /etc/systemd/system/serverpulse-agent.service
[Unit]
Description=ServerPulse Agent
After=network.target

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/serverpulse agent --config /root/.serverpulse/config.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

**启用并启动服务 | Enable and start services:**

```bash
# Hub 服务器上 | On the Hub server
sudo systemctl daemon-reload
sudo systemctl enable serverpulse-hub
sudo systemctl start serverpulse-hub

# Agent 服务器上 | On Agent servers
sudo systemctl daemon-reload
sudo systemctl enable serverpulse-agent
sudo systemctl start serverpulse-agent

# 查看状态 | Check status
sudo systemctl status serverpulse-hub
sudo systemctl status serverpulse-agent

# 查看日志 | View logs
sudo journalctl -u serverpulse-hub -f
sudo journalctl -u serverpulse-agent -f
```

---

## REST API 文档 | REST API Documentation

Hub 服务端提供以下 REST API 接口。

The Hub server provides the following REST API endpoints.

### 基础信息 | Base Info

- **Base URL**: `http://<hub-host>:<port>/api/v1`
- **认证方式**: `Authorization: Bearer <api-key>`（如果配置了 API Key）
- **Content-Type**: `application/json`

### API 端点 | API Endpoints

#### 1. 健康检查 | Health Check

检查 Hub 服务健康状态。

Check Hub service health status.

```
GET /api/v1/health
```

**响应示例 | Response Example:**

```json
{
  "status": "ok",
  "version": "1.0.0",
  "uptime": "72h15m30s",
  "servers_count": 5,
  "timestamp": "2024-01-01T12:00:00Z"
}
```

#### 2. 注册服务器 | Register Server

Agent 向 Hub 注册自身信息。

Agent registers itself with the Hub.

```
POST /api/v1/register
```

**请求体 | Request Body:**

```json
{
  "hostname": "server-01",
  "os": "linux",
  "arch": "amd64",
  "ip": "192.168.1.10",
  "tags": {
    "env": "production",
    "role": "web"
  }
}
```

**响应示例 | Response Example:**

```json
{
  "status": "ok",
  "server_id": "srv-abc123",
  "message": "Server registered successfully"
}
```

#### 3. 上报数据 | Report Data

Agent 向 Hub 上报采集数据。

Agent reports collected data to the Hub.

```
POST /api/v1/report
```

**请求体 | Request Body:**

```json
{
  "server_id": "srv-abc123",
  "hostname": "server-01",
  "timestamp": "2024-01-01T12:00:00Z",
  "metrics": {
    "cpu": {
      "usage": 23.5,
      "cores": 8,
      "load_avg_1": 1.2,
      "load_avg_5": 0.8,
      "load_avg_15": 0.5
    },
    "memory": {
      "usage": 67.2,
      "total": 17179869184,
      "available": 5620367872,
      "swap_usage": 2.5
    },
    "disk": [
      {
        "partition": "/",
        "usage": 45.3,
        "total": 536870912000,
        "free": 293601280000
      }
    ],
    "network": [
      {
        "interface": "eth0",
        "bytes_in": 13107200,
        "bytes_out": 3355443,
        "tcp_connections": 128
      }
    ]
  }
}
```

**响应示例 | Response Example:**

```json
{
  "status": "ok",
  "anomalies": [
    {
      "metric": "cpu",
      "value": 23.5,
      "z_score": 0.8,
      "is_anomaly": false
    }
  ]
}
```

#### 4. 获取服务器列表 | Get Server List

获取所有已注册的服务器列表。

Get the list of all registered servers.

```
GET /api/v1/servers
```

**响应示例 | Response Example:**

```json
{
  "servers": [
    {
      "server_id": "srv-abc123",
      "hostname": "server-01",
      "os": "linux",
      "arch": "amd64",
      "ip": "192.168.1.10",
      "status": "online",
      "last_report": "2024-01-01T12:00:00Z",
      "tags": {
        "env": "production",
        "role": "web"
      }
    },
    {
      "server_id": "srv-def456",
      "hostname": "server-02",
      "os": "linux",
      "arch": "arm64",
      "ip": "192.168.1.11",
      "status": "offline",
      "last_report": "2024-01-01T11:30:00Z",
      "tags": {
        "env": "production",
        "role": "database"
      }
    }
  ],
  "total": 2
}
```

#### 5. 获取服务器详情 | Get Server Details

获取指定服务器的详细指标数据。

Get detailed metrics for a specific server.

```
GET /api/v1/servers/{server_id}
```

**响应示例 | Response Example:**

```json
{
  "server_id": "srv-abc123",
  "hostname": "server-01",
  "status": "online",
  "last_report": "2024-01-01T12:00:00Z",
  "metrics": {
    "cpu": {
      "usage": 23.5,
      "cores": 8,
      "load_avg": [1.2, 0.8, 0.5]
    },
    "memory": {
      "usage": 67.2,
      "total_gb": 16.0,
      "available_gb": 5.2
    },
    "disk": [
      {
        "partition": "/",
        "usage": 45.3,
        "total_gb": 500,
        "free_gb": 273.5
      }
    ]
  },
  "alerts": [
    {
      "id": "alert-001",
      "rule": "high-cpu",
      "severity": "warning",
      "message": "CPU usage exceeded 80%",
      "triggered_at": "2024-01-01T11:45:00Z"
    }
  ]
}
```

#### 6. 获取告警列表 | Get Alerts

获取当前活跃告警列表。

Get the list of current active alerts.

```
GET /api/v1/alerts?severity=critical&server=srv-abc123
```

**查询参数 | Query Parameters:**

| 参数 | 说明 |
|------|------|
| `severity` | 按严重级别过滤 (info/warning/critical) |
| `server` | 按服务器 ID 过滤 |
| `limit` | 返回数量限制（默认 50）|

**响应示例 | Response Example:**

```json
{
  "alerts": [
    {
      "id": "alert-20240101-001",
      "server_id": "srv-abc123",
      "hostname": "server-01",
      "rule_name": "high-cpu",
      "severity": "critical",
      "metric": "cpu",
      "value": 95.2,
      "threshold": 90,
      "message": "CPU usage is 95.2%, exceeding threshold 90%",
      "triggered_at": "2024-01-01T12:00:00Z",
      "acknowledged": false
    }
  ],
  "total": 1
}
```

### 错误响应 | Error Responses

所有 API 在出错时返回统一格式：

All APIs return a unified format on error:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Server not found",
    "details": "Server ID 'srv-xxx' does not exist"
  }
}
```

**HTTP 状态码 | HTTP Status Codes:**

| 状态码 | 说明 |
|--------|------|
| `200 OK` | 请求成功 |
| `201 Created` | 资源创建成功 |
| `400 Bad Request` | 请求参数错误 |
| `401 Unauthorized` | 未授权（API Key 缺失或错误）|
| `404 Not Found` | 资源不存在 |
| `500 Internal Server Error` | 服务端内部错误 |

---

<div align="center">
  <sub>如有问题，请提交 <a href="https://github.com/gitstq/ServerPulse/issues">Issue</a></sub>
</div>

<div align="center">
  <img src="assets/logo.jpg" alt="ServerPulse Logo" width="120" height="120">

  # 🩺 ServerPulse

  **轻量级服务器智能监控引擎 | Lightweight Server Intelligent Monitoring Engine**

  [![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
  [![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
  [![Release](https://img.shields.io/badge/Release-v1.0.0-blue.svg)](https://github.com/gitstq/ServerPulse/releases/tag/v1.0.0)
  [![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey)]()
  [![Go Report Card](https://goreportcard.com/badge/github.com/gitstq/ServerPulse)](https://goreportcard.com/report/github.com/gitstq/ServerPulse)

  [简体中文](README.md) | [English](README_EN.md) | [繁體中文](README_TW.md)

  <p>
    <b>零依赖</b> · <b>AI异常检测</b> · <b>多协议探测</b> · <b>TUI仪表板</b> · <b>Hub+Agent架构</b> · <b>跨平台</b>
  </p>
</div>

---

## 🎉 项目介绍

ServerPulse 是一款轻量级、零外部依赖的服务器智能监控引擎。它集成了系统指标采集、多协议探测、AI异常检测、智能告警管理和实时TUI仪表板，为个人开发者和小团队提供开箱即用的服务器监控解决方案。

### 🎯 设计理念

- **极简主义**：零外部依赖，单二进制文件即可运行
- **智能驱动**：内置AI异常检测算法，自动识别性能异常
- **多服务器**：Hub+Agent架构，轻松管理多台服务器
- **终端优先**：精美的TUI仪表板，SSH连接即可查看
- **告警智能**：告警聚合防风暴，多通道通知

---

## ✨ 核心特性

### 📊 系统指标采集
- **CPU**：使用率、核心数、负载均值、温度
- **内存**：使用率、可用内存、Swap状态
- **磁盘**：分区使用率、I/O读写速率
- **网络**：带宽使用、连接数、流量统计
- **Docker**：容器状态、资源使用、镜像统计
- **GPU**：NVIDIA/AMD显卡使用率、显存、温度（如可用）

### 🔍 多协议探测
- **Ping**：ICMP连通性检测，延迟和丢包率
- **TCP**：端口可达性检查
- **HTTP**：健康检查端点，响应时间和状态码
- **SSH**：SSH服务连通性验证

### 🧠 AI异常检测
- **Z-Score算法**：基于统计学的异常值检测
- **移动平均**：平滑趋势分析，消除噪声干扰
- **可配置灵敏度**：Low / Medium / High 三级
- **自动学习**：根据历史数据动态调整基线

### 🚨 智能告警
- **阈值告警**：自定义CPU/内存/磁盘阈值
- **异常告警**：AI检测到异常自动触发
- **告警聚合**：60秒窗口内同类告警合并
- **多通道通知**：Webhook / 邮件（SMTP）

### 📺 TUI仪表板
- **实时仪表板**：CPU/内存/磁盘/网络实时图表
- **告警面板**：当前活跃告警列表
- **服务器视图**：多服务器状态一览
- **键盘导航**：Tab切换视图，方向键浏览

### 🏗️ Hub+Agent架构
- **Hub**：中心服务端，接收所有Agent数据，提供REST API
- **Agent**：轻量客户端，部署在被监控服务器上
- **REST API**：`/api/v1/register`、`/api/v1/report`、`/api/v1/servers`、`/api/v1/health`

### 📸 快照对比
- 创建系统状态快照
- 历史快照列表
- 两个快照间差异对比

---

## 🚀 快速开始

### 安装

**从源码编译：**
```bash
git clone https://github.com/gitstq/ServerPulse.git
cd ServerPulse
make build
sudo make install
```

**使用 Go Install：**
```bash
go install github.com/gitstq/ServerPulse/cmd/serverpulse@latest
```

**从 Release 下载：**
```bash
# 下载最新版本
wget https://github.com/gitstq/ServerPulse/releases/latest/download/serverpulse-linux-amd64
chmod +x serverpulse-linux-amd64
sudo mv serverpulse-linux-amd64 /usr/local/bin/serverpulse
```

### 快速使用

```bash
# 启动TUI仪表板（最常用）
serverpulse tui

# 快速系统检查
serverpulse check --all
serverpulse check --cpu --memory --json

# 启动Hub服务端
serverpulse serve --port 8080

# 启动Agent客户端
serverpulse agent --hub http://your-hub:8080 --interval 10s

# 创建状态快照
serverpulse snapshot create --name "baseline"
serverpulse snapshot list
serverpulse snapshot compare baseline current
```

---

## 📖 详细使用指南

详见 [docs/USAGE.md](docs/USAGE.md)

---

## 💡 设计思路与迭代规划

### 架构设计
ServerPulse 采用 Hub+Agent 分离架构：
- Agent 部署在被监控节点，负责数据采集和上报
- Hub 作为中心节点，负责数据聚合、存储、异常检测和告警
- 单机模式下，Agent 和 Hub 可在同一进程运行

### 技术选型
| 组件 | 技术 | 理由 |
|------|------|------|
| 核心语言 | Go | 高性能、跨平台编译、单二进制分发 |
| 系统指标 | gopsutil | 跨平台系统信息采集标准库 |
| 终端UI | Bubbletea | Go生态最成熟的TUI框架 |
| 存储 | SQLite | 零配置嵌入式数据库，WAL模式高性能 |
| CLI | Cobra | Go生态标准CLI框架 |

### 迭代规划
- **v1.1**：Web仪表板、Prometheus指标导出
- **v1.2**：分布式部署、集群监控
- **v2.0**：机器学习异常检测、预测性告警

---

## 📦 打包与部署指南

### 编译
```bash
make build          # 编译当前平台
make build-all      # 交叉编译全平台（Linux/macOS/Windows, amd64/arm64）
make install        # 安装到 $GOPATH/bin
```

### 配置
编辑 `~/.serverpulse/config.yaml` 或使用 `--config` 指定配置文件路径。

### Docker部署（规划中）
```bash
docker run -d --name serverpulse \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v ~/.serverpulse:/root/.serverpulse \
  gitstq/serverpulse:latest serve
```

---

## 🤝 贡献指南

详见 [CONTRIBUTING.md](CONTRIBUTING.md)

---

## 📄 开源协议

本项目基于 [MIT License](LICENSE) 开源。

---

<div align="center">
  <sub>Built with ❤️ by <a href="https://github.com/gitstq">gitstq</a></sub>
</div>

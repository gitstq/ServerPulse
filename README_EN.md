<div align="center">
  <img src="assets/logo.jpg" alt="ServerPulse Logo" width="120" height="120">

  # 🩺 ServerPulse

  **Lightweight Server Intelligent Monitoring Engine**

  [![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
  [![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
  [![Release](https://img.shields.io/badge/Release-v1.0.0-blue.svg)](https://github.com/gitstq/ServerPulse/releases/tag/v1.0.0)
  [![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey)]()
  [![Go Report Card](https://goreportcard.com/badge/github.com/gitstq/ServerPulse)](https://goreportcard.com/report/github.com/gitstq/ServerPulse)

  [简体中文](README.md) | [English](README_EN.md) | [繁體中文](README_TW.md)

  <p>
    <b>Zero Dependencies</b> · <b>AI Anomaly Detection</b> · <b>Multi-Protocol Probing</b> · <b>TUI Dashboard</b> · <b>Hub+Agent Architecture</b> · <b>Cross-Platform</b>
  </p>
</div>

---

## 🎉 Introduction

ServerPulse is a lightweight, zero-dependency intelligent server monitoring engine. It integrates system metrics collection, multi-protocol probing, AI anomaly detection, smart alert management, and a real-time TUI dashboard, providing an out-of-the-box server monitoring solution for individual developers and small teams.

### 🎯 Design Philosophy

- **Minimalism**: Zero external dependencies, runs as a single binary
- **Intelligence-Driven**: Built-in AI anomaly detection algorithms for automatic performance anomaly identification
- **Multi-Server**: Hub+Agent architecture for easy management of multiple servers
- **Terminal-First**: Beautiful TUI dashboard, viewable via SSH connection
- **Smart Alerts**: Alert aggregation to prevent storms, multi-channel notifications

---

## ✨ Core Features

### 📊 System Metrics Collection
- **CPU**: Usage, core count, load average, temperature
- **Memory**: Usage, available memory, swap status
- **Disk**: Partition usage, I/O read/write rates
- **Network**: Bandwidth usage, connection count, traffic statistics
- **Docker**: Container status, resource usage, image statistics
- **GPU**: NVIDIA/AMD GPU usage, VRAM, temperature (if available)

### 🔍 Multi-Protocol Probing
- **Ping**: ICMP connectivity detection, latency and packet loss rate
- **TCP**: Port reachability check
- **HTTP**: Health check endpoint, response time and status code
- **SSH**: SSH service connectivity verification

### 🧠 AI Anomaly Detection
- **Z-Score Algorithm**: Statistical-based outlier detection
- **Moving Average**: Smooth trend analysis to eliminate noise interference
- **Configurable Sensitivity**: Low / Medium / High levels
- **Auto-Learning**: Dynamic baseline adjustment based on historical data

### 🚨 Smart Alerts
- **Threshold Alerts**: Customizable CPU/memory/disk thresholds
- **Anomaly Alerts**: Automatically triggered when AI detects anomalies
- **Alert Aggregation**: Similar alerts merged within a 60-second window
- **Multi-Channel Notifications**: Webhook / Email (SMTP)

### 📺 TUI Dashboard
- **Real-Time Dashboard**: CPU/memory/disk/network real-time charts
- **Alert Panel**: Current active alert list
- **Server View**: Multi-server status overview
- **Keyboard Navigation**: Tab to switch views, arrow keys to browse

### 🏗️ Hub+Agent Architecture
- **Hub**: Central server that receives data from all Agents and provides REST API
- **Agent**: Lightweight client deployed on monitored servers
- **REST API**: `/api/v1/register`, `/api/v1/report`, `/api/v1/servers`, `/api/v1/health`

### 📸 Snapshot Comparison
- Create system state snapshots
- Historical snapshot list
- Diff comparison between two snapshots

---

## 🚀 Quick Start

### Installation

**Build from source:**
```bash
git clone https://github.com/gitstq/ServerPulse.git
cd ServerPulse
make build
sudo make install
```

**Using Go Install:**
```bash
go install github.com/gitstq/ServerPulse/cmd/serverpulse@latest
```

**Download from Release:**
```bash
# Download the latest version
wget https://github.com/gitstq/ServerPulse/releases/latest/download/serverpulse-linux-amd64
chmod +x serverpulse-linux-amd64
sudo mv serverpulse-linux-amd64 /usr/local/bin/serverpulse
```

### Quick Usage

```bash
# Launch TUI Dashboard (most common)
serverpulse tui

# Quick system check
serverpulse check --all
serverpulse check --cpu --memory --json

# Start Hub server
serverpulse serve --port 8080

# Start Agent client
serverpulse agent --hub http://your-hub:8080 --interval 10s

# Create state snapshots
serverpulse snapshot create --name "baseline"
serverpulse snapshot list
serverpulse snapshot compare baseline current
```

---

## 📖 Detailed Usage Guide

See [docs/USAGE.md](docs/USAGE.md) for details.

---

## 💡 Design & Roadmap

### Architecture Design
ServerPulse adopts a Hub+Agent separated architecture:
- Agent is deployed on monitored nodes, responsible for data collection and reporting
- Hub acts as the central node, responsible for data aggregation, storage, anomaly detection, and alerting
- In standalone mode, Agent and Hub can run in the same process

### Technology Stack
| Component | Technology | Rationale |
|-----------|-----------|-----------|
| Core Language | Go | High performance, cross-platform compilation, single binary distribution |
| System Metrics | gopsutil | Cross-platform system information collection standard library |
| Terminal UI | Bubbletea | The most mature TUI framework in the Go ecosystem |
| Storage | SQLite | Zero-configuration embedded database, high-performance WAL mode |
| CLI | Cobra | Standard CLI framework in the Go ecosystem |

### Roadmap
- **v1.1**: Web dashboard, Prometheus metrics export
- **v1.2**: Distributed deployment, cluster monitoring
- **v2.0**: Machine learning anomaly detection, predictive alerting

---

## 📦 Build & Deployment Guide

### Build
```bash
make build          # Build for current platform
make build-all      # Cross-compile for all platforms (Linux/macOS/Windows, amd64/arm64)
make install        # Install to $GOPATH/bin
```

### Configuration
Edit `~/.serverpulse/config.yaml` or use `--config` to specify the configuration file path.

### Docker Deployment (Planned)
```bash
docker run -d --name serverpulse \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v ~/.serverpulse:/root/.serverpulse \
  gitstq/serverpulse:latest serve
```

---

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

---

## 📄 License

This project is open-sourced under the [MIT License](LICENSE).

---

<div align="center">
  <sub>Built with ❤️ by <a href="https://github.com/gitstq">gitstq</a></sub>
</div>

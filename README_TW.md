<div align="center">
  <img src="assets/logo.jpg" alt="ServerPulse Logo" width="120" height="120">

  # 🩺 ServerPulse

  **輕量級伺服器智慧監控引擎 | Lightweight Server Intelligent Monitoring Engine**

  [![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
  [![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
  [![Release](https://img.shields.io/badge/Release-v1.0.0-blue.svg)](https://github.com/gitstq/ServerPulse/releases/tag/v1.0.0)
  [![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey)]()
  [![Go Report Card](https://goreportcard.com/badge/github.com/gitstq/ServerPulse)](https://goreportcard.com/report/github.com/gitstq/ServerPulse)

  [简体中文](README.md) | [English](README_EN.md) | [繁體中文](README_TW.md)

  <p>
    <b>零依賴</b> · <b>AI異常檢測</b> · <b>多協議探測</b> · <b>TUI儀表板</b> · <b>Hub+Agent架構</b> · <b>跨平台</b>
  </p>
</div>

---

## 🎉 專案介紹

ServerPulse 是一款輕量級、零外部依賴的伺服器智慧監控引擎。它整合了系統指標採集、多協議探測、AI異常檢測、智慧告警管理和即時TUI儀表板，為個人開發者和小團隊提供開箱即用的伺服器監控解決方案。

### 🎯 設計理念

- **極簡主義**：零外部依賴，單二進位檔即可運行
- **智慧驅動**：內建AI異常檢測演算法，自動識別效能異常
- **多伺服器**：Hub+Agent架構，輕鬆管理多台伺服器
- **終端優先**：精美的TUI儀表板，SSH連線即可查看
- **告警智慧**：告警聚合防風暴，多通道通知

---

## ✨ 核心特性

### 📊 系統指標採集
- **CPU**：使用率、核心數、負載均值、溫度
- **記憶體**：使用率、可用記憶體、Swap狀態
- **磁碟**：分區使用率、I/O讀寫速率
- **網路**：頻寬使用、連線數、流量統計
- **Docker**：容器狀態、資源使用、映像檔統計
- **GPU**：NVIDIA/AMD顯示卡使用率、VRAM、溫度（如可用）

### 🔍 多協議探測
- **Ping**：ICMP連通性檢測，延遲和丟包率
- **TCP**：埠可達性檢查
- **HTTP**：健康檢查端點，回應時間和狀態碼
- **SSH**：SSH服務連通性驗證

### 🧠 AI異常檢測
- **Z-Score演算法**：基於統計學的異常值檢測
- **移動平均**：平滑趨勢分析，消除雜訊干擾
- **可配置靈敏度**：Low / Medium / High 三級
- **自動學習**：根據歷史資料動態調整基線

### 🚨 智慧告警
- **閾值告警**：自訂CPU/記憶體/磁碟閾值
- **異常告警**：AI檢測到異常自動觸發
- **告警聚合**：60秒視窗內同類告警合併
- **多通道通知**：Webhook / 郵件（SMTP）

### 📺 TUI儀表板
- **即時儀表板**：CPU/記憶體/磁碟/網路即時圖表
- **告警面板**：當前活躍告警列表
- **伺服器檢視**：多伺服器狀態一覽
- **鍵盤導航**：Tab切換檢視，方向鍵瀏覽

### 🏗️ Hub+Agent架構
- **Hub**：中心伺服端，接收所有Agent資料，提供REST API
- **Agent**：輕量客戶端，部署在被監控伺服器上
- **REST API**：`/api/v1/register`、`/api/v1/report`、`/api/v1/servers`、`/api/v1/health`

### 📸 快照對比
- 建立系統狀態快照
- 歷史快照列表
- 兩個快照間差異對比

---

## 🚀 快速開始

### 安裝

**從原始碼編譯：**
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

**從 Release 下載：**
```bash
# 下載最新版本
wget https://github.com/gitstq/ServerPulse/releases/latest/download/serverpulse-linux-amd64
chmod +x serverpulse-linux-amd64
sudo mv serverpulse-linux-amd64 /usr/local/bin/serverpulse
```

### 快速使用

```bash
# 啟動TUI儀表板（最常用）
serverpulse tui

# 快速系統檢查
serverpulse check --all
serverpulse check --cpu --memory --json

# 啟動Hub伺服端
serverpulse serve --port 8080

# 啟動Agent客戶端
serverpulse agent --hub http://your-hub:8080 --interval 10s

# 建立狀態快照
serverpulse snapshot create --name "baseline"
serverpulse snapshot list
serverpulse snapshot compare baseline current
```

---

## 📖 詳細使用指南

詳見 [docs/USAGE.md](docs/USAGE.md)

---

## 💡 設計思路與迭代規劃

### 架構設計
ServerPulse 採用 Hub+Agent 分離架構：
- Agent 部署在被監控節點，負責資料採集和上報
- Hub 作為中心節點，負責資料聚合、儲存、異常檢測和告警
- 單機模式下，Agent 和 Hub 可在同一進程運行

### 技術選型
| 組件 | 技術 | 理由 |
|------|------|------|
| 核心語言 | Go | 高效能、跨平台編譯、單二進位檔分發 |
| 系統指標 | gopsutil | 跨平台系統資訊採集標準庫 |
| 終端UI | Bubbletea | Go生態最成熟的TUI框架 |
| 儲存 | SQLite | 零配置嵌入式資料庫，WAL模式高效能 |
| CLI | Cobra | Go生態標準CLI框架 |

### 迭代規劃
- **v1.1**：Web儀表板、Prometheus指標匯出
- **v1.2**：分散式部署、叢集監控
- **v2.0**：機器學習異常檢測、預測性告警

---

## 📦 打包與部署指南

### 編譯
```bash
make build          # 編譯當前平台
make build-all      # 交叉編譯全平台（Linux/macOS/Windows, amd64/arm64）
make install        # 安裝到 $GOPATH/bin
```

### 配置
編輯 `~/.serverpulse/config.yaml` 或使用 `--config` 指定設定檔路徑。

### Docker部署（規劃中）
```bash
docker run -d --name serverpulse \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v ~/.serverpulse:/root/.serverpulse \
  gitstq/serverpulse:latest serve
```

---

## 🤝 貢獻指南

詳見 [CONTRIBUTING.md](CONTRIBUTING.md)

---

## 📄 開源協議

本專案基於 [MIT License](LICENSE) 開源。

---

<div align="center">
  <sub>Built with ❤️ by <a href="https://github.com/gitstq">gitstq</a></sub>
</div>

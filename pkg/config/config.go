// Package config - 配置管理 / Configuration Management
// 负责加载、解析和管理 ServerPulse 的所有配置项
// Handles loading, parsing and managing all ServerPulse configurations
package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config - 全局配置结构 / Global configuration structure
type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Collector CollectorConfig `yaml:"collector"`
	Probe     ProbeConfig     `yaml:"probe"`
	Anomaly   AnomalyConfig   `yaml:"anomaly"`
	Alert     AlertConfig     `yaml:"alert"`
	Notify    NotifyConfig    `yaml:"notify"`
	Storage   StorageConfig   `yaml:"storage"`
	Hub       HubConfig       `yaml:"hub"`
	Agent     AgentConfig     `yaml:"agent"`
	TUI       TUIConfig       `yaml:"tui"`
	Snapshot  SnapshotConfig  `yaml:"snapshot"`
}

// ServerConfig - 服务器标识配置 / Server identity configuration
type ServerConfig struct {
	Name string `yaml:"name"` // 服务器名称 / Server name
	ID   string `yaml:"id"`   // 唯一标识，留空自动生成 / Unique ID, auto-generated if empty
}

// CollectorConfig - 采集器配置 / Collector configuration
type CollectorConfig struct {
	Interval time.Duration `yaml:"interval"` // 采集间隔 / Collection interval
	CPU      bool          `yaml:"cpu"`       // CPU采集开关 / Enable CPU collection
	Memory   bool          `yaml:"memory"`    // 内存采集开关 / Enable memory collection
	Disk     bool          `yaml:"disk"`      // 磁盘采集开关 / Enable disk collection
	Network  bool          `yaml:"network"`   // 网络采集开关 / Enable network collection
	Docker   bool          `yaml:"docker"`    // Docker采集开关 / Enable Docker collection
}

// ProbeConfig - 探测配置 / Probe configuration
type ProbeConfig struct {
	Interval time.Duration `yaml:"interval"` // 探测间隔 / Probe interval
	Targets  []ProbeTarget `yaml:"targets"`  // 探测目标列表 / Probe target list
}

// ProbeTarget - 探测目标 / Probe target definition
type ProbeTarget struct {
	Name    string        `yaml:"name"`    // 目标名称 / Target name
	Type    string        `yaml:"type"`    // 探测类型: ping, tcp, http, ssh / Probe type
	Host    string        `yaml:"host"`    // 目标主机 / Target host
	Port    int           `yaml:"port"`    // 目标端口(TCP/SSH) / Target port
	URL     string        `yaml:"url"`     // 目标URL(HTTP) / Target URL
	Timeout time.Duration `yaml:"timeout"` // 超时时间 / Timeout
}

// AnomalyConfig - 异常检测配置 / Anomaly detection configuration
type AnomalyConfig struct {
	Enabled        bool    `yaml:"enabled"`         // 启用异常检测 / Enable anomaly detection
	WindowSize     int     `yaml:"window_size"`      // 滑动窗口大小 / Sliding window size
	ZScoreThreshold float64 `yaml:"zscore_threshold"` // Z-Score阈值 / Z-Score threshold
	Sensitivity    string  `yaml:"sensitivity"`     // 灵敏度 / Sensitivity level
}

// AlertConfig - 告警配置 / Alert configuration
type AlertConfig struct {
	Enabled         bool              `yaml:"enabled"`          // 启用告警 / Enable alerts
	AggregateWindow time.Duration     `yaml:"aggregate_window"` // 聚合窗口 / Aggregation window
	Cooldown        time.Duration     `yaml:"cooldown"`         // 冷却时间 / Cooldown period
	Thresholds      []ThresholdRule   `yaml:"thresholds"`       // 阈值规则 / Threshold rules
}

// ThresholdRule - 阈值告警规则 / Threshold alert rule
type ThresholdRule struct {
	Metric   string  `yaml:"metric"`   // 指标名称 / Metric name
	Warning  float64 `yaml:"warning"`  // 警告阈值 / Warning threshold
	Critical float64 `yaml:"critical"` // 严重阈值 / Critical threshold
}

// NotifyConfig - 通知配置 / Notification configuration
type NotifyConfig struct {
	Webhook WebhookConfig `yaml:"webhook"` // Webhook配置 / Webhook config
	Email   EmailConfig   `yaml:"email"`   // 邮件配置 / Email config
}

// WebhookConfig - Webhook通知配置 / Webhook notification configuration
type WebhookConfig struct {
	Enabled bool   `yaml:"enabled"` // 启用Webhook / Enable webhook
	URL     string `yaml:"url"`     // Webhook URL
	Secret  string `yaml:"secret"`  // 签名密钥 / Signing secret
}

// EmailConfig - 邮件通知配置 / Email notification configuration
type EmailConfig struct {
	Enabled  bool     `yaml:"enabled"`   // 启用邮件通知 / Enable email
	SMTPHost string   `yaml:"smtp_host"` // SMTP服务器 / SMTP server
	SMTPPort int      `yaml:"smtp_port"` // SMTP端口 / SMTP port
	Username string   `yaml:"username"`  // SMTP用户名 / SMTP username
	Password string   `yaml:"password"`  // SMTP密码 / SMTP password
	From     string   `yaml:"from"`      // 发件人 / Sender address
	To       []string `yaml:"to"`        // 收件人列表 / Recipient list
}

// StorageConfig - 存储配置 / Storage configuration
type StorageConfig struct {
	Path          string `yaml:"path"`           // SQLite数据库路径 / Database path
	RetentionDays int    `yaml:"retention_days"`  // 数据保留天数 / Data retention days
}

// HubConfig - Hub服务配置 / Hub server configuration
type HubConfig struct {
	Enabled bool   `yaml:"enabled"` // 是否作为Hub运行 / Run as Hub
	Bind    string `yaml:"bind"`    // 监听地址 / Listen address
	APIKey  string `yaml:"api_key"` // API密钥 / API key
}

// AgentConfig - Agent配置 / Agent configuration
type AgentConfig struct {
	Enabled        bool          `yaml:"enabled"`         // 是否作为Agent运行 / Run as Agent
	HubURL         string        `yaml:"hub_url"`         // Hub地址 / Hub URL
	APIKey         string        `yaml:"api_key"`         // Hub API密钥 / Hub API key
	ReportInterval time.Duration `yaml:"report_interval"` // 上报间隔 / Report interval
}

// TUIConfig - TUI仪表板配置 / TUI dashboard configuration
type TUIConfig struct {
	RefreshRate time.Duration `yaml:"refresh_rate"` // 刷新频率 / Refresh rate
	Theme       string        `yaml:"theme"`        // 主题名称 / Theme name
}

// SnapshotConfig - 快照配置 / Snapshot configuration
type SnapshotConfig struct {
	Auto      bool          `yaml:"auto"`       // 自动快照 / Auto snapshot
	Interval  time.Duration `yaml:"interval"`   // 自动快照间隔 / Auto snapshot interval
	MaxCount  int           `yaml:"max_count"`  // 最大快照数量 / Max snapshot count
}

// Default - 返回默认配置 / Return default configuration
func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Name: "localhost",
			ID:   "",
		},
		Collector: CollectorConfig{
			Interval: 10 * time.Second,
			CPU:      true,
			Memory:   true,
			Disk:     true,
			Network:  true,
			Docker:   false,
		},
		Probe: ProbeConfig{
			Interval: 30 * time.Second,
			Targets: []ProbeTarget{
				{Name: "google-dns", Type: "ping", Host: "8.8.8.8"},
				{Name: "local-ssh", Type: "tcp", Host: "127.0.0.1", Port: 22},
				{Name: "local-http", Type: "http", URL: "http://localhost:8080/health", Timeout: 5 * time.Second},
			},
		},
		Anomaly: AnomalyConfig{
			Enabled:         true,
			WindowSize:      20,
			ZScoreThreshold: 2.5,
			Sensitivity:     "medium",
		},
		Alert: AlertConfig{
			Enabled:         true,
			AggregateWindow: 60 * time.Second,
			Cooldown:        300 * time.Second,
			Thresholds: []ThresholdRule{
				{Metric: "cpu_usage", Warning: 80, Critical: 95},
				{Metric: "memory_usage", Warning: 85, Critical: 95},
				{Metric: "disk_usage", Warning: 85, Critical: 95},
			},
		},
		Notify: NotifyConfig{
			Webhook: WebhookConfig{Enabled: false},
			Email:   EmailConfig{SMTPPort: 587},
		},
		Storage: StorageConfig{
			Path:          "./data/serverpulse.db",
			RetentionDays: 30,
		},
		Hub: HubConfig{
			Enabled: false,
			Bind:    ":9090",
		},
		Agent: AgentConfig{
			Enabled:        false,
			HubURL:         "http://localhost:9090",
			ReportInterval: 10 * time.Second,
		},
		TUI: TUIConfig{
			RefreshRate: 1 * time.Second,
			Theme:       "default",
		},
		Snapshot: SnapshotConfig{
			Auto:     false,
			Interval: 3600 * time.Second,
			MaxCount: 100,
		},
	}
}

// LoadFromFile - 从YAML文件加载配置 / Load configuration from YAML file
func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadFromBytes(data)
}

// LoadFromBytes - 从字节加载配置 / Load configuration from bytes
func LoadFromBytes(data []byte) (*Config, error) {
	cfg := Default()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// SaveToFile - 将配置保存到YAML文件 / Save configuration to YAML file
func (c *Config) SaveToFile(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

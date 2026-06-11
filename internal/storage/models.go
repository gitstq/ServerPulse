// Package storage - 数据存储模型 / Data Storage Models
// 定义所有持久化数据的结构体
// Defines all data structures for persistence
package storage

import "time"

// MetricData - 指标数据记录 / Metric data record
type MetricData struct {
	ID        int64     `json:"id" db:"id"`
	ServerID  string    `json:"server_id" db:"server_id"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	Metric    string    `json:"metric" db:"metric"`       // 指标名称 / Metric name (cpu_usage, memory_usage, ...)
	Value     float64   `json:"value" db:"value"`         // 指标值 / Metric value
	Unit      string    `json:"unit" db:"unit"`           // 单位 / Unit (%, MB, Mbps, ...)
	Labels    string    `json:"labels" db:"labels"`       // 附加标签JSON / Additional labels JSON
}

// ProbeResult - 探测结果记录 / Probe result record
type ProbeResult struct {
	ID        int64     `json:"id" db:"id"`
	ServerID  string    `json:"server_id" db:"server_id"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	Target    string    `json:"target" db:"target"`       // 目标名称 / Target name
	Type      string    `json:"type" db:"type"`           // 探测类型 / Probe type (ping, tcp, http, ssh)
	Status    string    `json:"status" db:"status"`       // 状态: ok, fail, timeout / Status
	Latency   float64   `json:"latency" db:"latency"`     // 延迟(ms) / Latency in milliseconds
	Message   string    `json:"message" db:"message"`     // 附加信息 / Additional message
}

// AlertRecord - 告警记录 / Alert record
type AlertRecord struct {
	ID        int64     `json:"id" db:"id"`
	ServerID  string    `json:"server_id" db:"server_id"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	Level     string    `json:"level" db:"level"`         // 级别: warning, critical / Severity level
	Type      string    `json:"type" db:"type"`           // 类型: threshold, anomaly / Alert type
	Metric    string    `json:"metric" db:"metric"`       // 关联指标 / Related metric
	Value     float64   `json:"value" db:"value"`         // 触发值 / Triggered value
	Threshold float64   `json:"threshold" db:"threshold"` // 阈值 / Threshold
	Message   string    `json:"message" db:"message"`     // 告警消息 / Alert message
	Acked     bool      `json:"acked" db:"acked"`         // 是否已确认 / Whether acknowledged
}

// ServerInfo - 服务器信息 / Server information
type ServerInfo struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	IP        string    `json:"ip" db:"ip"`
	OS        string    `json:"os" db:"os"`
	Arch      string    `json:"arch" db:"arch"`
	Status    string    `json:"status" db:"status"`       // online, offline / Status
	LastSeen  time.Time `json:"last_seen" db:"last_seen"` // 最后在线时间 / Last seen time
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// SnapshotData - 快照数据 / Snapshot data
type SnapshotData struct {
	ID        int64     `json:"id" db:"id"`
	ServerID  string    `json:"server_id" db:"server_id"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	Name      string    `json:"name" db:"name"`       // 快照名称 / Snapshot name
	Data      string    `json:"data" db:"data"`       // 快照数据JSON / Snapshot data JSON
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

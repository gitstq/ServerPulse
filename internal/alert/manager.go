// Package alert - 告警管理器 / Alert Manager
// 管理阈值告警和异常告警的生成、存储和分发
// Manages generation, storage and distribution of threshold and anomaly alerts
package alert

import (
	"fmt"
	"sync"
	"time"

	"github.com/gitstq/ServerPulse/pkg/config"
)

// Alert - 告警信息 / Alert information
type Alert struct {
	ID        string    `json:"id"`         // 告警唯一ID / Unique alert ID
	ServerID  string    `json:"server_id"`  // 服务器ID / Server ID
	Timestamp time.Time `json:"timestamp"`  // 告警时间 / Alert timestamp
	Level     string    `json:"level"`      // 级别: warning, critical / Level
	Type      string    `json:"type"`      // 类型: threshold, anomaly / Type
	Metric    string    `json:"metric"`     // 关联指标 / Related metric
	Value     float64   `json:"value"`      // 触发值 / Triggered value
	Threshold float64   `json:"threshold"` // 阈值 / Threshold
	Message   string    `json:"message"`    // 告警消息 / Alert message
	Acked     bool      `json:"acked"`      // 是否已确认 / Whether acknowledged
}

// AlertCallback - 告警回调函数类型 / Alert callback function type
type AlertCallback func(alert Alert)

// Manager - 告警管理器 / Alert manager
type Manager struct {
	mu           sync.RWMutex
	alerts       []Alert              // 活跃告警列表 / Active alert list
	thresholds   map[string]config.ThresholdRule // 阈值规则映射 / Threshold rule map
	cooldowns    map[string]time.Time // 告警冷却时间 / Alert cooldown times
	cooldownDur  time.Duration        // 冷却持续时间 / Cooldown duration
	callbacks    []AlertCallback     // 告警回调 / Alert callbacks
	aggregator   *Aggregator         // 告警聚合器 / Alert aggregator
}

// NewManager - 创建告警管理器 / Create alert manager
func NewManager(cfg config.AlertConfig) *Manager {
	m := &Manager{
		alerts:      make([]Alert, 0),
		thresholds:  make(map[string]config.ThresholdRule),
		cooldowns:   make(map[string]time.Time),
		cooldownDur: cfg.Cooldown,
		callbacks:   make([]AlertCallback, 0),
		aggregator:  NewAggregator(cfg.AggregateWindow),
	}

	// 加载阈值规则 / Load threshold rules
	for _, rule := range cfg.Thresholds {
		m.thresholds[rule.Metric] = rule
	}

	return m
}

// AddCallback 添加告警回调 / Add alert callback
func (m *Manager) AddCallback(cb AlertCallback) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callbacks = append(m.callbacks, cb)
}

// CheckThreshold - 检查阈值告警 / Check threshold alerts
// 对比当前值与阈值规则，生成告警
// Compare current value against threshold rules, generate alerts
func (m *Manager) CheckThreshold(serverID, metric string, value float64) *Alert {
	m.mu.Lock()
	defer m.mu.Unlock()

	rule, ok := m.thresholds[metric]
	if !ok {
		return nil
	}

	// 检查冷却时间 / Check cooldown
	cooldownKey := fmt.Sprintf("%s:%s", serverID, metric)
	if lastTrigger, ok := m.cooldowns[cooldownKey]; ok {
		if time.Since(lastTrigger) < m.cooldownDur {
			return nil
		}
	}

	var alertLevel string
	var threshold float64

	if value >= rule.Critical {
		alertLevel = "critical"
		threshold = rule.Critical
	} else if value >= rule.Warning {
		alertLevel = "warning"
		threshold = rule.Warning
	} else {
		return nil // 未触发阈值 / Below threshold
	}

	alert := Alert{
		ID:        fmt.Sprintf("%s-%s-%d", serverID, metric, time.Now().UnixNano()),
		ServerID:  serverID,
		Timestamp: time.Now(),
		Level:     alertLevel,
		Type:      "threshold",
		Metric:    metric,
		Value:     value,
		Threshold: threshold,
		Message:   fmt.Sprintf("[%s] 阈值告警: %s = %.2f%% (阈值: %.1f%%)", alertLevel, metric, value, threshold),
		Acked:     false,
	}

	// 更新冷却时间 / Update cooldown
	m.cooldowns[cooldownKey] = time.Now()

	// 添加到活跃告警 / Add to active alerts
	m.alerts = append(m.alerts, alert)

	// 通知回调 / Notify callbacks
	for _, cb := range m.callbacks {
		go cb(alert)
	}

	return &alert
}

// AddAnomalyAlert - 添加异常告警 / Add anomaly alert
func (m *Manager) AddAnomalyAlert(serverID, metric string, value float64, zScore float64, severity string) *Alert {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查冷却时间 / Check cooldown
	cooldownKey := fmt.Sprintf("anomaly:%s:%s", serverID, metric)
	if lastTrigger, ok := m.cooldowns[cooldownKey]; ok {
		if time.Since(lastTrigger) < m.cooldownDur {
			return nil
		}
	}

	alertLevel := "warning"
	if severity == "high" {
		alertLevel = "critical"
	}

	alert := Alert{
		ID:        fmt.Sprintf("anomaly-%s-%s-%d", serverID, metric, time.Now().UnixNano()),
		ServerID:  serverID,
		Timestamp: time.Now(),
		Level:     alertLevel,
		Type:      "anomaly",
		Metric:    metric,
		Value:     value,
		Threshold: zScore,
		Message:   fmt.Sprintf("[%s] 异常告警: %s = %.2f, Z-Score = %.2f", alertLevel, metric, value, zScore),
		Acked:     false,
	}

	// 更新冷却时间 / Update cooldown
	m.cooldowns[cooldownKey] = time.Now()

	// 添加到活跃告警 / Add to active alerts
	m.alerts = append(m.alerts, alert)

	// 通知回调 / Notify callbacks
	for _, cb := range m.callbacks {
		go cb(alert)
	}

	return &alert
}

// GetAlerts - 获取活跃告警列表 / Get active alert list
func (m *Manager) GetAlerts() []Alert {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]Alert, len(m.alerts))
	copy(result, m.alerts)
	return result
}

// GetActiveAlerts - 获取未确认的活跃告警 / Get unacknowledged active alerts
func (m *Manager) GetActiveAlerts() []Alert {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var active []Alert
	for _, a := range m.alerts {
		if !a.Acked {
			active = append(active, a)
		}
	}
	return active
}

// AckAlert - 确认告警 / Acknowledge alert
func (m *Manager) AckAlert(alertID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range m.alerts {
		if m.alerts[i].ID == alertID {
			m.alerts[i].Acked = true
			return true
		}
	}
	return false
}

// ClearAlerts - 清除已确认的告警 / Clear acknowledged alerts
func (m *Manager) ClearAlerts() {
	m.mu.Lock()
	defer m.mu.Unlock()

	var remaining []Alert
	for _, a := range m.alerts {
		if !a.Acked {
			remaining = append(remaining, a)
		}
	}
	m.alerts = remaining
}

// AlertCount - 获取告警数量 / Get alert count
func (m *Manager) AlertCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.alerts)
}

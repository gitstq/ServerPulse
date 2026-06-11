// Package alert - 告警聚合器 / Alert Aggregator
// 将短时间内的多个告警聚合为一条通知，避免告警风暴
// Aggregates multiple alerts within a short time into one notification to prevent alert storms
package alert

import (
	"sync"
	"time"
)

// AggregatedAlert - 聚合后的告警 / Aggregated alert
type AggregatedAlert struct {
	ID        string    `json:"id"`         // 聚合ID / Aggregation ID
	Timestamp time.Time `json:"timestamp"`  // 首次告警时间 / First alert time
	Count     int       `json:"count"`      // 包含的告警数量 / Number of alerts included
	Level     string    `json:"level"`      // 最高告警级别 / Highest alert level
	ServerID  string    `json:"server_id"`  // 服务器ID / Server ID
	Alerts    []Alert   `json:"alerts"`     // 原始告警列表 / Original alert list
	Summary   string    `json:"summary"`    // 聚合摘要 / Aggregation summary
}

// Aggregator - 告警聚合器 / Alert aggregator
type Aggregator struct {
	mu           sync.Mutex
	window       time.Duration // 聚合窗口 / Aggregation window
	pendingAlerts []Alert      // 待聚合的告警 / Pending alerts
	lastFlush    time.Time     // 上次刷新时间 / Last flush time
}

// NewAggregator - 创建告警聚合器 / Create alert aggregator
func NewAggregator(window time.Duration) *Aggregator {
	if window <= 0 {
		window = 60 * time.Second
	}
	return &Aggregator{
		window:        window,
		pendingAlerts: make([]Alert, 0),
		lastFlush:     time.Now(),
	}
}

// Add - 添加告警到聚合器 / Add alert to aggregator
func (a *Aggregator) Add(alert Alert) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.pendingAlerts = append(a.pendingAlerts, alert)
}

// Flush - 刷新聚合器，返回聚合结果 / Flush aggregator, return aggregated results
// 当聚合窗口到期或待聚合告警数量达到阈值时触发
// Triggered when aggregation window expires or pending alert count reaches threshold
func (a *Aggregator) Flush() *AggregatedAlert {
	a.mu.Lock()
	defer a.mu.Unlock()

	if len(a.pendingAlerts) == 0 {
		return nil
	}

	// 检查是否需要刷新 / Check if flush is needed
	if time.Since(a.lastFlush) < a.window && len(a.pendingAlerts) < 10 {
		return nil
	}

	// 创建聚合告警 / Create aggregated alert
	agg := &AggregatedAlert{
		ID:        time.Now().Format("20060102-150405"),
		Timestamp: a.pendingAlerts[0].Timestamp,
		Count:     len(a.pendingAlerts),
		Alerts:    a.pendingAlerts,
	}

	// 确定最高告警级别 / Determine highest alert level
	agg.Level = "warning"
	agg.ServerID = a.pendingAlerts[0].ServerID
	for _, al := range a.pendingAlerts {
		if al.Level == "critical" {
			agg.Level = "critical"
		}
	}

	// 生成摘要 / Generate summary
	agg.Summary = a.generateSummary()

	// 清空待聚合列表 / Clear pending alerts
	a.pendingAlerts = make([]Alert, 0)
	a.lastFlush = time.Now()

	return agg
}

// ForceFlush - 强制刷新聚合器 / Force flush aggregator
func (a *Aggregator) ForceFlush() *AggregatedAlert {
	a.mu.Lock()
	defer a.mu.Unlock()

	if len(a.pendingAlerts) == 0 {
		return nil
	}

	agg := &AggregatedAlert{
		ID:        time.Now().Format("20060102-150405"),
		Timestamp: a.pendingAlerts[0].Timestamp,
		Count:     len(a.pendingAlerts),
		Alerts:    a.pendingAlerts,
	}

	agg.Level = "warning"
	agg.ServerID = a.pendingAlerts[0].ServerID
	for _, al := range a.pendingAlerts {
		if al.Level == "critical" {
			agg.Level = "critical"
		}
	}

	agg.Summary = a.generateSummary()
	a.pendingAlerts = make([]Alert, 0)
	a.lastFlush = time.Now()

	return agg
}

// generateSummary - 生成聚合摘要 / Generate aggregation summary
func (a *Aggregator) generateSummary() string {
	if len(a.pendingAlerts) == 0 {
		return ""
	}

	// 统计各级别告警数量 / Count alerts by level
	warningCount := 0
	criticalCount := 0
	metrics := make(map[string]bool)

	for _, al := range a.pendingAlerts {
		switch al.Level {
		case "critical":
			criticalCount++
		case "warning":
			warningCount++
		}
		metrics[al.Metric] = true
	}

	summary := "告警聚合: "
	if criticalCount > 0 {
		summary += "严重:" + itoa(criticalCount) + " "
	}
	if warningCount > 0 {
		summary += "警告:" + itoa(warningCount) + " "
	}
	summary += "涉及指标: " + itoa(len(metrics)) + "个"

	return summary
}

// itoa - 整数转字符串 / Integer to string
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	s := ""
	n := i
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

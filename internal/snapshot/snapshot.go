// Package snapshot - 快照对比功能 / Snapshot Comparison Feature
// 创建系统状态快照并支持历史状态对比
// Creates system state snapshots and supports historical state comparison
package snapshot

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gitstq/ServerPulse/internal/collector"
	"github.com/gitstq/ServerPulse/internal/storage"
)

// Snapshot - 系统状态快照 / System state snapshot
type Snapshot struct {
	ID        string               `json:"id"`
	ServerID  string               `json:"server_id"`
	Name      string               `json:"name"`
	Timestamp time.Time            `json:"timestamp"`
	Metrics   *collector.SystemMetrics `json:"metrics"`
}

// Comparison - 快照对比结果 / Snapshot comparison result
type Comparison struct {
	SnapshotA    *Snapshot  `json:"snapshot_a"`     // 快照A / Snapshot A
	SnapshotB    *Snapshot  `json:"snapshot_b"`     // 快照B / Snapshot B
	Differences  []DiffItem `json:"differences"`    // 差异列表 / Difference list
	Summary      string     `json:"summary"`        // 对比摘要 / Comparison summary
}

// DiffItem - 单个指标差异 / Single metric difference
type DiffItem struct {
	Metric    string  `json:"metric"`     // 指标名称 / Metric name
	ValueA    float64 `json:"value_a"`    // 快照A的值 / Value in snapshot A
	ValueB    float64 `json:"value_b"`    // 快照B的值 / Value in snapshot B
	Change    float64 `json:"change"`     // 变化量 / Change amount
	ChangePct float64 `json:"change_pct"` // 变化百分比 / Change percentage
	Direction string  `json:"direction"`  // 变化方向: up, down, same / Change direction
}

// Manager - 快照管理器 / Snapshot manager
type Manager struct {
	store    *storage.SQLiteStorage
	maxCount int
}

// NewManager - 创建快照管理器 / Create snapshot manager
func NewManager(store *storage.SQLiteStorage, maxCount int) *Manager {
	if maxCount <= 0 {
		maxCount = 100
	}
	return &Manager{
		store:    store,
		maxCount: maxCount,
	}
}

// Create - 创建快照 / Create snapshot
func (m *Manager) Create(serverID, name string, metrics *collector.SystemMetrics) (*Snapshot, error) {
	snap := &Snapshot{
		ID:        fmt.Sprintf("%s-%d", serverID, time.Now().UnixNano()),
		ServerID:  serverID,
		Name:      name,
		Timestamp: metrics.Timestamp,
		Metrics:   metrics,
	}

	// 序列化快照数据 / Serialize snapshot data
	data, err := json.Marshal(metrics)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	// 存储到数据库 / Store to database
	if m.store != nil {
		snapData := &storage.SnapshotData{
			ServerID:  serverID,
			Timestamp: metrics.Timestamp,
			Name:      name,
			Data:      string(data),
			CreatedAt: time.Now(),
		}
		if err := m.store.InsertSnapshot(snapData); err != nil {
			return nil, fmt.Errorf("failed to store snapshot: %w", err)
		}
	}

	return snap, nil
}

// Compare - 对比两个快照 / Compare two snapshots
func (m *Manager) Compare(a, b *Snapshot) *Comparison {
	result := &Comparison{
		SnapshotA: a,
		SnapshotB: b,
	}

	if a.Metrics == nil || b.Metrics == nil {
		result.Summary = "无法对比：快照数据不完整"
		return result
	}

	// 对比CPU / Compare CPU
	result.Differences = append(result.Differences, DiffItem{
		Metric:    "cpu_usage",
		ValueA:    a.Metrics.CPU.UsagePercent,
		ValueB:    b.Metrics.CPU.UsagePercent,
		Change:    b.Metrics.CPU.UsagePercent - a.Metrics.CPU.UsagePercent,
		ChangePct: calcChangePct(a.Metrics.CPU.UsagePercent, b.Metrics.CPU.UsagePercent),
		Direction: getDirection(a.Metrics.CPU.UsagePercent, b.Metrics.CPU.UsagePercent),
	})

	// 对比内存 / Compare Memory
	result.Differences = append(result.Differences, DiffItem{
		Metric:    "memory_usage",
		ValueA:    a.Metrics.Memory.UsagePercent,
		ValueB:    b.Metrics.Memory.UsagePercent,
		Change:    b.Metrics.Memory.UsagePercent - a.Metrics.Memory.UsagePercent,
		ChangePct: calcChangePct(a.Metrics.Memory.UsagePercent, b.Metrics.Memory.UsagePercent),
		Direction: getDirection(a.Metrics.Memory.UsagePercent, b.Metrics.Memory.UsagePercent),
	})

	// 对比磁盘 / Compare Disk
	if len(a.Metrics.Disk) > 0 && len(b.Metrics.Disk) > 0 {
		result.Differences = append(result.Differences, DiffItem{
			Metric:    "disk_usage",
			ValueA:    a.Metrics.Disk[0].UsagePercent,
			ValueB:    b.Metrics.Disk[0].UsagePercent,
			Change:    b.Metrics.Disk[0].UsagePercent - a.Metrics.Disk[0].UsagePercent,
			ChangePct: calcChangePct(a.Metrics.Disk[0].UsagePercent, b.Metrics.Disk[0].UsagePercent),
			Direction: getDirection(a.Metrics.Disk[0].UsagePercent, b.Metrics.Disk[0].UsagePercent),
		})
	}

	// 生成摘要 / Generate summary
	result.Summary = generateComparisonSummary(result.Differences)

	return result
}

// List - 列出所有快照 / List all snapshots
func (m *Manager) List(serverID string) ([]storage.SnapshotData, error) {
	if m.store == nil {
		return nil, fmt.Errorf("storage not initialized")
	}
	return m.store.GetSnapshots(serverID, m.maxCount)
}

// calcChangePct - 计算变化百分比 / Calculate change percentage
func calcChangePct(oldVal, newVal float64) float64 {
	if oldVal == 0 {
		return 0
	}
	return (newVal - oldVal) / oldVal * 100
}

// getDirection - 获取变化方向 / Get change direction
func getDirection(oldVal, newVal float64) string {
	if newVal > oldVal {
		return "up"
	} else if newVal < oldVal {
		return "down"
	}
	return "same"
}

// generateComparisonSummary - 生成对比摘要 / Generate comparison summary
func generateComparisonSummary(diffs []DiffItem) string {
	ups := 0
	downs := 0
	for _, d := range diffs {
		switch d.Direction {
		case "up":
			ups++
		case "down":
			downs++
		}
	}
	return fmt.Sprintf("对比完成: %d项指标, %d项上升, %d项下降", len(diffs), ups, downs)
}

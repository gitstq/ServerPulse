// Package collector - 磁盘采集器 / Disk Collector
package collector

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
)

// DiskMetrics - 单个磁盘/分区的指标 / Single disk/partition metrics
type DiskMetrics struct {
	Device      string  `json:"device"`       // 设备名 / Device name
	Mountpoint  string  `json:"mountpoint"`   // 挂载点 / Mount point
	Fstype      string  `json:"fstype"`       // 文件系统类型 / Filesystem type
	Total       uint64  `json:"total"`        // 总容量(bytes) / Total capacity
	Used        uint64  `json:"used"`         // 已使用(bytes) / Used space
	Free        uint64  `json:"free"`         // 可用(bytes) / Free space
	UsagePercent float64 `json:"usage_percent"` // 使用率(%) / Usage percentage
}

// DiskCollector - 磁盘采集器 / Disk collector
type DiskCollector struct{}

// Name 返回采集器名称 / Returns collector name
func (c *DiskCollector) Name() string {
	return "disk"
}

// Collect 采集所有磁盘分区指标 / Collect all disk partition metrics
func (c *DiskCollector) Collect() ([]DiskMetrics, error) {
	var metrics []DiskMetrics

	// 获取所有分区 / Get all partitions
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, fmt.Errorf("failed to get disk partitions: %w", err)
	}

	for _, part := range partitions {
		usage, err := disk.Usage(part.Mountpoint)
		if err != nil {
			// 跳过无法读取的分区 / Skip unreadable partitions
			continue
		}

		metrics = append(metrics, DiskMetrics{
			Device:       part.Device,
			Mountpoint:   part.Mountpoint,
			Fstype:       part.Fstype,
			Total:        usage.Total,
			Used:         usage.Used,
			Free:         usage.Free,
			UsagePercent: usage.UsedPercent,
		})
	}

	return metrics, nil
}

// CollectAsMetrics - 实现Collector接口 / Implement Collector interface
func (c *DiskCollector) CollectAsMetrics() ([]Metric, error) {
	disks, err := c.Collect()
	if err != nil {
		return nil, err
	}
	var metrics []Metric
	ts := time.Now()
	for _, d := range disks {
		labels := fmt.Sprintf(`{"device":"%s","mount":"%s"}`, d.Device, d.Mountpoint)
		metrics = append(metrics,
			Metric{Name: "disk_usage", Value: d.UsagePercent, Unit: "%", Timestamp: ts, Labels: labels},
			Metric{Name: "disk_total", Value: float64(d.Total), Unit: "bytes", Timestamp: ts, Labels: labels},
			Metric{Name: "disk_used", Value: float64(d.Used), Unit: "bytes", Timestamp: ts, Labels: labels},
			Metric{Name: "disk_free", Value: float64(d.Free), Unit: "bytes", Timestamp: ts, Labels: labels},
		)
	}
	return metrics, nil
}

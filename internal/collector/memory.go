// Package collector - 内存采集器 / Memory Collector
package collector

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/mem"
)

// MemoryMetrics - 内存指标数据 / Memory metrics data
type MemoryMetrics struct {
	Total        uint64  `json:"total"`         // 总内存(bytes) / Total memory
	Used         uint64  `json:"used"`          // 已使用(bytes) / Used memory
	Available    uint64  `json:"available"`     // 可用(bytes) / Available memory
	UsagePercent float64 `json:"usage_percent"` // 使用率(%) / Usage percentage
	SwapTotal    uint64  `json:"swap_total"`    // Swap总量(bytes) / Swap total
	SwapUsed     uint64  `json:"swap_used"`     // Swap已使用(bytes) / Swap used
}

// MemoryCollector - 内存采集器 / Memory collector
type MemoryCollector struct{}

// Name 返回采集器名称 / Returns collector name
func (c *MemoryCollector) Name() string {
	return "memory"
}

// Collect 采集内存指标 / Collect memory metrics
func (c *MemoryCollector) Collect() (MemoryMetrics, error) {
	metrics := MemoryMetrics{}

	// 获取虚拟内存信息 / Get virtual memory info
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return metrics, fmt.Errorf("failed to get virtual memory: %w", err)
	}

	metrics.Total = vmem.Total
	metrics.Used = vmem.Used
	metrics.Available = vmem.Available
	metrics.UsagePercent = vmem.UsedPercent

	// 获取Swap信息 / Get swap info
	swap, err := mem.SwapMemory()
	if err != nil {
		return metrics, fmt.Errorf("failed to get swap memory: %w", err)
	}
	metrics.SwapTotal = swap.Total
	metrics.SwapUsed = swap.Used

	return metrics, nil
}

// CollectAsMetrics - 实现Collector接口 / Implement Collector interface
func (c *MemoryCollector) CollectAsMetrics() ([]Metric, error) {
	mem, err := c.Collect()
	if err != nil {
		return nil, err
	}
	return []Metric{
		{Name: "memory_usage", Value: mem.UsagePercent, Unit: "%", Timestamp: time.Now()},
		{Name: "memory_total", Value: float64(mem.Total), Unit: "bytes", Timestamp: time.Now()},
		{Name: "memory_used", Value: float64(mem.Used), Unit: "bytes", Timestamp: time.Now()},
		{Name: "memory_available", Value: float64(mem.Available), Unit: "bytes", Timestamp: time.Now()},
		{Name: "swap_usage", Value: float64(mem.SwapUsed) / float64(mem.SwapTotal+1) * 100, Unit: "%", Timestamp: time.Now()},
	}, nil
}

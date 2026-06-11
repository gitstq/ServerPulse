// Package collector - CPU采集器 / CPU Collector
package collector

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

// CPUMetrics - CPU指标数据 / CPU metrics data
type CPUMetrics struct {
	UsagePercent float64 `json:"usage_percent"` // CPU使用率(%) / CPU usage percentage
	Cores        int     `json:"cores"`         // 逻辑核心数 / Logical core count
	ModelName    string  `json:"model_name"`    // CPU型号 / CPU model name
}

// CPUCollector - CPU采集器 / CPU collector
type CPUCollector struct {
	prevIdle     float64
	prevTotal    float64
	lastUsage    float64
}

// Name 返回采集器名称 / Returns collector name
func (c *CPUCollector) Name() string {
	return "cpu"
}

// Collect 采集CPU指标 / Collect CPU metrics
func (c *CPUCollector) Collect() (CPUMetrics, error) {
	metrics := CPUMetrics{}

	// 获取CPU使用率（取1秒内的平均值）/ Get CPU usage (average over 1 second)
	percents, err := cpu.Percent(1*time.Second, false)
	if err != nil {
		return metrics, fmt.Errorf("failed to get cpu percent: %w", err)
	}
	if len(percents) > 0 {
		metrics.UsagePercent = percents[0]
	}
	c.lastUsage = metrics.UsagePercent

	// 获取CPU信息 / Get CPU info
	info, err := cpu.Info()
	if err != nil {
		return metrics, fmt.Errorf("failed to get cpu info: %w", err)
	}
	if len(info) > 0 {
		metrics.ModelName = info[0].ModelName
		metrics.Cores = int(info[0].Cores)
	}

	// 获取逻辑核心数 / Get logical core count
	counts, err := cpu.Counts(true)
	if err == nil {
		metrics.Cores = counts
	}

	return metrics, nil
}

// CollectAsMetrics - 实现Collector接口 / Implement Collector interface
func (c *CPUCollector) CollectAsMetrics() ([]Metric, error) {
	cpu, err := c.Collect()
	if err != nil {
		return nil, err
	}
	return []Metric{
		{Name: "cpu_usage", Value: cpu.UsagePercent, Unit: "%", Timestamp: time.Now()},
		{Name: "cpu_cores", Value: float64(cpu.Cores), Unit: "cores", Timestamp: time.Now()},
	}, nil
}

// Package collector - 网络采集器 / Network Collector
package collector

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/net"
)

// NetworkMetric - 单个网络接口指标 / Single network interface metrics
type NetworkMetric struct {
	Name      string  `json:"name"`       // 接口名称 / Interface name
	BytesSent uint64  `json:"bytes_sent"` // 发送字节数 / Bytes sent
	BytesRecv uint64  `json:"bytes_recv"` // 接收字节数 / Bytes received
	PacketsSent uint64 `json:"packets_sent"` // 发送包数 / Packets sent
	PacketsRecv uint64 `json:"packets_recv"` // 接收包数 / Packets received
	Errin     uint64  `json:"errin"`      // 接收错误 / Receive errors
	Errout    uint64  `json:"errout"`     // 发送错误 / Send errors
	Dropin    uint64  `json:"dropin"`     // 接收丢弃 / Receive drops
	Dropout   uint64  `json:"dropout"`    // 发送丢弃 / Send drops
}

// NetworkCollector - 网络采集器 / Network collector
type NetworkCollector struct {
	prevStats map[string]NetworkMetric
}

// Name 返回采集器名称 / Returns collector name
func (c *NetworkCollector) Name() string {
	return "network"
}

// Collect 采集所有网络接口指标 / Collect all network interface metrics
func (c *NetworkCollector) Collect() ([]NetworkMetric, error) {
	var metrics []NetworkMetric

	// 获取所有网络接口IO统计 / Get all network interface IO counters
	counters, err := net.IOCounters(true) // pernic=true 获取每个接口 / per interface
	if err != nil {
		return nil, fmt.Errorf("failed to get network counters: %w", err)
	}

	for _, counter := range counters {
		metrics = append(metrics, NetworkMetric{
			Name:        counter.Name,
			BytesSent:   counter.BytesSent,
			BytesRecv:   counter.BytesRecv,
			PacketsSent: counter.PacketsSent,
			PacketsRecv: counter.PacketsRecv,
			Errin:       counter.Errin,
			Errout:      counter.Errout,
			Dropin:       counter.Dropin,
			Dropout:      counter.Dropout,
		})
	}

	// 保存当前值用于下次计算速率 / Save current values for rate calculation
	c.prevStats = make(map[string]NetworkMetric)
	for _, m := range metrics {
		c.prevStats[m.Name] = m
	}

	return metrics, nil
}

// CollectAsMetrics - 实现Collector接口 / Implement Collector interface
func (c *NetworkCollector) CollectAsMetrics() ([]Metric, error) {
	nets, err := c.Collect()
	if err != nil {
		return nil, err
	}
	var metrics []Metric
	ts := time.Now()
	for _, n := range nets {
		labels := fmt.Sprintf(`{"interface":"%s"}`, n.Name)
		metrics = append(metrics,
			Metric{Name: "network_bytes_sent", Value: float64(n.BytesSent), Unit: "bytes", Timestamp: ts, Labels: labels},
			Metric{Name: "network_bytes_recv", Value: float64(n.BytesRecv), Unit: "bytes", Timestamp: ts, Labels: labels},
			Metric{Name: "network_packets_sent", Value: float64(n.PacketsSent), Unit: "packets", Timestamp: ts, Labels: labels},
			Metric{Name: "network_packets_recv", Value: float64(n.PacketsRecv), Unit: "packets", Timestamp: ts, Labels: labels},
		)
	}
	return metrics, nil
}

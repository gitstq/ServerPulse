// Package collector - 系统指标采集器 / System Metrics Collectors
// 采集CPU、内存、磁盘、网络等系统指标
// Collects CPU, memory, disk, network and other system metrics
package collector

import "time"

// Metric - 单个指标数据点 / Single metric data point
type Metric struct {
	Name      string    `json:"name"`       // 指标名称 / Metric name
	Value     float64   `json:"value"`      // 指标值 / Metric value
	Unit      string    `json:"unit"`       // 单位 / Unit
	Timestamp time.Time `json:"timestamp"`  // 采集时间 / Collection time
	Labels    string    `json:"labels"`     // 附加标签 / Additional labels
}

// SystemMetrics - 系统综合指标 / Comprehensive system metrics
type SystemMetrics struct {
	Timestamp time.Time `json:"timestamp"` // 采集时间 / Collection time
	ServerID  string    `json:"server_id"` // 服务器ID / Server ID
	CPU       CPUMetrics    `json:"cpu"`       // CPU指标 / CPU metrics
	Memory    MemoryMetrics `json:"memory"`    // 内存指标 / Memory metrics
	Disk      []DiskMetrics `json:"disk"`      // 磁盘指标列表 / Disk metrics list
	Network   []NetworkMetric `json:"network"`  // 网络指标列表 / Network metrics list
	Docker    []DockerContainer `json:"docker"` // Docker容器列表 / Docker container list
}

// Collector - 采集器接口 / Collector interface
type Collector interface {
	// Name 返回采集器名称 / Returns collector name
	Name() string
	// Collect 执行一次采集 / Perform a single collection
	Collect() ([]Metric, error)
}

// SystemCollector - 系统综合采集器 / Comprehensive system collector
// 协调所有子采集器的工作 / Coordinates all sub-collectors
type SystemCollector struct {
	cpuCollector     *CPUCollector
	memoryCollector  *MemoryCollector
	diskCollector    *DiskCollector
	networkCollector *NetworkCollector
	dockerCollector  *DockerCollector
	serverID         string
	enableCPU        bool
	enableMemory     bool
	enableDisk       bool
	enableNetwork    bool
	enableDocker     bool
}

// NewSystemCollector - 创建系统综合采集器 / Create comprehensive system collector
func NewSystemCollector(serverID string, enableCPU, enableMemory, enableDisk, enableNetwork, enableDocker bool) *SystemCollector {
	return &SystemCollector{
		cpuCollector:     &CPUCollector{},
		memoryCollector:  &MemoryCollector{},
		diskCollector:    &DiskCollector{},
		networkCollector: &NetworkCollector{},
		dockerCollector:  &DockerCollector{},
		serverID:         serverID,
		enableCPU:        enableCPU,
		enableMemory:     enableMemory,
		enableDisk:       enableDisk,
		enableNetwork:    enableNetwork,
		enableDocker:     enableDocker,
	}
}

// CollectAll - 采集所有启用的指标 / Collect all enabled metrics
func (sc *SystemCollector) CollectAll() (*SystemMetrics, error) {
	result := &SystemMetrics{
		Timestamp: time.Now(),
		ServerID:  sc.serverID,
	}

	var errs []string

	// 采集CPU / Collect CPU
	if sc.enableCPU {
		cpu, err := sc.cpuCollector.Collect()
		if err != nil {
			errs = append(errs, "CPU: "+err.Error())
		} else {
			result.CPU = cpu
		}
	}

	// 采集内存 / Collect Memory
	if sc.enableMemory {
		mem, err := sc.memoryCollector.Collect()
		if err != nil {
			errs = append(errs, "Memory: "+err.Error())
		} else {
			result.Memory = mem
		}
	}

	// 采集磁盘 / Collect Disk
	if sc.enableDisk {
		disks, err := sc.diskCollector.Collect()
		if err != nil {
			errs = append(errs, "Disk: "+err.Error())
		} else {
			result.Disk = disks
		}
	}

	// 采集网络 / Collect Network
	if sc.enableNetwork {
		nets, err := sc.networkCollector.Collect()
		if err != nil {
			errs = append(errs, "Network: "+err.Error())
		} else {
			result.Network = nets
		}
	}

	// 采集Docker / Collect Docker
	if sc.enableDocker {
		containers, err := sc.dockerCollector.Collect()
		if err != nil {
			errs = append(errs, "Docker: "+err.Error())
		} else {
			result.Docker = containers
		}
	}

	if len(errs) > 0 {
		return result, nil // 返回已采集的数据，错误仅记录 / Return collected data, errors logged
	}

	return result, nil
}

// ToMetricList - 将系统指标转换为指标列表 / Convert system metrics to metric list
func (sc *SystemCollector) ToMetricList(sm *SystemMetrics) []Metric {
	var metrics []Metric
	ts := sm.Timestamp

	if sc.enableCPU {
		metrics = append(metrics,
			Metric{Name: "cpu_usage", Value: sm.CPU.UsagePercent, Unit: "%", Timestamp: ts},
			Metric{Name: "cpu_cores", Value: float64(sm.CPU.Cores), Unit: "cores", Timestamp: ts},
		)
	}

	if sc.enableMemory {
		metrics = append(metrics,
			Metric{Name: "memory_usage", Value: sm.Memory.UsagePercent, Unit: "%", Timestamp: ts},
			Metric{Name: "memory_total", Value: float64(sm.Memory.Total), Unit: "bytes", Timestamp: ts},
			Metric{Name: "memory_used", Value: float64(sm.Memory.Used), Unit: "bytes", Timestamp: ts},
			Metric{Name: "memory_available", Value: float64(sm.Memory.Available), Unit: "bytes", Timestamp: ts},
		)
	}

	if sc.enableDisk && len(sm.Disk) > 0 {
		// 使用第一个磁盘的总使用率 / Use first disk overall usage
		metrics = append(metrics,
			Metric{Name: "disk_usage", Value: sm.Disk[0].UsagePercent, Unit: "%", Timestamp: ts},
			Metric{Name: "disk_total", Value: float64(sm.Disk[0].Total), Unit: "bytes", Timestamp: ts},
			Metric{Name: "disk_used", Value: float64(sm.Disk[0].Used), Unit: "bytes", Timestamp: ts},
			Metric{Name: "disk_free", Value: float64(sm.Disk[0].Free), Unit: "bytes", Timestamp: ts},
		)
	}

	if sc.enableNetwork && len(sm.Network) > 0 {
		// 汇总所有接口的速率 / Sum rates across all interfaces
		var totalBytesSent, totalBytesRecv uint64
		for _, n := range sm.Network {
			totalBytesSent += n.BytesSent
			totalBytesRecv += n.BytesRecv
		}
		metrics = append(metrics,
			Metric{Name: "network_bytes_sent", Value: float64(totalBytesSent), Unit: "bytes", Timestamp: ts},
			Metric{Name: "network_bytes_recv", Value: float64(totalBytesRecv), Unit: "bytes", Timestamp: ts},
		)
	}

	return metrics
}

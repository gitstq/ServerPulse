// Package collector - Docker容器采集器 / Docker Container Collector
package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// DockerContainer - Docker容器信息 / Docker container information
type DockerContainer struct {
	ID        string  `json:"id"`         // 容器ID / Container ID
	Name      string  `json:"name"`       // 容器名称 / Container name
	Image     string  `json:"image"`      // 镜像名称 / Image name
	Status    string  `json:"status"`     // 状态 / Status (running, exited, ...)
	CPUPercent float64 `json:"cpu_percent"` // CPU使用率(%) / CPU usage percentage
	MemoryMB  float64 `json:"memory_mb"`  // 内存使用(MB) / Memory usage in MB
	NetworkIn float64 `json:"network_in"`  // 网络入流量(bytes) / Network input bytes
	NetworkOut float64 `json:"network_out"` // 网络出流量(bytes) / Network output bytes
}

// DockerCollector - Docker采集器 / Docker collector
type DockerCollector struct {
	httpClient *http.Client
	socketPath string
}

// NewDockerCollector - 创建Docker采集器 / Create Docker collector
func NewDockerCollector() *DockerCollector {
	return &DockerCollector{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				// 使用Unix socket连接Docker / Use Unix socket to connect to Docker
			},
		},
		socketPath: "/var/run/docker.sock",
	}
}

// Name 返回采集器名称 / Returns collector name
func (c *DockerCollector) Name() string {
	return "docker"
}

// Collect 采集Docker容器指标 / Collect Docker container metrics
func (c *DockerCollector) Collect() ([]DockerContainer, error) {
	// 检查Docker socket是否存在 / Check if Docker socket exists
	if _, err := os.Stat(c.socketPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("docker socket not found at %s", c.socketPath)
	}

	// 使用Docker API获取容器列表 / Use Docker API to get container list
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost/containers/json?all=true", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker request: %w", err)
	}

	// 使用Unix socket传输 / Use Unix socket transport
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to docker daemon: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("docker api returned status: %d", resp.StatusCode)
	}

	var rawContainers []json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&rawContainers); err != nil {
		return nil, fmt.Errorf("failed to decode docker response: %w", err)
	}

	containers := make([]DockerContainer, 0, len(rawContainers))
	for _, raw := range rawContainers {
		var basic struct {
			ID     string `json:"Id"`
			Names  []string `json:"Names"`
			Image  string `json:"Image"`
			State  string `json:"State"`
			Status string `json:"Status"`
		}
		if err := json.Unmarshal(raw, &basic); err != nil {
			continue
		}

		name := ""
		if len(basic.Names) > 0 {
			name = basic.Names[0]
		}

		containers = append(containers, DockerContainer{
			ID:     basic.ID[:12],
			Name:   name,
			Image:  basic.Image,
			Status: basic.Status,
		})
	}

	return containers, nil
}

// CollectAsMetrics - 实现Collector接口 / Implement Collector interface
func (c *DockerCollector) CollectAsMetrics() ([]Metric, error) {
	containers, err := c.Collect()
	if err != nil {
		return nil, err
	}
	var metrics []Metric
	ts := time.Now()
	for _, ctr := range containers {
		labels := fmt.Sprintf(`{"container_id":"%s","name":"%s","image":"%s"}`, ctr.ID, ctr.Name, ctr.Image)
		metrics = append(metrics,
			Metric{Name: "docker_container_cpu", Value: ctr.CPUPercent, Unit: "%", Timestamp: ts, Labels: labels},
			Metric{Name: "docker_container_memory", Value: ctr.MemoryMB, Unit: "MB", Timestamp: ts, Labels: labels},
		)
	}
	return metrics, nil
}

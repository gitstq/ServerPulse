// Package server - Agent客户端 / Agent Client
// 采集本地系统数据并定期上报到Hub
// Collects local system data and periodically reports to Hub
package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gitstq/ServerPulse/internal/collector"
	"github.com/gitstq/ServerPulse/pkg/config"
)

// Agent - Agent客户端 / Agent client
type Agent struct {
	cfg       config.AgentConfig
	serverCfg config.ServerConfig
	collector *collector.SystemCollector
	client    *http.Client
	stopCh    chan struct{}
}

// NewAgent - 创建Agent客户端 / Create Agent client
func NewAgent(agentCfg config.AgentConfig, serverCfg config.ServerConfig) *Agent {
	return &Agent{
		cfg: agentCfg,
		serverCfg: serverCfg,
		collector: collector.NewSystemCollector(
			serverCfg.ID,
			true, true, true, true, false,
		),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		stopCh: make(chan struct{}),
	}
}

// Start - 启动Agent / Start Agent
func (a *Agent) Start() error {
	// 先注册到Hub / Register with Hub first
	if err := a.register(); err != nil {
		return fmt.Errorf("failed to register with hub: %w", err)
	}

	log.Printf("[Agent] 已注册到Hub: %s", a.cfg.HubURL)

	// 启动数据上报协程 / Start data reporting goroutine
	go a.reportLoop()

	return nil
}

// Stop - 停止Agent / Stop Agent
func (a *Agent) Stop() {
	close(a.stopCh)
	log.Println("[Agent] 已停止")
}

// register - 向Hub注册 / Register with Hub
func (a *Agent) register() error {
	hostname := a.serverCfg.Name
	if hostname == "" {
		hostname = "unknown"
	}

	regData := AgentInfo{
		ID:      a.serverCfg.ID,
		Name:    hostname,
		Address: "",
		Status:  "online",
	}

	body, err := json.Marshal(regData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", a.cfg.HubURL+"/api/v1/register", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if a.cfg.APIKey != "" {
		req.Header.Set("X-API-Key", a.cfg.APIKey)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("registration failed with status: %d", resp.StatusCode)
	}

	return nil
}

// reportLoop - 数据上报循环 / Data reporting loop
func (a *Agent) reportLoop() {
	ticker := time.NewTicker(a.cfg.ReportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.report()
		case <-a.stopCh:
			return
		}
	}
}

// report - 执行一次数据上报 / Execute a single data report
func (a *Agent) report() {
	// 采集系统指标 / Collect system metrics
	metrics, err := a.collector.CollectAll()
	if err != nil {
		log.Printf("[Agent] 采集失败: %v", err)
		return
	}

	// 构建上报数据 / Build report data
	report := map[string]interface{}{
		"server_id": a.serverCfg.ID,
		"timestamp":  metrics.Timestamp.Format(time.RFC3339),
		"cpu":        metrics.CPU,
		"memory":     metrics.Memory,
		"disk":       metrics.Disk,
		"network":    metrics.Network,
	}

	body, err := json.Marshal(report)
	if err != nil {
		log.Printf("[Agent] 序列化失败: %v", err)
		return
	}

	// 发送数据到Hub / Send data to Hub
	req, err := http.NewRequest("POST", a.cfg.HubURL+"/api/v1/report", bytes.NewReader(body))
	if err != nil {
		log.Printf("[Agent] 创建请求失败: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	if a.cfg.APIKey != "" {
		req.Header.Set("X-API-Key", a.cfg.APIKey)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		log.Printf("[Agent] 上报失败: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		log.Printf("[Agent] 上报返回错误状态: %d", resp.StatusCode)
	}
}

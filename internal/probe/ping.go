// Package probe - 多协议探测 / Multi-Protocol Probing
// 支持Ping、TCP、HTTP、SSH等多种探测方式
// Supports Ping, TCP, HTTP, SSH and other probe methods
package probe

import (
	"fmt"
	"os/exec"
	"time"
)

// ProbeStatus - 探测状态 / Probe status
type ProbeStatus string

const (
	StatusOK      ProbeStatus = "ok"      // 正常 / Normal
	StatusFail    ProbeStatus = "fail"    // 失败 / Failed
	StatusTimeout ProbeStatus = "timeout" // 超时 / Timeout
)

// ProbeResult - 探测结果 / Probe result
type ProbeResult struct {
	Target    string      `json:"target"`    // 目标名称 / Target name
	Type      string      `json:"type"`      // 探测类型 / Probe type
	Host      string      `json:"host"`      // 目标主机 / Target host
	Status    ProbeStatus `json:"status"`    // 探测状态 / Probe status
	Latency   float64     `json:"latency"`   // 延迟(ms) / Latency in milliseconds
	Message   string      `json:"message"`   // 附加信息 / Additional message
	Timestamp time.Time   `json:"timestamp"` // 探测时间 / Probe timestamp
}

// Prober - 探测器接口 / Prober interface
type Prober interface {
	// Type 返回探测类型 / Returns probe type
	Type() string
	// Probe 执行一次探测 / Execute a single probe
	Probe(host string, port int, timeout time.Duration) ProbeResult
}

// PingProber - Ping探测器 / Ping prober
// 使用系统ping命令进行ICMP探测
// Uses system ping command for ICMP probing
type PingProber struct{}

// Type 返回探测类型 / Returns probe type
func (p *PingProber) Type() string {
	return "ping"
}

// Probe 执行Ping探测 / Execute Ping probe
// 调用系统ping命令并解析结果
// Calls system ping command and parses result
func (p *PingProber) Probe(host string, port int, timeout time.Duration) ProbeResult {
	result := ProbeResult{
		Target: host,
		Type:   "ping",
		Host:   host,
		Status: StatusFail,
	}

	if host == "" {
		result.Message = "empty host"
		return result
	}

	// 使用系统ping命令 / Use system ping command
	cmd := exec.Command("ping", "-c", "1", "-W", fmt.Sprintf("%.0f", timeout.Seconds()), host)
	start := time.Now()
	output, err := cmd.CombinedOutput()
	elapsed := time.Since(start)
	result.Latency = float64(elapsed.Milliseconds())
	result.Timestamp = time.Now()

	if err != nil {
		result.Status = StatusFail
		result.Message = fmt.Sprintf("ping failed: %v", err)
		return result
	}

	result.Status = StatusOK
	result.Message = fmt.Sprintf("ping ok, latency: %v, output: %s", elapsed, string(output))
	return result
}

// ProbeManager - 探测管理器 / Probe manager
// 管理所有探测目标的调度和执行 / Manages scheduling and execution of all probe targets
type ProbeManager struct {
	pingProber *PingProber
	tcpProber  *TCPProber
	httpProber *HTTPProber
	sshProber  *SSHProber
}

// NewProbeManager - 创建探测管理器 / Create probe manager
func NewProbeManager() *ProbeManager {
	return &ProbeManager{
		pingProber: &PingProber{},
		tcpProber:  &TCPProber{},
		httpProber: &HTTPProber{},
		sshProber:  &SSHProber{},
	}
}

// Probe - 根据类型执行探测 / Execute probe by type
func (pm *ProbeManager) Probe(targetType, host string, port int, url string, timeout time.Duration) ProbeResult {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	switch targetType {
	case "ping":
		return pm.pingProber.Probe(host, 0, timeout)
	case "tcp":
		return pm.tcpProber.Probe(host, port, timeout)
	case "http":
		return pm.httpProber.ProbeURL(url, timeout)
	case "ssh":
		return pm.sshProber.Probe(host, port, timeout)
	default:
		return ProbeResult{
			Target:  host,
			Type:    targetType,
			Host:    host,
			Status:  StatusFail,
			Message: "unknown probe type: " + targetType,
		}
	}
}

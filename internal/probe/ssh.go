// Package probe - SSH连通性检查 / SSH Connectivity Check Probe
package probe

import (
	"fmt"
	"net"
	"time"
)

// SSHProber - SSH连通性探测器 / SSH connectivity prober
// 通过尝试建立TCP连接到SSH默认端口(22)来检测SSH服务可用性
// Checks SSH service availability by attempting TCP connection to default SSH port
type SSHProber struct {
	defaultPort int
}

// NewSSHProber - 创建SSH探测器 / Create SSH prober
func NewSSHProber() *SSHProber {
	return &SSHProber{
		defaultPort: 22,
	}
}

// Type 返回探测类型 / Returns probe type
func (p *SSHProber) Type() string {
	return "ssh"
}

// Probe 执行SSH连通性检查 / Execute SSH connectivity check
// 步骤：
// 1. 建立TCP连接到目标SSH端口
// 2. 读取SSH协议标识（以"SSH-"开头的banner）
// 3. 根据响应判断SSH服务状态
//
// Steps:
// 1. Establish TCP connection to target SSH port
// 2. Read SSH protocol identification (banner starting with "SSH-")
// 3. Determine SSH service status based on response
func (p *SSHProber) Probe(host string, port int, timeout time.Duration) ProbeResult {
	if port <= 0 {
		port = p.defaultPort
	}

	result := ProbeResult{
		Target:  fmt.Sprintf("%s:%d", host, port),
		Type:    "ssh",
		Host:    host,
		Status:  StatusFail,
	}

	address := fmt.Sprintf("%s:%d", host, port)
	start := time.Now()

	// 建立TCP连接 / Establish TCP connection
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		elapsed := time.Since(start)
		result.Latency = float64(elapsed.Milliseconds())
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			result.Status = StatusTimeout
			result.Message = fmt.Sprintf("SSH connection timeout after %v", elapsed)
		} else {
			result.Status = StatusFail
			result.Message = fmt.Sprintf("SSH connection refused: %v", err)
		}
		result.Timestamp = time.Now()
		return result
	}
	defer conn.Close()

	// 设置读取超时 / Set read timeout
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))

	// 读取SSH banner / Read SSH banner
	buf := make([]byte, 128)
	n, err := conn.Read(buf)
	elapsed := time.Since(start)
	result.Latency = float64(elapsed.Milliseconds())
	result.Timestamp = time.Now()

	if err != nil {
		result.Status = StatusFail
		result.Message = fmt.Sprintf("SSH banner read failed: %v", err)
		return result
	}

	banner := string(buf[:n])
	if len(banner) >= 4 && banner[:4] == "SSH-" {
		result.Status = StatusOK
		// 截取banner到换行符 / Truncate banner at newline
		for i, c := range banner {
			if c == '\n' {
				banner = banner[:i]
				break
			}
		}
		result.Message = fmt.Sprintf("SSH service available: %s, latency: %v", banner, elapsed)
	} else {
		result.Status = StatusFail
		result.Message = fmt.Sprintf("invalid SSH banner: %s", banner)
	}

	return result
}

// Package probe - TCP端口探测 / TCP Port Probe
package probe

import (
	"fmt"
	"net"
	"time"
)

// TCPProber - TCP端口探测器 / TCP port prober
type TCPProber struct{}

// Type 返回探测类型 / Returns probe type
func (p *TCPProber) Type() string {
	return "tcp"
}

// Probe 执行TCP端口探测 / Execute TCP port probe
// 尝试建立TCP连接来检测端口是否可达
// Attempts to establish a TCP connection to check port reachability
func (p *TCPProber) Probe(host string, port int, timeout time.Duration) ProbeResult {
	result := ProbeResult{
		Target:  fmt.Sprintf("%s:%d", host, port),
		Type:    "tcp",
		Host:    host,
		Status:  StatusFail,
	}

	if port <= 0 {
		result.Message = "invalid port number"
		return result
	}

	start := time.Now()
	address := fmt.Sprintf("%s:%d", host, port)

	// 尝试TCP连接 / Attempt TCP connection
	conn, err := net.DialTimeout("tcp", address, timeout)
	elapsed := time.Since(start)
	result.Latency = float64(elapsed.Milliseconds())

	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			result.Status = StatusTimeout
			result.Message = fmt.Sprintf("connection timeout after %v", elapsed)
		} else {
			result.Status = StatusFail
			result.Message = fmt.Sprintf("connection refused: %v", err)
		}
		return result
	}

	// 连接成功，关闭连接 / Connection successful, close it
	conn.Close()
	result.Status = StatusOK
	result.Message = fmt.Sprintf("port %d is open, latency: %v", port, elapsed)
	result.Timestamp = time.Now()

	return result
}

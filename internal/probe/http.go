// Package probe - HTTP健康检查 / HTTP Health Check Probe
package probe

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

// HTTPProber - HTTP健康检查探测器 / HTTP health check prober
type HTTPProber struct {
	client *http.Client
}

// NewHTTPProber - 创建HTTP探测器 / Create HTTP prober
func NewHTTPProber() *HTTPProber {
	return &HTTPProber{
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // 允许自签名证书 / Allow self-signed certs
				},
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// 最多跟随5次重定向 / Follow up to 5 redirects
				if len(via) >= 5 {
					return fmt.Errorf("too many redirects")
				}
				return nil
			},
		},
	}
}

// Type 返回探测类型 / Returns probe type
func (p *HTTPProber) Type() string {
	return "http"
}

// Probe 执行HTTP探测 / Execute HTTP probe
func (p *HTTPProber) Probe(host string, port int, timeout time.Duration) ProbeResult {
	url := fmt.Sprintf("http://%s:%d", host, port)
	return p.ProbeURL(url, timeout)
}

// ProbeURL 对指定URL执行HTTP探测 / Execute HTTP probe against specified URL
// 发送GET请求并检查响应状态码
// Sends GET request and checks response status code
func (p *HTTPProber) ProbeURL(url string, timeout time.Duration) ProbeResult {
	result := ProbeResult{
		Target: url,
		Type:   "http",
		Host:   url,
		Status: StatusFail,
	}

	if url == "" {
		result.Message = "empty URL"
		return result
	}

	if timeout > 0 {
		p.client.Timeout = timeout
	}

	start := time.Now()
	resp, err := p.client.Get(url)
	elapsed := time.Since(start)
	result.Latency = float64(elapsed.Milliseconds())
	result.Timestamp = time.Now()

	if err != nil {
		result.Status = StatusFail
		result.Message = fmt.Sprintf("request failed: %v", err)
		return result
	}
	defer resp.Body.Close()

	// 检查状态码 / Check status code
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		result.Status = StatusOK
		result.Message = fmt.Sprintf("HTTP %d OK, latency: %v", resp.StatusCode, elapsed)
	} else if resp.StatusCode >= 500 {
		result.Status = StatusFail
		result.Message = fmt.Sprintf("HTTP %d Server Error", resp.StatusCode)
	} else {
		result.Status = StatusFail
		result.Message = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}

	return result
}

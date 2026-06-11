// Package alert - 通知器 / Notifier
// 支持Webhook和邮件通知方式
// Supports Webhook and Email notification methods
package alert

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"github.com/gitstq/ServerPulse/pkg/config"
)

// Notifier - 通知器 / Notifier
type Notifier struct {
	webhookCfg config.WebhookConfig
	emailCfg   config.EmailConfig
	httpClient *http.Client
}

// NewNotifier - 创建通知器 / Create notifier
func NewNotifier(webhookCfg config.WebhookConfig, emailCfg config.EmailConfig) *Notifier {
	return &Notifier{
		webhookCfg: webhookCfg,
		emailCfg:   emailCfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Notify - 发送通知 / Send notification
// 根据配置尝试所有启用的通知方式
// Attempts all enabled notification methods based on configuration
func (n *Notifier) Notify(alert Alert) error {
	var errs []string

	// Webhook通知 / Webhook notification
	if n.webhookCfg.Enabled {
		if err := n.sendWebhook(alert); err != nil {
			errs = append(errs, "webhook: "+err.Error())
		}
	}

	// 邮件通知 / Email notification
	if n.emailCfg.Enabled {
		if err := n.sendEmail(alert); err != nil {
			errs = append(errs, "email: "+err.Error())
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("notification errors: %s", strings.Join(errs, "; "))
	}
	return nil
}

// NotifyAggregated - 发送聚合告警通知 / Send aggregated alert notification
func (n *Notifier) NotifyAggregated(agg *AggregatedAlert) error {
	var errs []string

	if n.webhookCfg.Enabled {
		if err := n.sendWebhookAggregated(agg); err != nil {
			errs = append(errs, "webhook: "+err.Error())
		}
	}

	if n.emailCfg.Enabled {
		if err := n.sendEmailAggregated(agg); err != nil {
			errs = append(errs, "email: "+err.Error())
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("notification errors: %s", strings.Join(errs, "; "))
	}
	return nil
}

// sendWebhook - 发送Webhook通知 / Send Webhook notification
func (n *Notifier) sendWebhook(alert Alert) error {
	if n.webhookCfg.URL == "" {
		return fmt.Errorf("webhook URL is empty")
	}

	// 构建通知载荷 / Build notification payload
	payload := map[string]interface{}{
		"id":        alert.ID,
		"server_id": alert.ServerID,
		"timestamp": alert.Timestamp.Format(time.RFC3339),
		"level":     alert.Level,
		"type":      alert.Type,
		"metric":    alert.Metric,
		"value":     alert.Value,
		"threshold": alert.Threshold,
		"message":   alert.Message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// 创建HTTP请求 / Create HTTP request
	req, err := http.NewRequest("POST", n.webhookCfg.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 添加签名（如果配置了密钥）/ Add signature if secret is configured
	if n.webhookCfg.Secret != "" {
		signature := n.signPayload(body, n.webhookCfg.Secret)
		req.Header.Set("X-ServerPulse-Signature", signature)
	}

	// 发送请求 / Send request
	resp, err := n.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status: %d", resp.StatusCode)
	}

	return nil
}

// sendWebhookAggregated - 发送聚合告警Webhook通知 / Send aggregated alert Webhook notification
func (n *Notifier) sendWebhookAggregated(agg *AggregatedAlert) error {
	if n.webhookCfg.URL == "" {
		return fmt.Errorf("webhook URL is empty")
	}

	payload := map[string]interface{}{
		"id":        agg.ID,
		"timestamp": agg.Timestamp.Format(time.RFC3339),
		"level":     agg.Level,
		"count":     agg.Count,
		"server_id": agg.ServerID,
		"summary":   agg.Summary,
		"type":      "aggregated",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", n.webhookCfg.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	if n.webhookCfg.Secret != "" {
		signature := n.signPayload(body, n.webhookCfg.Secret)
		req.Header.Set("X-ServerPulse-Signature", signature)
	}

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status: %d", resp.StatusCode)
	}

	return nil
}

// sendEmail - 发送邮件通知 / Send email notification
func (n *Notifier) sendEmail(alert Alert) error {
	if len(n.emailCfg.To) == 0 {
		return fmt.Errorf("no email recipients configured")
	}

	subject := fmt.Sprintf("[ServerPulse] [%s] %s: %s = %.2f", alert.Level, alert.Type, alert.Metric, alert.Value)
	body := fmt.Sprintf(
		"告警通知 / Alert Notification\n\n"+
			"级别 / Level: %s\n"+
			"类型 / Type: %s\n"+
			"指标 / Metric: %s\n"+
			"当前值 / Value: %.2f\n"+
			"阈值 / Threshold: %.2f\n"+
			"时间 / Time: %s\n"+
			"消息 / Message: %s\n",
		alert.Level, alert.Type, alert.Metric, alert.Value, alert.Threshold,
		alert.Timestamp.Format(time.RFC3339), alert.Message,
	)

	return n.sendMail(subject, body)
}

// sendEmailAggregated - 发送聚合告警邮件通知 / Send aggregated alert email notification
func (n *Notifier) sendEmailAggregated(agg *AggregatedAlert) error {
	if len(n.emailCfg.To) == 0 {
		return fmt.Errorf("no email recipients configured")
	}

	subject := fmt.Sprintf("[ServerPulse] [%s] 聚合告警: %d条告警", agg.Level, agg.Count)
	body := fmt.Sprintf(
		"聚合告警通知 / Aggregated Alert Notification\n\n"+
			"级别 / Level: %s\n"+
			"告警数量 / Count: %d\n"+
			"摘要 / Summary: %s\n"+
			"时间 / Time: %s\n",
		agg.Level, agg.Count, agg.Summary,
		agg.Timestamp.Format(time.RFC3339),
	)

	return n.sendMail(subject, body)
}

// sendMail - 底层邮件发送 / Low-level mail sending
func (n *Notifier) sendMail(subject, body string) error {
	addr := fmt.Sprintf("%s:%d", n.emailCfg.SMTPHost, n.emailCfg.SMTPPort)

	// 构建邮件内容 / Build email content
	var msg bytes.Buffer
	msg.WriteString("From: " + n.emailCfg.From + "\r\n")
	msg.WriteString("To: " + strings.Join(n.emailCfg.To, ",") + "\r\n")
	msg.WriteString("Subject: " + subject + "\r\n")
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)

	// 发送邮件 / Send email
	auth := smtp.PlainAuth("", n.emailCfg.Username, n.emailCfg.Password, n.emailCfg.SMTPHost)
	return smtp.SendMail(addr, auth, n.emailCfg.From, n.emailCfg.To, msg.Bytes())
}

// signPayload - 对载荷进行HMAC-SHA256签名 / HMAC-SHA256 sign payload
func (n *Notifier) signPayload(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

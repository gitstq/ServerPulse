// Package tui - 告警视图 / Alerts View
// 显示当前活跃告警和历史告警列表
// Displays current active alerts and historical alert list
package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AlertItem - 告警项 / Alert item
type AlertItem struct {
	ID        string
	Level     string // warning, critical
	Type      string // threshold, anomaly
	Metric    string
	Value     float64
	Threshold float64
	Message   string
	Timestamp time.Time
	Acked     bool
}

// AlertsModel - 告警视图模型 / Alerts view model
type AlertsModel struct {
	width    int
	height   int
	alerts   []AlertItem
	selected int
}

// NewAlertsModel - 创建告警视图模型 / Create alerts view model
func NewAlertsModel() AlertsModel {
	return AlertsModel{
		alerts: make([]AlertItem, 0),
	}
}

// Update - 处理消息 / Handle messages
func (m AlertsModel) Update(msg tea.Msg) (AlertsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.alerts)-1 {
				m.selected++
			}
		case "a":
			// 确认选中的告警 / Acknowledge selected alert
			if m.selected < len(m.alerts) {
				m.alerts[m.selected].Acked = true
			}
		}
	}
	return m, nil
}

// SetAlerts - 设置告警列表 / Set alert list
func (m *AlertsModel) SetAlerts(alerts []AlertItem) {
	m.alerts = alerts
	if m.selected >= len(alerts) {
		m.selected = 0
	}
}

// GetActiveCount - 获取活跃告警数量 / Get active alert count
func (m AlertsModel) GetActiveCount() int {
	count := 0
	for _, a := range m.alerts {
		if !a.Acked {
			count++
		}
	}
	return count
}

// View - 渲染告警视图 / Render alerts view
func (m AlertsModel) View() string {
	// 标题 / Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF6B6B")).
		MarginBottom(1)

	activeCount := m.GetActiveCount()
	totalCount := len(m.alerts)
	title := fmt.Sprintf("告警列表 (活跃: %d / 总计: %d)", activeCount, totalCount)

	var sections []string
	sections = append(sections, titleStyle.Render(title))

	if len(m.alerts) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))
		sections = append(sections, emptyStyle.Render("暂无告警 / No alerts"))
	} else {
		// 渲染告警列表 / Render alert list
		for i, alert := range m.alerts {
			sections = append(sections, m.renderAlertItem(i, alert))
		}
	}

	// 操作提示 / Action hints
	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		MarginTop(1)
	hints := "↑↓/jk: 选择 | a: 确认告警"
	sections = append(sections, hintStyle.Render(hints))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderAlertItem - 渲染单个告警项 / Render single alert item
func (m AlertsModel) renderAlertItem(index int, alert AlertItem) string {
	// 根据级别选择颜色 / Choose color based on level
	var levelColor, levelText string
	switch alert.Level {
	case "critical":
		levelColor = "#FF0000"
		levelText = "严重"
	case "warning":
		levelColor = "#FFA500"
		levelText = "警告"
	default:
		levelColor = "#888888"
		levelText = "信息"
	}

	// 根据类型选择标签 / Choose label based on type
	var typeText string
	switch alert.Type {
	case "threshold":
		typeText = "阈值"
	case "anomaly":
		typeText = "异常"
	default:
		typeText = alert.Type
	}

	// 构建告警行 / Build alert line
	levelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(levelColor))

	timeStr := alert.Timestamp.Format("15:04:05")
	line := fmt.Sprintf("[%s] [%s] %s: %s = %.2f (阈值: %.2f)",
		timeStr, levelText, typeText, alert.Metric, alert.Value, alert.Threshold)

	// 选中状态 / Selected state
	if index == m.selected {
		selectedStyle := lipgloss.NewStyle().
			Background(lipgloss.Color("#333333"))
		return selectedStyle.Render(levelStyle.Render(line))
	}

	if alert.Acked {
		// 已确认的告警显示为灰色 / Acknowledged alerts shown in gray
		ackedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
		return ackedStyle.Render(line)
	}

	return levelStyle.Render(line)
}

// unused - 避免编译器警告 / Avoid compiler warnings
var _ = strings.TrimSpace

// Package tui - 服务器列表视图 / Servers View
// 显示已注册的服务器列表和状态
// Displays registered server list and status
package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ServerItem - 服务器项 / Server item
type ServerItem struct {
	ID       string
	Name     string
	IP       string
	OS       string
	Status   string // online, offline
	LastSeen time.Time
}

// ServersModel - 服务器列表视图模型 / Servers view model
type ServersModel struct {
	width    int
	height   int
	servers  []ServerItem
	selected int
}

// NewServersModel - 创建服务器列表模型 / Create servers view model
func NewServersModel() ServersModel {
	return ServersModel{
		servers: make([]ServerItem, 0),
	}
}

// Update - 处理消息 / Handle messages
func (m ServersModel) Update(msg tea.Msg) (ServersModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.servers)-1 {
				m.selected++
			}
		}
	}
	return m, nil
}

// SetServers - 设置服务器列表 / Set server list
func (m *ServersModel) SetServers(servers []ServerItem) {
	m.servers = servers
	if m.selected >= len(servers) {
		m.selected = 0
	}
}

// View - 渲染服务器列表视图 / Render servers view
func (m ServersModel) View() string {
	// 标题 / Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00BFFF")).
		MarginBottom(1)

	title := fmt.Sprintf("服务器列表 (%d 台)", len(m.servers))
	var sections []string
	sections = append(sections, titleStyle.Render(title))

	if len(m.servers) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))
		sections = append(sections, emptyStyle.Render("暂无注册服务器 / No registered servers"))
	} else {
		// 表头 / Table header
		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#AAAAAA"))

		header := fmt.Sprintf("  %-20s %-15s %-20s %-10s %-15s",
			"名称 Name", "IP地址 IP", "操作系统 OS", "状态 Status", "最后在线 LastSeen")
		sections = append(sections, headerStyle.Render(header))

		// 分隔线 / Separator
		separator := "  " + "─────────────────────────────────────────────────────────────────────────────────────"
		sections = append(sections, separator)

		// 服务器列表 / Server list
		for i, server := range m.servers {
			sections = append(sections, m.renderServerItem(i, server))
		}
	}

	// 操作提示 / Action hints
	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		MarginTop(1)
	hints := "↑↓/jk: 选择服务器 | Tab: 返回仪表板"
	sections = append(sections, hintStyle.Render(hints))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderServerItem - 渲染单个服务器项 / Render single server item
func (m ServersModel) renderServerItem(index int, server ServerItem) string {
	// 状态文本 / Status text
	var statusText string
	switch server.Status {
	case "online":
		statusText = "在线"
	case "offline":
		statusText = "离线"
	default:
		statusText = "未知"
	}

	// 最后在线时间 / Last seen time
	lastSeen := "从未"
	if !server.LastSeen.IsZero() {
		lastSeen = server.LastSeen.Format("2006-01-02 15:04:05")
	}

	line := fmt.Sprintf("  %-20s %-15s %-20s %-10s %-15s",
		server.Name, server.IP, server.OS, statusText, lastSeen)

	// 选中状态 / Selected state
	if index == m.selected {
		selectedStyle := lipgloss.NewStyle().Background(lipgloss.Color("#333333"))
		return selectedStyle.Render(line)
	}

	return line
}

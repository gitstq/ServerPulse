// Package tui - TUI主应用 / TUI Main Application
// 基于Bubbletea的终端UI仪表板，支持实时指标显示、告警列表、服务器切换
// Terminal UI dashboard based on Bubbletea, supports real-time metrics, alerts, server switching
package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// 视图索引 / View indices
const (
	viewDashboard int = iota
	viewAlerts
	viewServers
)

// Model - TUI主模型 / TUI main model
type Model struct {
	currentView    int              // 当前视图 / Current view
	width          int              // 终端宽度 / Terminal width
	height         int              // 终端高度 / Terminal height
	dashboardModel DashboardModel   // 仪表板模型 / Dashboard model
	alertsModel    AlertsModel      // 告警视图模型 / Alerts view model
	serversModel   ServersModel     // 服务器列表模型 / Servers view model
	refreshRate    time.Duration    // 刷新频率 / Refresh rate
	lastRefresh    time.Time        // 上次刷新 / Last refresh
	quitting       bool             // 是否退出 / Whether quitting
	err            error            // 错误 / Error
}

// tickMsg - 定时刷新消息 / Tick message for periodic refresh
type tickMsg time.Time

// NewModel - 创建TUI模型 / Create TUI model
func NewModel(refreshRate time.Duration) Model {
	if refreshRate <= 0 {
		refreshRate = 1 * time.Second
	}
	return Model{
		currentView:  viewDashboard,
		refreshRate:  refreshRate,
		lastRefresh:  time.Now(),
		dashboardModel: NewDashboardModel(),
		alertsModel:    NewAlertsModel(),
		serversModel:   NewServersModel(),
	}
}

// Init - 初始化TUI / Initialize TUI
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(m.refreshRate),
	)
}

// tickCmd - 创建定时命令 / Create tick command
func tickCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update - 处理消息更新 / Handle message updates
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.dashboardModel.width = msg.Width
		m.dashboardModel.height = msg.Height
		m.alertsModel.width = msg.Width
		m.alertsModel.height = msg.Height
		m.serversModel.width = msg.Width
		m.serversModel.height = msg.Height
	case tickMsg:
		m.lastRefresh = time.Now()
		// 更新仪表板数据 / Update dashboard data
		m.dashboardModel.UpdateData()
	case teaErrMsg:
		m.err = msg
	}

	// 根据当前视图分发消息 / Dispatch message based on current view
	switch m.currentView {
	case viewDashboard:
		_, cmd := m.dashboardModel.Update(msg)
		return m, cmd
	case viewAlerts:
		_, cmd := m.alertsModel.Update(msg)
		return m, cmd
	case viewServers:
		_, cmd := m.serversModel.Update(msg)
		return m, cmd
	}

	return m, tickCmd(m.refreshRate)
}

// teaErrMsg - 错误消息 / Error message
type teaErrMsg struct{ err error }

func (e teaErrMsg) Error() string { return e.err.Error() }

// handleKeyMsg - 处理按键消息 / Handle key messages
func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	case "1", "tab":
		m.currentView = viewDashboard
	case "2":
		m.currentView = viewAlerts
	case "3":
		m.currentView = viewServers
	case "r":
		// 手动刷新 / Manual refresh
		m.dashboardModel.UpdateData()
	}
	return m, nil
}

// View - 渲染TUI视图 / Render TUI view
func (m Model) View() string {
	if m.quitting {
		return "ServerPulse 已退出 / ServerPulse exited\n"
	}

	if m.err != nil {
		return fmt.Sprintf("错误: %v\n", m.err)
	}

	// 渲染标题栏 / Render title bar
	title := m.renderTitle()

	// 渲染标签栏 / Render tab bar
	tabs := m.renderTabs()

	// 渲染当前视图内容 / Render current view content
	var content string
	switch m.currentView {
	case viewDashboard:
		content = m.dashboardModel.View()
	case viewAlerts:
		content = m.alertsModel.View()
	case viewServers:
		content = m.serversModel.View()
	}

	// 渲染底部状态栏 / Render bottom status bar
	status := m.renderStatus()

	// 组合所有部分 / Combine all parts
	return lipgloss.JoinVertical(lipgloss.Left, title, tabs, content, status)
}

// renderTitle - 渲染标题 / Render title
func (m Model) renderTitle() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Background(lipgloss.Color("#1A1A2E")).
		Padding(0, 2)

	return titleStyle.Render(" ServerPulse - 轻量级服务器智能监控引擎 ")
}

// renderTabs - 渲染标签栏 / Render tab bar
func (m Model) renderTabs() string {
	tabNames := []string{"1: 仪表板 Dashboard", "2: 告警 Alerts", "3: 服务器 Servers"}

	var tabs []string
	for i, name := range tabNames {
		style := lipgloss.NewStyle().Padding(0, 2)
		if i == m.currentView {
			style = style.Bold(true).Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#7D56F4"))
		} else {
			style = style.Foreground(lipgloss.Color("#888888"))
		}
		tabs = append(tabs, style.Render(name))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

// renderStatus - 渲染状态栏 / Render status bar
func (m Model) renderStatus() string {
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Padding(0, 2)

	alertCount := m.alertsModel.GetActiveCount()
	alertInfo := ""
	if alertCount > 0 {
		alertInfo = fmt.Sprintf(" | 告警: %d", alertCount)
	}

	uptime := time.Since(m.lastRefresh)
	info := fmt.Sprintf("刷新: %v | 按 Tab/1-3 切换视图 | q 退出%s", uptime.Round(time.Millisecond), alertInfo)

	return statusStyle.Render(info)
}

// formatBar - 格式化进度条 / Format progress bar
func formatBar(percent float64, width int) string {
	if width <= 0 {
		width = 30
	}
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	filled := int(percent / 100 * float64(width))
	empty := width - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)

	// 根据使用率选择颜色 / Choose color based on usage
	color := "#00FF00" // 绿色 / Green
	if percent >= 80 {
		color = "#FFA500" // 橙色 / Orange
	}
	if percent >= 95 {
		color = "#FF0000" // 红色 / Red
	}

	style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	return style.Render(bar)
}

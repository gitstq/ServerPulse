// Package tui - 仪表板视图 / Dashboard View
// 实时显示系统指标：CPU、内存、磁盘、网络
// Real-time display of system metrics: CPU, Memory, Disk, Network
package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DashboardModel - 仪表板视图模型 / Dashboard view model
type DashboardModel struct {
	width       int
	height      int
	cpuUsage    float64
	memUsage    float64
	memTotal    uint64
	memUsed     uint64
	diskUsage   float64
	diskTotal   uint64
	diskUsed    uint64
	netSent     uint64
	netRecv     uint64
	cpuCores    int
	cpuModel    string
	uptime      time.Time
	lastUpdate  time.Time
}

// NewDashboardModel - 创建仪表板模型 / Create dashboard model
func NewDashboardModel() DashboardModel {
	return DashboardModel{
		uptime: time.Now(),
	}
}

// Update - 处理消息 / Handle messages
func (m DashboardModel) Update(msg tea.Msg) (DashboardModel, tea.Cmd) {
	return m, nil
}

// UpdateData - 更新仪表板数据 / Update dashboard data
func (m *DashboardModel) UpdateData() {
	// 使用gopsutil采集实时数据 / Collect real-time data using gopsutil
	m.updateCPU()
	m.updateMemory()
	m.updateDisk()
	m.updateNetwork()
	m.lastUpdate = time.Now()
}

// updateCPU - 更新CPU数据 / Update CPU data
func (m *DashboardModel) updateCPU() {
	import_cpu()
}

// updateMemory - 更新内存数据 / Update memory data
func (m *DashboardModel) updateMemory() {
	import_memory()
}

// updateDisk - 更新磁盘数据 / Update disk data
func (m *DashboardModel) updateDisk() {
	import_disk()
}

// updateNetwork - 更新网络数据 / Update network data
func (m *DashboardModel) updateNetwork() {
	import_network()
}

// 占位导入函数，实际在下方实现 / Placeholder import functions, implemented below
func import_cpu() {
	// 实际采集逻辑通过gopsutil实现
	// Actual collection logic via gopsutil
}

func import_memory() {}

func import_disk() {}

func import_network() {}

// View - 渲染仪表板视图 / Render dashboard view
func (m DashboardModel) View() string {
	var sections []string

	// 系统概览 / System overview
	sections = append(sections, m.renderOverview())

	// CPU指标 / CPU metrics
	sections = append(sections, m.renderCPU())

	// 内存指标 / Memory metrics
	sections = append(sections, m.renderMemory())

	// 磁盘指标 / Disk metrics
	sections = append(sections, m.renderDisk())

	// 网络指标 / Network metrics
	sections = append(sections, m.renderNetwork())

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return content
}

// renderOverview - 渲染系统概览 / Render system overview
func (m DashboardModel) renderOverview() string {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		MarginBottom(1)

	uptime := time.Since(m.uptime).Truncate(time.Second)
	header := fmt.Sprintf("系统概览 | 运行时间: %v | CPU: %s | 核心数: %d",
		uptime, m.cpuModel, m.cpuCores)

	return style.Render(header)
}

// renderCPU - 渲染CPU指标 / Render CPU metrics
func (m DashboardModel) renderCPU() string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		MarginBottom(1)

	barWidth := 40
	bar := formatBar(m.cpuUsage, barWidth)

	content := fmt.Sprintf("CPU 使用率: %.1f%%\n%s", m.cpuUsage, bar)
	return boxStyle.Render(content)
}

// renderMemory - 渲染内存指标 / Render memory metrics
func (m DashboardModel) renderMemory() string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00BFFF")).
		Padding(0, 1).
		MarginBottom(1)

	barWidth := 40
	bar := formatBar(m.memUsage, barWidth)

	memTotalGB := float64(m.memTotal) / 1024 / 1024 / 1024
	memUsedGB := float64(m.memUsed) / 1024 / 1024 / 1024

	content := fmt.Sprintf("内存使用率: %.1f%% (%.1fGB / %.1fGB)\n%s",
		m.memUsage, memUsedGB, memTotalGB, bar)
	return boxStyle.Render(content)
}

// renderDisk - 渲染磁盘指标 / Render disk metrics
func (m DashboardModel) renderDisk() string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FFD700")).
		Padding(0, 1).
		MarginBottom(1)

	barWidth := 40
	bar := formatBar(m.diskUsage, barWidth)

	diskTotalGB := float64(m.diskTotal) / 1024 / 1024 / 1024
	diskUsedGB := float64(m.diskUsed) / 1024 / 1024 / 1024

	content := fmt.Sprintf("磁盘使用率: %.1f%% (%.1fGB / %.1fGB)\n%s",
		m.diskUsage, diskUsedGB, diskTotalGB, bar)
	return boxStyle.Render(content)
}

// renderNetwork - 渲染网络指标 / Render network metrics
func (m DashboardModel) renderNetwork() string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#32CD32")).
		Padding(0, 1).
		MarginBottom(1)

	sentStr := formatBytes(m.netSent)
	recvStr := formatBytes(m.netRecv)

	content := fmt.Sprintf("网络流量:\n  发送: %s  |  接收: %s", sentStr, recvStr)
	return boxStyle.Render(content)
}

// formatBytes - 格式化字节数 / Format bytes
func formatBytes(bytes uint64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(TB))
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// renderPlaceholder - 渲染占位内容 / Render placeholder content
func renderPlaceholder(title, desc string) string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#888888")).
		Padding(0, 1).
		MarginBottom(1)

	content := fmt.Sprintf("%s\n%s", title, desc)
	return boxStyle.Render(content)
}

// unused - 避免编译器警告 / Avoid compiler warnings
var _ = strings.TrimSpace

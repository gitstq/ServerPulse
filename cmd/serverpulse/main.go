// Package main - ServerPulse CLI入口 / ServerPulse CLI Entry Point
// 轻量级服务器智能监控引擎
// Lightweight Intelligent Server Monitoring Engine
//
// 支持子命令 / Supported subcommands:
//   serve    - 启动监控服务（Hub模式）/ Start monitoring service (Hub mode)
//   agent    - 启动Agent客户端 / Start Agent client
//   tui      - 启动TUI仪表板 / Start TUI dashboard
//   check    - 执行一次性检查 / Execute one-time check
//   snapshot - 快照管理 / Snapshot management
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/gitstq/ServerPulse/internal/alert"
	"github.com/gitstq/ServerPulse/internal/anomaly"
	"github.com/gitstq/ServerPulse/internal/collector"
	"github.com/gitstq/ServerPulse/internal/probe"
	"github.com/gitstq/ServerPulse/internal/server"
	"github.com/gitstq/ServerPulse/internal/snapshot"
	"github.com/gitstq/ServerPulse/internal/storage"
	"github.com/gitstq/ServerPulse/internal/tui"
	"github.com/gitstq/ServerPulse/pkg/config"
)

// 版本信息 / Version information
var (
	Version   = "1.0.0"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

// 全局配置 / Global configuration
var (
	cfgFile string
	cfg     *config.Config
)

func main() {
	// 设置根命令 / Set up root command
	rootCmd := &cobra.Command{
		Use:   "serverpulse",
		Short: "ServerPulse - 轻量级服务器智能监控引擎",
		Long: `ServerPulse - Lightweight Intelligent Server Monitoring Engine

支持系统指标采集、多协议探测、AI异常检测、智能告警和TUI仪表板。
Supports system metrics collection, multi-protocol probing, AI anomaly detection,
intelligent alerting and TUI dashboard.`,
		Version: Version,
	}

	// 全局标志 / Global flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "配置文件路径 / Config file path (default: configs/serverpulse.yaml)")

	// 添加子命令 / Add subcommands
	rootCmd.AddCommand(serveCmd())
	rootCmd.AddCommand(agentCmd())
	rootCmd.AddCommand(tuiCmd())
	rootCmd.AddCommand(checkCmd())
	rootCmd.AddCommand(snapshotCmd())
	rootCmd.AddCommand(versionCmd())

	// 执行 / Execute
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// loadConfig - 加载配置 / Load configuration
func loadConfig() *config.Config {
	if cfg != nil {
		return cfg
	}

	// 尝试加载配置文件 / Try to load config file
	path := cfgFile
	if path == "" {
		// 默认路径 / Default path
		if _, err := os.Stat("configs/serverpulse.yaml"); err == nil {
			path = "configs/serverpulse.yaml"
		} else if _, err := os.Stat("/etc/serverpulse/serverpulse.yaml"); err == nil {
			path = "/etc/serverpulse/serverpulse.yaml"
		}
	}

	if path != "" {
		var err error
		cfg, err = config.LoadFromFile(path)
		if err != nil {
			log.Printf("警告: 加载配置文件失败，使用默认配置: %v", err)
			cfg = config.Default()
		}
	} else {
		cfg = config.Default()
	}

	return cfg
}

// versionCmd - 版本命令 / Version command
func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "显示版本信息 / Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("ServerPulse %s\n", Version)
			fmt.Printf("  Git Commit: %s\n", GitCommit)
			fmt.Printf("  Build Date: %s\n", BuildDate)
			fmt.Printf("  Go Version: %s\n", runtime.Version())
			fmt.Printf("  OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		},
	}
}

// serveCmd - 启动监控服务命令 / Start monitoring service command
func serveCmd() *cobra.Command {
	var port int
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "启动监控服务 / Start monitoring service",
		Long:  "启动ServerPulse监控服务，包含指标采集、异常检测和告警功能。\nStart ServerPulse monitoring service with metrics collection, anomaly detection and alerting.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := loadConfig()

			// 如果指定了端口，覆盖配置 / Override config if port specified
			if port > 0 {
				c.Hub.Bind = fmt.Sprintf(":%d", port)
			}

			log.Printf("[ServerPulse] v%s 启动中...", Version)
			log.Printf("[ServerPulse] 服务器名称: %s", c.Server.Name)
			log.Printf("[ServerPulse] 采集间隔: %v", c.Collector.Interval)
			log.Printf("[ServerPulse] 存储路径: %s", c.Storage.Path)

			// 初始化存储 / Initialize storage
			store, err := storage.NewSQLiteStorage(c.Storage.Path)
			if err != nil {
				return fmt.Errorf("初始化存储失败: %w", err)
			}
			defer store.Close()

			// 初始化采集器 / Initialize collector
			sysCollector := collector.NewSystemCollector(
				c.Server.ID,
				c.Collector.CPU,
				c.Collector.Memory,
				c.Collector.Disk,
				c.Collector.Network,
				c.Collector.Docker,
			)

			// 初始化异常检测器 / Initialize anomaly detector
			detector := anomaly.NewDetector(c.Anomaly.WindowSize, c.Anomaly.ZScoreThreshold)

			// 初始化告警管理器 / Initialize alert manager
			alertMgr := alert.NewManager(c.Alert)

			// 初始化通知器 / Initialize notifier
			notifier := alert.NewNotifier(c.Notify.Webhook, c.Notify.Email)
			alertMgr.AddCallback(func(a alert.Alert) {
				log.Printf("[告警] %s", a.Message)
				if err := notifier.Notify(a); err != nil {
					log.Printf("[通知] 发送失败: %v", err)
				}
			})

			// 启动Hub服务（如果启用）/ Start Hub service (if enabled)
			if c.Hub.Enabled {
				hub := server.NewHub(c.Hub, store)
				if err := hub.Start(); err != nil {
					return fmt.Errorf("启动Hub失败: %w", err)
				}
				defer hub.Stop()
			}

			// 启动采集循环 / Start collection loop
			log.Println("[ServerPulse] 开始采集指标...")
			go func() {
				ticker := time.NewTicker(c.Collector.Interval)
				defer ticker.Stop()
				for range ticker.C {
					// 采集数据 / Collect data
					metrics, err := sysCollector.CollectAll()
					if err != nil {
						log.Printf("[采集] 错误: %v", err)
						continue
					}

					// 存储指标数据 / Store metric data
					metricList := sysCollector.ToMetricList(metrics)
					storageMetrics := make([]storage.MetricData, len(metricList))
					for i, m := range metricList {
						storageMetrics[i] = storage.MetricData{
							ServerID:  c.Server.ID,
							Timestamp: m.Timestamp,
							Metric:    m.Name,
							Value:     m.Value,
							Unit:      m.Unit,
							Labels:    m.Labels,
						}
					}
					if err := store.InsertMetrics(storageMetrics); err != nil {
						log.Printf("[存储] 写入失败: %v", err)
					}

					// 异常检测 / Anomaly detection
					if c.Anomaly.Enabled {
						for _, m := range metricList {
							result := detector.Detect(m.Name, m.Value)
							if result != nil && result.IsAnomaly {
								log.Printf("[异常] %s", result.Message)
								alertMgr.AddAnomalyAlert(c.Server.ID, m.Name, m.Value, result.ZScore, result.Severity)
							}
						}
					}

					// 阈值告警检查 / Threshold alert check
					if c.Alert.Enabled {
						for _, m := range metricList {
							alertMgr.CheckThreshold(c.Server.ID, m.Name, m.Value)
						}
					}
				}
			}()

			// 启动探测循环 / Start probe loop
			if len(c.Probe.Targets) > 0 {
				log.Println("[ServerPulse] 开始执行探测...")
				probeMgr := probe.NewProbeManager()
				go func() {
					ticker := time.NewTicker(c.Probe.Interval)
					defer ticker.Stop()
					for range ticker.C {
						for _, target := range c.Probe.Targets {
							result := probeMgr.Probe(target.Type, target.Host, target.Port, target.URL, target.Timeout)
							log.Printf("[探测] %s (%s): %s - %s",
								target.Name, target.Type, result.Status, result.Message)

							// 存储探测结果 / Store probe result
							store.InsertProbeResult(&storage.ProbeResult{
								ServerID:  c.Server.ID,
								Timestamp: result.Timestamp,
								Target:    target.Name,
								Type:      target.Type,
								Status:    string(result.Status),
								Latency:   result.Latency,
								Message:   result.Message,
							})
						}
					}
				}()
			}

			log.Println("[ServerPulse] 服务已启动，按 Ctrl+C 停止")

			// 等待退出信号 / Wait for exit signal
			select {}
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 0, "Hub监听端口 / Hub listen port")
	return cmd
}

// agentCmd - 启动Agent命令 / Start Agent command
func agentCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "agent",
		Short: "启动Agent客户端 / Start Agent client",
		Long:  "启动ServerPulse Agent，采集本地数据并上报到Hub。\nStart ServerPulse Agent to collect local data and report to Hub.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := loadConfig()

			if !c.Agent.Enabled && c.Agent.HubURL == "" {
				return fmt.Errorf("Agent未启用，请在配置文件中设置 agent.enabled: true")
			}

			log.Printf("[Agent] v%s 启动中...", Version)
			log.Printf("[Agent] Hub地址: %s", c.Agent.HubURL)

			agent := server.NewAgent(c.Agent, c.Server)
			if err := agent.Start(); err != nil {
				return fmt.Errorf("启动Agent失败: %w", err)
			}

			log.Println("[Agent] 已启动，按 Ctrl+C 停止")
			select {}
		},
	}
}

// tuiCmd - 启动TUI命令 / Start TUI command
func tuiCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tui",
		Short: "启动TUI仪表板 / Start TUI dashboard",
		Long:  "启动ServerPulse终端UI仪表板，实时显示系统指标和告警。\nStart ServerPulse terminal UI dashboard for real-time metrics and alerts.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := loadConfig()

			// 初始化TUI模型 / Initialize TUI model
			appModel := tui.NewModel(c.TUI.RefreshRate)

			// 启动Bubbletea程序 / Start Bubbletea program
			p := tea.NewProgram(appModel, tea.WithAltScreen())

			_, err := p.Run()
			if err != nil {
				return fmt.Errorf("TUI运行错误: %w", err)
			}

			return nil
		},
	}
}

// checkCmd - 一次性检查命令 / One-time check command
func checkCmd() *cobra.Command {
	var (
		checkCPU    bool
		checkMemory bool
		checkDisk   bool
		checkNet    bool
		checkAll    bool
		jsonOutput  bool
	)

	cmd := &cobra.Command{
		Use:   "check",
		Short: "执行一次性系统检查 / Execute one-time system check",
		Long:  "采集当前系统指标并输出，用于快速诊断。\nCollect and display current system metrics for quick diagnostics.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := loadConfig()

			// 确定检查项 / Determine check items
			if checkAll {
				checkCPU = true
				checkMemory = true
				checkDisk = true
				checkNet = true
			}

			sysCollector := collector.NewSystemCollector(
				c.Server.ID,
				checkCPU || !checkMemory && !checkDisk && !checkNet,
				checkMemory,
				checkDisk,
				checkNet,
				false,
			)

			metrics, err := sysCollector.CollectAll()
			if err != nil {
				return fmt.Errorf("采集失败: %w", err)
			}

			if jsonOutput {
				data, _ := json.MarshalIndent(metrics, "", "  ")
				fmt.Println(string(data))
				return nil
			}

			// 文本输出 / Text output
			fmt.Println("=== ServerPulse 系统检查报告 ===")
			fmt.Printf("时间: %s\n", metrics.Timestamp.Format("2006-01-02 15:04:05"))
			fmt.Println()

			if checkCPU || (!checkMemory && !checkDisk && !checkNet) {
				fmt.Printf("CPU 使用率: %.1f%%\n", metrics.CPU.UsagePercent)
				fmt.Printf("CPU 型号:   %s\n", metrics.CPU.ModelName)
				fmt.Printf("CPU 核心数: %d\n", metrics.CPU.Cores)
				fmt.Println()
			}

			if checkMemory {
				fmt.Printf("内存使用率: %.1f%%\n", metrics.Memory.UsagePercent)
				fmt.Printf("内存总量:   %.2f GB\n", float64(metrics.Memory.Total)/1024/1024/1024)
				fmt.Printf("内存已用:   %.2f GB\n", float64(metrics.Memory.Used)/1024/1024/1024)
				fmt.Printf("内存可用:   %.2f GB\n", float64(metrics.Memory.Available)/1024/1024/1024)
				fmt.Println()
			}

			if checkDisk {
				for _, d := range metrics.Disk {
					fmt.Printf("磁盘 [%s] (%s):\n", d.Mountpoint, d.Device)
					fmt.Printf("  使用率: %.1f%%\n", d.UsagePercent)
					fmt.Printf("  总量:   %.2f GB\n", float64(d.Total)/1024/1024/1024)
					fmt.Printf("  已用:   %.2f GB\n", float64(d.Used)/1024/1024/1024)
					fmt.Printf("  可用:   %.2f GB\n", float64(d.Free)/1024/1024/1024)
				}
				fmt.Println()
			}

			if checkNet {
				for _, n := range metrics.Network {
					fmt.Printf("网络 [%s]:\n", n.Name)
					fmt.Printf("  发送: %s\n", formatBytes(n.BytesSent))
					fmt.Printf("  接收: %s\n", formatBytes(n.BytesRecv))
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&checkCPU, "cpu", false, "检查CPU / Check CPU")
	cmd.Flags().BoolVar(&checkMemory, "memory", false, "检查内存 / Check memory")
	cmd.Flags().BoolVar(&checkDisk, "disk", false, "检查磁盘 / Check disk")
	cmd.Flags().BoolVar(&checkNet, "network", false, "检查网络 / Check network")
	cmd.Flags().BoolVar(&checkAll, "all", false, "检查所有 / Check all")
	cmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "JSON格式输出 / JSON output")
	return cmd
}

// snapshotCmd - 快照管理命令 / Snapshot management command
func snapshotCmd() *cobra.Command {
	var (
		snapshotName string
		snapshotID1  string
		snapshotID2  string
	)

	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "快照管理 / Snapshot management",
		Long:  "创建系统状态快照并支持历史对比。\nCreate system state snapshots and support historical comparison.",
	}

	// create子命令 / create subcommand
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "创建快照 / Create snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := loadConfig()

			store, err := storage.NewSQLiteStorage(c.Storage.Path)
			if err != nil {
				return fmt.Errorf("初始化存储失败: %w", err)
			}
			defer store.Close()

			sysCollector := collector.NewSystemCollector(c.Server.ID, true, true, true, true, false)
			metrics, err := sysCollector.CollectAll()
			if err != nil {
				return fmt.Errorf("采集失败: %w", err)
			}

			name := snapshotName
			if name == "" {
				name = metrics.Timestamp.Format("2006-01-02_15-04-05")
			}

			snapMgr := snapshot.NewManager(store, c.Snapshot.MaxCount)
			snap, err := snapMgr.Create(c.Server.ID, name, metrics)
			if err != nil {
				return fmt.Errorf("创建快照失败: %w", err)
			}

			fmt.Printf("快照已创建: %s (%s)\n", snap.Name, snap.ID)
			return nil
		},
	}
	createCmd.Flags().StringVarP(&snapshotName, "name", "n", "", "快照名称 / Snapshot name")

	// list子命令 / list subcommand
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "列出快照 / List snapshots",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := loadConfig()

			store, err := storage.NewSQLiteStorage(c.Storage.Path)
			if err != nil {
				return fmt.Errorf("初始化存储失败: %w", err)
			}
			defer store.Close()

			snapMgr := snapshot.NewManager(store, c.Snapshot.MaxCount)
			snaps, err := snapMgr.List(c.Server.ID)
			if err != nil {
				return err
			}

			if len(snaps) == 0 {
				fmt.Println("暂无快照 / No snapshots")
				return nil
			}

			fmt.Printf("共 %d 个快照:\n", len(snaps))
			for _, s := range snaps {
				fmt.Printf("  [%d] %s - %s (%s)\n", s.ID, s.Name, s.Timestamp.Format("2006-01-02 15:04:05"), s.ServerID)
			}
			return nil
		},
	}

	// compare子命令 / compare subcommand
	compareCmd := &cobra.Command{
		Use:   "compare",
		Short: "对比两个快照 / Compare two snapshots",
		RunE: func(cmd *cobra.Command, args []string) error {
			if snapshotID1 == "" || snapshotID2 == "" {
				return fmt.Errorf("请指定两个快照ID / Please specify two snapshot IDs")
			}
			fmt.Printf("对比快照 %s 和 %s\n", snapshotID1, snapshotID2)
			// 实际对比逻辑需要从存储加载快照数据
			// Actual comparison logic requires loading snapshot data from storage
			return nil
		},
	}
	compareCmd.Flags().StringVar(&snapshotID1, "id1", "", "第一个快照ID / First snapshot ID")
	compareCmd.Flags().StringVar(&snapshotID2, "id2", "", "第二个快照ID / Second snapshot ID")

	cmd.AddCommand(createCmd)
	cmd.AddCommand(listCmd)
	cmd.AddCommand(compareCmd)
	return cmd
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

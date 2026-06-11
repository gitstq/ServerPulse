// Package server - Hub+Agent架构 / Hub+Agent Architecture
// Hub: 接收Agent上报的数据，提供集中管理
// Agent: 采集本地数据并上报到Hub
//
// Hub: Receives data reported by Agents, provides centralized management
// Agent: Collects local data and reports to Hub
package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gitstq/ServerPulse/internal/storage"
	"github.com/gitstq/ServerPulse/pkg/config"
)

// Hub - Hub服务端 / Hub server
// 接收来自多个Agent的数据上报，存储并管理
// Receives data reports from multiple Agents, stores and manages them
type Hub struct {
	cfg        config.HubConfig
	storage    *storage.SQLiteStorage
	servers    map[string]*AgentInfo // 已注册的Agent / Registered agents
	mu         sync.RWMutex
	httpServer *http.Server
}

// AgentInfo - Agent注册信息 / Agent registration info
type AgentInfo struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Address  string    `json:"address"`
	LastSeen time.Time `json:"last_seen"`
	Status   string    `json:"status"`
}

// NewHub - 创建Hub服务 / Create Hub service
func NewHub(cfg config.HubConfig, store *storage.SQLiteStorage) *Hub {
	return &Hub{
		cfg:     cfg,
		storage: store,
		servers: make(map[string]*AgentInfo),
	}
}

// Start - 启动Hub服务 / Start Hub service
func (h *Hub) Start() error {
	mux := http.NewServeMux()

	// 注册API路由 / Register API routes
	mux.HandleFunc("/api/v1/report", h.handleReport)
	mux.HandleFunc("/api/v1/register", h.handleRegister)
	mux.HandleFunc("/api/v1/servers", h.handleServers)
	mux.HandleFunc("/api/v1/health", h.handleHealth)

	h.httpServer = &http.Server{
		Addr:         h.cfg.Bind,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// 启动HTTP监听 / Start HTTP listener
	listener, err := net.Listen("tcp", h.cfg.Bind)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", h.cfg.Bind, err)
	}

	log.Printf("[Hub] 服务启动，监听地址: %s", h.cfg.Bind)

	go func() {
		if err := h.httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Printf("[Hub] 服务错误: %v", err)
		}
	}()

	return nil
}

// Stop - 停止Hub服务 / Stop Hub service
func (h *Hub) Stop() error {
	if h.httpServer != nil {
		return h.httpServer.Close()
	}
	return nil
}

// handleHealth - 健康检查 / Health check
func (h *Hub) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"agents":    len(h.servers),
	})
}

// handleRegister - 处理Agent注册 / Handle Agent registration
func (h *Hub) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 验证API Key / Verify API key
	if h.cfg.APIKey != "" {
		if r.Header.Get("X-API-Key") != h.cfg.APIKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	var info AgentInfo
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	info.LastSeen = time.Now()
	info.Status = "online"

	h.mu.Lock()
	h.servers[info.ID] = &info
	h.mu.Unlock()

	// 存储服务器信息 / Store server info
	if h.storage != nil {
		h.storage.UpsertServer(&storage.ServerInfo{
			ID:        info.ID,
			Name:      info.Name,
			Status:    "online",
			LastSeen:  info.LastSeen,
			UpdatedAt: time.Now(),
		})
	}

	log.Printf("[Hub] Agent注册: %s (%s)", info.Name, info.ID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "registered"})
}

// handleReport - 处理Agent数据上报 / Handle Agent data report
func (h *Hub) handleReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 验证API Key / Verify API key
	if h.cfg.APIKey != "" {
		if r.Header.Get("X-API-Key") != h.cfg.APIKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	// 更新Agent最后在线时间 / Update Agent last seen time
	var report struct {
		ServerID string `json:"server_id"`
	}
	if err := json.Unmarshal(body, &report); err == nil && report.ServerID != "" {
		h.mu.Lock()
		if agent, ok := h.servers[report.ServerID]; ok {
			agent.LastSeen = time.Now()
		}
		h.mu.Unlock()
	}

	// 存储数据 / Store data (简化处理，实际应解析完整数据）
	log.Printf("[Hub] 收到数据上报: %d bytes", len(body))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// handleServers - 返回已注册的Agent列表 / Return registered Agent list
func (h *Hub) handleServers(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	servers := make([]*AgentInfo, 0, len(h.servers))
	for _, info := range h.servers {
		servers = append(servers, info)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(servers)
}

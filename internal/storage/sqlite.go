// Package storage - SQLite存储层 / SQLite Storage Layer
// 提供数据持久化功能，使用SQLite存储历史指标、探测结果和告警记录
// Provides data persistence using SQLite for historical metrics, probe results and alert records
package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStorage - SQLite存储实现 / SQLite storage implementation
type SQLiteStorage struct {
	db *sql.DB
}

// NewSQLiteStorage - 创建SQLite存储 / Create SQLite storage
// 自动创建数据库文件和表结构
// Automatically creates database file and table schema
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	// 确保目录存在 / Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// 打开数据库连接 / Open database connection
	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 测试连接 / Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	storage := &SQLiteStorage{db: db}

	// 初始化表结构 / Initialize table schema
	if err := storage.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return storage, nil
}

// initSchema - 初始化数据库表结构 / Initialize database schema
func (s *SQLiteStorage) initSchema() error {
	queries := []string{
		// 指标数据表 / Metrics data table
		`CREATE TABLE IF NOT EXISTS metrics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			server_id TEXT NOT NULL,
			timestamp DATETIME NOT NULL,
			metric TEXT NOT NULL,
			value REAL NOT NULL,
			unit TEXT DEFAULT '',
			labels TEXT DEFAULT '{}'
		)`,
		// 指标索引 / Metrics index
		`CREATE INDEX IF NOT EXISTS idx_metrics_server_metric_ts ON metrics(server_id, metric, timestamp)`,
		`CREATE INDEX IF NOT EXISTS idx_metrics_timestamp ON metrics(timestamp)`,

		// 探测结果表 / Probe results table
		`CREATE TABLE IF NOT EXISTS probe_results (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			server_id TEXT NOT NULL,
			timestamp DATETIME NOT NULL,
			target TEXT NOT NULL,
			type TEXT NOT NULL,
			status TEXT NOT NULL,
			latency REAL DEFAULT 0,
			message TEXT DEFAULT ''
		)`,
		// 探测结果索引 / Probe results index
		`CREATE INDEX IF NOT EXISTS idx_probe_server_ts ON probe_results(server_id, timestamp)`,

		// 告警记录表 / Alert records table
		`CREATE TABLE IF NOT EXISTS alerts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			server_id TEXT NOT NULL,
			timestamp DATETIME NOT NULL,
			level TEXT NOT NULL,
			type TEXT NOT NULL,
			metric TEXT NOT NULL,
			value REAL DEFAULT 0,
			threshold REAL DEFAULT 0,
			message TEXT DEFAULT '',
			acked INTEGER DEFAULT 0
		)`,
		// 告警索引 / Alerts index
		`CREATE INDEX IF NOT EXISTS idx_alerts_server_ts ON alerts(server_id, timestamp)`,

		// 服务器信息表 / Server info table
		`CREATE TABLE IF NOT EXISTS servers (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			ip TEXT DEFAULT '',
			os TEXT DEFAULT '',
			arch TEXT DEFAULT '',
			status TEXT DEFAULT 'online',
			last_seen DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		)`,

		// 快照数据表 / Snapshot data table
		`CREATE TABLE IF NOT EXISTS snapshots (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			server_id TEXT NOT NULL,
			timestamp DATETIME NOT NULL,
			name TEXT NOT NULL,
			data TEXT NOT NULL,
			created_at DATETIME NOT NULL
		)`,
		// 快照索引 / Snapshot index
		`CREATE INDEX IF NOT EXISTS idx_snapshots_server_ts ON snapshots(server_id, timestamp)`,
	}

	for _, query := range queries {
		if _, err := s.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute schema query: %w", err)
		}
	}

	return nil
}

// Close - 关闭数据库连接 / Close database connection
func (s *SQLiteStorage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// InsertMetric - 插入指标数据 / Insert metric data
func (s *SQLiteStorage) InsertMetric(m *MetricData) error {
	_, err := s.db.Exec(
		`INSERT INTO metrics (server_id, timestamp, metric, value, unit, labels) VALUES (?, ?, ?, ?, ?, ?)`,
		m.ServerID, m.Timestamp, m.Metric, m.Value, m.Unit, m.Labels,
	)
	return err
}

// InsertMetrics - 批量插入指标数据 / Batch insert metric data
func (s *SQLiteStorage) InsertMetrics(metrics []MetricData) error {
	if len(metrics) == 0 {
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO metrics (server_id, timestamp, metric, value, unit, labels) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, m := range metrics {
		if _, err := stmt.Exec(m.ServerID, m.Timestamp, m.Metric, m.Value, m.Unit, m.Labels); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// QueryMetrics - 查询指标数据 / Query metric data
func (s *SQLiteStorage) QueryMetrics(serverID, metric string, since time.Time, limit int) ([]MetricData, error) {
	if limit <= 0 {
		limit = 100
	}

	query := `SELECT id, server_id, timestamp, metric, value, unit, labels FROM metrics WHERE 1=1`
	var args []interface{}

	if serverID != "" {
		query += " AND server_id = ?"
		args = append(args, serverID)
	}
	if metric != "" {
		query += " AND metric = ?"
		args = append(args, metric)
	}
	if !since.IsZero() {
		query += " AND timestamp >= ?"
		args = append(args, since)
	}

	query += " ORDER BY timestamp DESC LIMIT ?"
	args = append(args, limit)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []MetricData
	for rows.Next() {
		var m MetricData
		if err := rows.Scan(&m.ID, &m.ServerID, &m.Timestamp, &m.Metric, &m.Value, &m.Unit, &m.Labels); err != nil {
			return nil, err
		}
		results = append(results, m)
	}

	return results, nil
}

// InsertProbeResult - 插入探测结果 / Insert probe result
func (s *SQLiteStorage) InsertProbeResult(r *ProbeResult) error {
	_, err := s.db.Exec(
		`INSERT INTO probe_results (server_id, timestamp, target, type, status, latency, message) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		r.ServerID, r.Timestamp, r.Target, r.Type, r.Status, r.Latency, r.Message,
	)
	return err
}

// QueryProbeResults - 查询探测结果 / Query probe results
func (s *SQLiteStorage) QueryProbeResults(serverID string, since time.Time, limit int) ([]ProbeResult, error) {
	if limit <= 0 {
		limit = 100
	}

	query := `SELECT id, server_id, timestamp, target, type, status, latency, message FROM probe_results WHERE 1=1`
	var args []interface{}

	if serverID != "" {
		query += " AND server_id = ?"
		args = append(args, serverID)
	}
	if !since.IsZero() {
		query += " AND timestamp >= ?"
		args = append(args, since)
	}

	query += " ORDER BY timestamp DESC LIMIT ?"
	args = append(args, limit)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ProbeResult
	for rows.Next() {
		var r ProbeResult
		if err := rows.Scan(&r.ID, &r.ServerID, &r.Timestamp, &r.Target, &r.Type, &r.Status, &r.Latency, &r.Message); err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, nil
}

// InsertAlert - 插入告警记录 / Insert alert record
func (s *SQLiteStorage) InsertAlert(a *AlertRecord) error {
	acked := 0
	if a.Acked {
		acked = 1
	}
	_, err := s.db.Exec(
		`INSERT INTO alerts (server_id, timestamp, level, type, metric, value, threshold, message, acked) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		a.ServerID, a.Timestamp, a.Level, a.Type, a.Metric, a.Value, a.Threshold, a.Message, acked,
	)
	return err
}

// QueryAlerts - 查询告警记录 / Query alert records
func (s *SQLiteStorage) QueryAlerts(serverID string, since time.Time, limit int) ([]AlertRecord, error) {
	if limit <= 0 {
		limit = 100
	}

	query := `SELECT id, server_id, timestamp, level, type, metric, value, threshold, message, acked FROM alerts WHERE 1=1`
	var args []interface{}

	if serverID != "" {
		query += " AND server_id = ?"
		args = append(args, serverID)
	}
	if !since.IsZero() {
		query += " AND timestamp >= ?"
		args = append(args, since)
	}

	query += " ORDER BY timestamp DESC LIMIT ?"
	args = append(args, limit)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []AlertRecord
	for rows.Next() {
		var a AlertRecord
		var acked int
		if err := rows.Scan(&a.ID, &a.ServerID, &a.Timestamp, &a.Level, &a.Type, &a.Metric, &a.Value, &a.Threshold, &a.Message, &acked); err != nil {
			return nil, err
		}
		a.Acked = acked == 1
		results = append(results, a)
	}

	return results, nil
}

// UpsertServer - 更新或插入服务器信息 / Upsert server info
func (s *SQLiteStorage) UpsertServer(info *ServerInfo) error {
	_, err := s.db.Exec(
		`INSERT OR REPLACE INTO servers (id, name, ip, os, arch, status, last_seen, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		info.ID, info.Name, info.IP, info.OS, info.Arch, info.Status, info.LastSeen, info.UpdatedAt,
	)
	return err
}

// GetServers - 获取所有服务器信息 / Get all server info
func (s *SQLiteStorage) GetServers() ([]ServerInfo, error) {
	rows, err := s.db.Query(`SELECT id, name, ip, os, arch, status, last_seen, updated_at FROM servers ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ServerInfo
	for rows.Next() {
		var info ServerInfo
		if err := rows.Scan(&info.ID, &info.Name, &info.IP, &info.OS, &info.Arch, &info.Status, &info.LastSeen, &info.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, info)
	}

	return results, nil
}

// InsertSnapshot - 插入快照 / Insert snapshot
func (s *SQLiteStorage) InsertSnapshot(snap *SnapshotData) error {
	_, err := s.db.Exec(
		`INSERT INTO snapshots (server_id, timestamp, name, data, created_at) VALUES (?, ?, ?, ?, ?)`,
		snap.ServerID, snap.Timestamp, snap.Name, snap.Data, snap.CreatedAt,
	)
	return err
}

// GetSnapshots - 获取快照列表 / Get snapshot list
func (s *SQLiteStorage) GetSnapshots(serverID string, limit int) ([]SnapshotData, error) {
	if limit <= 0 {
		limit = 50
	}

	query := `SELECT id, server_id, timestamp, name, data, created_at FROM snapshots WHERE 1=1`
	var args []interface{}

	if serverID != "" {
		query += " AND server_id = ?"
		args = append(args, serverID)
	}

	query += " ORDER BY timestamp DESC LIMIT ?"
	args = append(args, limit)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []SnapshotData
	for rows.Next() {
		var snap SnapshotData
		if err := rows.Scan(&snap.ID, &snap.ServerID, &snap.Timestamp, &snap.Name, &snap.Data, &snap.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, snap)
	}

	return results, nil
}

// CleanupOldData - 清理过期数据 / Cleanup expired data
func (s *SQLiteStorage) CleanupOldData(retentionDays int) error {
	if retentionDays <= 0 {
		return nil
	}

	cutoff := time.Now().AddDate(0, 0, -retentionDays)

	// 清理指标数据 / Cleanup metrics
	if _, err := s.db.Exec(`DELETE FROM metrics WHERE timestamp < ?`, cutoff); err != nil {
		return fmt.Errorf("failed to cleanup metrics: %w", err)
	}

	// 清理探测结果 / Cleanup probe results
	if _, err := s.db.Exec(`DELETE FROM probe_results WHERE timestamp < ?`, cutoff); err != nil {
		return fmt.Errorf("failed to cleanup probe results: %w", err)
	}

	// 清理告警记录 / Cleanup alert records
	if _, err := s.db.Exec(`DELETE FROM alerts WHERE timestamp < ?`, cutoff); err != nil {
		return fmt.Errorf("failed to cleanup alerts: %w", err)
	}

	return nil
}

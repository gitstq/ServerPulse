// Package anomaly - 异常检测引擎 / Anomaly Detection Engine
// 基于Z-Score和移动平均的统计异常检测
// Statistical anomaly detection based on Z-Score and moving average
package anomaly

import (
	"fmt"
	"sync"
	"time"
)

// AnomalyResult - 异常检测结果 / Anomaly detection result
type AnomalyResult struct {
	Metric    string    `json:"metric"`     // 指标名称 / Metric name
	Value     float64   `json:"value"`      // 当前值 / Current value
	ZScore    float64   `json:"z_score"`    // Z-Score值 / Z-Score value
	Mean      float64   `json:"mean"`       // 窗口均值 / Window mean
	StdDev    float64   `json:"std_dev"`    // 窗口标准差 / Window standard deviation
	IsAnomaly bool      `json:"is_anomaly"` // 是否异常 / Whether anomaly
	Severity  string    `json:"severity"`   // 严重程度: low, medium, high / Severity level
	Timestamp time.Time `json:"timestamp"`  // 检测时间 / Detection time
	Message   string    `json:"message"`    // 描述信息 / Description message
}

// Detector - 异常检测器 / Anomaly detector
type Detector struct {
	mu             sync.RWMutex
	models         map[string]*MovingAverage // 每个指标的移动平均模型 / Moving average model per metric
	windowSize     int                      // 滑动窗口大小 / Sliding window size
	zScoreThreshold float64                  // Z-Score阈值 / Z-Score threshold
}

// NewDetector - 创建异常检测器 / Create anomaly detector
func NewDetector(windowSize int, zScoreThreshold float64) *Detector {
	if windowSize <= 0 {
		windowSize = 20
	}
	if zScoreThreshold <= 0 {
		zScoreThreshold = 2.5
	}
	return &Detector{
		models:         make(map[string]*MovingAverage),
		windowSize:     windowSize,
		zScoreThreshold: zScoreThreshold,
	}
}

// PushData 推送数据点到检测器 / Push data point to detector
func (d *Detector) PushData(metric string, value float64) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.models[metric]; !ok {
		d.models[metric] = NewMovingAverage(d.windowSize)
	}
	d.models[metric].Push(value, time.Now())
}

// Detect 对指定指标执行异常检测 / Execute anomaly detection on specified metric
// 使用Z-Score算法：
//   Z-Score = (当前值 - 均值) / 标准差
//   当 |Z-Score| > 阈值时判定为异常
//
// Uses Z-Score algorithm:
//   Z-Score = (current value - mean) / standard deviation
//   Anomaly detected when |Z-Score| > threshold
func (d *Detector) Detect(metric string, value float64) *AnomalyResult {
	d.mu.Lock()
	defer d.mu.Unlock()

	result := &AnomalyResult{
		Metric:    metric,
		Value:     value,
		Timestamp: time.Now(),
	}

	model, ok := d.models[metric]
	if !ok {
		// 首次数据，创建模型 / First data point, create model
		d.models[metric] = NewMovingAverage(d.windowSize)
		d.models[metric].Push(value, time.Now())
		result.Message = "insufficient data for detection"
		return result
	}

	// 推入新数据 / Push new data
	model.Push(value, time.Now())

	// 至少需要3个数据点才能计算标准差 / Need at least 3 data points for stddev
	if model.Len() < 3 {
		result.Message = fmt.Sprintf("collecting data (%d/%d)", model.Len(), d.windowSize)
		return result
	}

	mean := model.Mean()
	stdDev := model.StdDev()

	result.Mean = mean
	result.StdDev = stdDev

	// 计算Z-Score / Calculate Z-Score
	if stdDev > 0 {
		result.ZScore = (value - mean) / stdDev
	} else {
		result.ZScore = 0
		result.Message = "standard deviation is zero, cannot compute Z-Score"
		return result
	}

	// 判断是否异常 / Determine if anomaly
	absZScore := result.ZScore
	if absZScore < 0 {
		absZScore = -absZScore
	}

	if absZScore > d.zScoreThreshold {
		result.IsAnomaly = true
		// 根据Z-Score大小确定严重程度 / Determine severity based on Z-Score magnitude
		switch {
		case absZScore > 4.0:
			result.Severity = "high"
			result.Message = fmt.Sprintf("严重异常: Z-Score=%.2f (阈值=%.1f), 值=%.2f, 均值=%.2f",
				result.ZScore, d.zScoreThreshold, value, mean)
		case absZScore > 3.0:
			result.Severity = "medium"
			result.Message = fmt.Sprintf("中度异常: Z-Score=%.2f (阈值=%.1f), 值=%.2f, 均值=%.2f",
				result.ZScore, d.zScoreThreshold, value, mean)
		default:
			result.Severity = "low"
			result.Message = fmt.Sprintf("轻微异常: Z-Score=%.2f (阈值=%.1f), 值=%.2f, 均值=%.2f",
				result.ZScore, d.zScoreThreshold, value, mean)
		}
	} else {
		result.IsAnomaly = false
		result.Severity = "normal"
		result.Message = fmt.Sprintf("正常: Z-Score=%.2f, 值=%.2f, 均值=%.2f",
			result.ZScore, value, mean)
	}

	return result
}

// GetModelStats 获取指定指标的统计信息 / Get statistics for specified metric
func (d *Detector) GetModelStats(metric string) (mean, stdDev float64, count int) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	model, ok := d.models[metric]
	if !ok {
		return 0, 0, 0
	}
	return model.Mean(), model.StdDev(), model.Len()
}

// Reset 重置指定指标的模型 / Reset model for specified metric
func (d *Detector) Reset(metric string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.models, metric)
}

// ResetAll 重置所有模型 / Reset all models
func (d *Detector) ResetAll() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.models = make(map[string]*MovingAverage)
}

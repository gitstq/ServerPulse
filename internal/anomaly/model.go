// Package anomaly - 统计模型定义 / Statistical Model Definitions
// 定义异常检测使用的统计模型和参数
// Defines statistical models and parameters for anomaly detection
package anomaly

import "time"

// DataPoint - 数据点 / Data point
type DataPoint struct {
	Value     float64   `json:"value"`     // 指标值 / Metric value
	Timestamp time.Time `json:"timestamp"` // 时间戳 / Timestamp
}

// MovingAverage - 移动平均模型 / Moving average model
type MovingAverage struct {
	Window []DataPoint // 滑动窗口数据 / Sliding window data
	Size   int         // 窗口大小 / Window size
}

// NewMovingAverage - 创建移动平均模型 / Create moving average model
func NewMovingAverage(windowSize int) *MovingAverage {
	if windowSize <= 0 {
		windowSize = 20
	}
	return &MovingAverage{
		Window: make([]DataPoint, 0, windowSize),
		Size:   windowSize,
	}
}

// Push 添加数据点到窗口 / Add data point to window
func (ma *MovingAverage) Push(value float64, ts time.Time) {
	ma.Window = append(ma.Window, DataPoint{Value: value, Timestamp: ts})
	// 保持窗口大小 / Maintain window size
	if len(ma.Window) > ma.Size {
		ma.Window = ma.Window[len(ma.Window)-ma.Size:]
	}
}

// Mean 计算窗口均值 / Calculate window mean
func (ma *MovingAverage) Mean() float64 {
	if len(ma.Window) == 0 {
		return 0
	}
	sum := 0.0
	for _, dp := range ma.Window {
		sum += dp.Value
	}
	return sum / float64(len(ma.Window))
}

// StdDev 计算窗口标准差 / Calculate window standard deviation
func (ma *MovingAverage) StdDev() float64 {
	if len(ma.Window) < 2 {
		return 0
	}
	mean := ma.Mean()
	variance := 0.0
	for _, dp := range ma.Window {
		diff := dp.Value - mean
		variance += diff * diff
	}
	variance /= float64(len(ma.Window) - 1)
	return sqrt(variance)
}

// Len 返回窗口中数据点数量 / Return number of data points in window
func (ma *MovingAverage) Len() int {
	return len(ma.Window)
}

// Values 返回所有值 / Return all values
func (ma *MovingAverage) Values() []float64 {
	values := make([]float64, len(ma.Window))
	for i, dp := range ma.Window {
		values[i] = dp.Value
	}
	return values
}

// sqrt 计算平方根 / Calculate square root
func sqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	// 牛顿迭代法 / Newton's method
	z := x
	for i := 0; i < 100; i++ {
		z = (z + x/z) / 2
	}
	return z
}

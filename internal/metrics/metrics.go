package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

// Metrics collects application metrics.
type Metrics struct {
	mu            sync.RWMutex
	counters      map[string]*atomic.Uint64
	gauges        map[string]*atomic.Int64
	histograms    map[string]*Histogram
	startTime     time.Time
}

// New creates a new metrics collector.
func New() *Metrics {
	return &Metrics{
		counters:   make(map[string]*atomic.Uint64),
		gauges:     make(map[string]*atomic.Int64),
		histograms: make(map[string]*Histogram),
		startTime:  time.Now(),
	}
}

// Counter increments a counter.
func (m *Metrics) Counter(name string) {
	m.mu.Lock()
	counter, exists := m.counters[name]
	if !exists {
		counter = &atomic.Uint64{}
		m.counters[name] = counter
	}
	m.mu.Unlock()
	counter.Add(1)
}

// SetGauge sets a gauge value.
func (m *Metrics) SetGauge(name string, value int64) {
	m.mu.Lock()
	gauge, exists := m.gauges[name]
	if !exists {
		gauge = &atomic.Int64{}
		m.gauges[name] = gauge
	}
	m.mu.Unlock()
	gauge.Store(value)
}

// Histogram records a value in a histogram.
func (m *Metrics) Histogram(name string, value float64) {
	m.mu.Lock()
	hist, exists := m.histograms[name]
	if !exists {
		hist = NewHistogram()
		m.histograms[name] = hist
	}
	m.mu.Unlock()
	hist.Record(value)
}

// GetCounter gets a counter value.
func (m *Metrics) GetCounter(name string) uint64 {
	m.mu.RLock()
	counter, exists := m.counters[name]
	m.mu.RUnlock()
	if !exists {
		return 0
	}
	return counter.Load()
}

// GetGauge gets a gauge value.
func (m *Metrics) GetGauge(name string) int64 {
	m.mu.RLock()
	gauge, exists := m.gauges[name]
	m.mu.RUnlock()
	if !exists {
		return 0
	}
	return gauge.Load()
}

// GetHistogram gets a histogram.
func (m *Metrics) GetHistogram(name string) *Histogram {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.histograms[name]
}

// Uptime returns the uptime since metrics started.
func (m *Metrics) Uptime() time.Duration {
	return time.Since(m.startTime)
}

// Histogram tracks value distribution.
type Histogram struct {
	mu     sync.RWMutex
	values []float64
	count  uint64
	sum    float64
	min    float64
	max    float64
}

// NewHistogram creates a new histogram.
func NewHistogram() *Histogram {
	return &Histogram{
		values: make([]float64, 0),
		min:    0,
		max:    0,
	}
}

// Record records a value.
func (h *Histogram) Record(value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.values = append(h.values, value)
	h.count++
	h.sum += value

	if h.count == 1 || value < h.min {
		h.min = value
	}
	if h.count == 1 || value > h.max {
		h.max = value
	}
}

// Mean returns the mean value.
func (h *Histogram) Mean() float64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.count == 0 {
		return 0
	}
	return h.sum / float64(h.count)
}

// Min returns the minimum value.
func (h *Histogram) Min() float64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.min
}

// Max returns the maximum value.
func (h *Histogram) Max() float64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.max
}

// Count returns the number of recorded values.
func (h *Histogram) Count() uint64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.count
}

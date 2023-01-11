package ginmetrics

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/metrics/influxdb"
)

type MetricType int

const (
	defaultSlowTime = 200
)

var (
	monitor *Monitor
	once    sync.Once
)

// Monitor is an object that uses to set gin server monitor.
type Monitor struct {
	slowTime int
	registry metrics.Registry
}

// GetMonitor used to get global Monitor object,
// this function returns a singleton object.
func GetMonitor() *Monitor {
	if monitor == nil {
		once.Do(func() {
			metrics.Enabled = true
			monitor = &Monitor{
				slowTime: defaultSlowTime,
				registry: metrics.NewRegistry(),
			}
		})
	}
	return monitor
}

// SetSlowTime set slowTime property. slowTime uint is milli second. slowTime is used to determine whether
// the request is slow. For "gin_slow_request_total" metric.
func (m *Monitor) SetSlowTime(slowTime int) {
	m.slowTime = slowTime
}

func (m *Monitor) SetRegistry(registry metrics.Registry) {
	m.registry = registry
}

func (m *Monitor) ReportToInfluxDB(d time.Duration, url, database, username, password, namespace string) {
	go influxdb.InfluxDB(m.registry, d, url, database, username, password, namespace)
}

// m.ReportToLogger(metrics.DefaultRegistry, 10*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
func (m *Monitor) ReportToLogger(freq time.Duration, l metrics.Logger) {
	go metrics.Log(m.registry, freq, l)
}

func (m *Monitor) GetOrRegisterCounter(name string) metrics.Counter {
	return metrics.GetOrRegisterCounter(name, m.registry)
}

func (m *Monitor) GetOrRegisterGauge(name string) metrics.Gauge {
	return metrics.GetOrRegisterGauge(name, m.registry)
}

func (m *Monitor) GetOrRegisterGaugeFloat64(name string) metrics.GaugeFloat64 {
	return metrics.GetOrRegisterGaugeFloat64(name, m.registry)
}

func (m *Monitor) GetOrRegisterMeter(name string) metrics.Meter {
	return metrics.GetOrRegisterMeter(name, m.registry)
}

func (m *Monitor) GetOrRegisterHistogram(name string) metrics.Histogram {
	return metrics.GetOrRegisterHistogram(name, m.registry, metrics.NewExpDecaySample(1028, 0.015))
}

func (m *Monitor) GetOrRegisterTimer(name string) metrics.Timer {
	return metrics.GetOrRegisterTimer(name, m.registry)
}

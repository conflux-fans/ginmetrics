package ginmetrics

import (
	"testing"

	"github.com/ethereum/go-ethereum/metrics"
	"gotest.tools/assert"
)

func TestCounterMetrics(t *testing.T) {
	metrics.Enable()
	r := metrics.NewRegistry()
	metrics.GetOrRegisterCounter("counte_metrics", r).Inc(10)
	assert.Equal(t, metrics.GetOrRegisterCounter("counte_metrics", r).Snapshot().Count(), int64(10))
}

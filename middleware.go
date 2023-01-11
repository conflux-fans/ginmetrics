package ginmetrics

import (
	"net/http"
	"time"

	"github.com/conflux-fans/ginmetrics/bloom"
	"github.com/gin-gonic/gin"
)

var (
	metricRequestTotal     = "gin_request_total"       //"all the server received request num."
	metricRequestFailTotal = "gin_request_fail_total"  //"all the server received request num."
	metricRequestUVTotal   = "gin_request_uv_total"    //"all the server received ip num."
	metricRequestBody      = "gin_request_body_total"  //"the server received request body size, unit byte"
	metricResponseBody     = "gin_response_body_total" //"the server send response body size, unit byte"
	metricRequestDuration  = "gin_request_duration"    //"the time server took to handle the request."
	metricSlowRequest      = "gin_slow_request_total"  //"the server handled slow requests counter"

	bloomFilter *bloom.BloomFilter = bloom.NewBloomFilter()
)

// Use set gin metrics middleware
func (m *Monitor) Use(r gin.IRoutes) {
	r.Use(m.monitorInterceptor)
}

// monitorInterceptor as gin monitor middleware.
func (m *Monitor) monitorInterceptor(ctx *gin.Context) {
	startTime := time.Now()

	// execute normal process.
	ctx.Next()

	// after request
	m.ginMetricHandle(ctx, startTime)
}

func (m *Monitor) ginMetricHandle(ctx *gin.Context, start time.Time) {
	r := ctx.Request
	w := ctx.Writer

	// set uv
	if clientIP := ctx.ClientIP(); !bloomFilter.Contains(clientIP) {
		bloomFilter.Add(clientIP)
		m.GetOrRegisterCounter(metricRequestUVTotal).Inc(1)
	}

	// set request total
	m.GetOrRegisterCounter(metricRequestTotal + "/all").Inc(1)
	m.GetOrRegisterCounter(metricRequestTotal + "/" + ctx.Request.Method + ctx.FullPath()).Inc(1)

	if w.Status() >= http.StatusBadRequest {
		m.GetOrRegisterCounter(metricRequestFailTotal + "/all").Inc(1)
		m.GetOrRegisterCounter(metricRequestFailTotal + "/" + ctx.Request.Method + ctx.FullPath()).Inc(1)
	}

	// set request body size
	// since r.ContentLength can be negative (in some occasions) guard the operation
	if r.ContentLength >= 0 {
		m.GetOrRegisterCounter(metricRequestBody).Inc(r.ContentLength)
	}

	// set slow request
	latency := time.Since(start)
	// fmt.Printf("latency %v\n", latency)
	if (int(latency.Milliseconds())) > m.slowTime {
		m.GetOrRegisterCounter(metricSlowRequest).Inc(1)
	}

	// set request duration
	m.GetOrRegisterTimer(metricRequestDuration + "/all").Update(latency)
	m.GetOrRegisterTimer(metricRequestDuration + "/" + ctx.Request.Method + ctx.FullPath()).Update(latency)

	// m.GetOrRegisterGauge(metricRequestDuration + "/" + ctx.Request.Method + ctx.FullPath()).Update(latency.Nanoseconds())

	// set response size
	if w.Size() > 0 {
		m.GetOrRegisterCounter(metricResponseBody).Inc(int64(w.Size()))
	}
}

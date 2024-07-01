package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metric struct {
	requestDuration *prometheus.HistogramVec
	requestCount    *prometheus.CounterVec
}

func New() *Metric {
	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "request_duration_seconds",
		Help:    "measure request duration in seconds, segmented by method, path, and status code.",
		Buckets: []float64{.005, .025, .5, 2},
	}, []string{"method", "path", "code"})

	requestCount := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "request_count",
		Help: "count of requests, segmented by method, path, and status code.",
	}, []string{"method", "path", "code"})

	prometheus.MustRegister(
		requestDuration,
		requestCount,
	)

	return &Metric{
		requestDuration: requestDuration,
		requestCount:    requestCount,
	}
}

func (m *Metric) ObserveRequestDuration(method string, route string, code int, elapsed time.Duration) {
	m.requestDuration.WithLabelValues(method, route, strconv.Itoa(code)).Observe(elapsed.Seconds())
}

func (m *Metric) IncRequestCount(method string, route string, code int) {
	m.requestCount.WithLabelValues(method, route, strconv.Itoa(code)).Inc()
}

// @Tags SRE
// @Summary expose prometheus metric
// @Id metrics
// @Router /metrics [get]
// @version 1.0
// @Success 200
func (m *Metric) Handler() http.Handler {
	return promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
		Registry:           prometheus.DefaultRegisterer,
		DisableCompression: true,
	})
}

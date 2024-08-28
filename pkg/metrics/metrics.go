package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// metricNamespace is metrics namespace.
const metricNamespace = "moco_proxy"

// MustRegisterGauge function returns new gauge vec.
func MustRegisterGauge(subsystem, name, help string, labelNames ...string) *prometheus.GaugeVec {
	return promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricNamespace,
		Subsystem: subsystem,
		Name:      name,
		Help:      help,
	}, labelNames)
}

// MustRegisterCounter function returns new counter vec.
func MustRegisterCounter(
	subsystem, name, help string,
	labelNames ...string,
) *prometheus.CounterVec {
	return promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: metricNamespace,
		Subsystem: subsystem,
		Name:      name,
		Help:      help,
	}, labelNames)
}

// MustRegisterHistogram function returns new histogram vec.
func MustRegisterHistogram(
	subsystem, name, help string,
	buckets []float64,
	labelNames ...string,
) *prometheus.HistogramVec {
	return promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: metricNamespace,
		Subsystem: subsystem,
		Name:      name,
		Help:      help,
		Buckets:   buckets,
	}, labelNames)
}

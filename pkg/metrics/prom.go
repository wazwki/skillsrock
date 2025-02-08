package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ObserveRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"handler"},
	)
)

func init() {
	prometheus.MustRegister(ObserveRequestDuration)
}

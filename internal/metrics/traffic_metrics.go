package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type TrafficMetrics struct {
	trafficCounter *prometheus.CounterVec
}

func NewTrafficMetrics() (*TrafficMetrics, error) {
	result := new(TrafficMetrics)
	result.trafficCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of requests.",
		},
		[]string{"path"},
	)
	prometheus.Register(result.trafficCounter)
	return result, nil
}

func (tm *TrafficMetrics) RequestHandled(urlPath string) {
	tm.trafficCounter.WithLabelValues(urlPath).Inc()
}

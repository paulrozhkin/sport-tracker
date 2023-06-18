package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
)

type TrafficMetrics struct {
	trafficCounter *prometheus.CounterVec
	httpDuration   *prometheus.HistogramVec
}

type TrafficTimer struct {
	timer *prometheus.Timer
}

func NewTrafficMetrics() (*TrafficMetrics, error) {
	result := new(TrafficMetrics)
	result.trafficCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of requests.",
		},
		[]string{"path", "status"},
	)

	result.httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_response_time_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})

	return result, nil
}

func (tm *TrafficMetrics) RequestHandled(urlPath string, status int) {
	tm.trafficCounter.WithLabelValues(urlPath, strconv.Itoa(status)).Inc()
}

func (tm *TrafficMetrics) GetTimer(urlPath string) *TrafficTimer {
	timer := prometheus.NewTimer(tm.httpDuration.WithLabelValues(urlPath))
	return &TrafficTimer{timer: timer}
}

func (timer *TrafficTimer) ObserveDuration() {
	timer.timer.ObserveDuration()
}

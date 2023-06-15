package http_server

import (
	"github.com/paulrozhkin/sport-tracker/internal/metrics"
	"net/http"
)

type TrafficMiddleware struct {
	trafficMetrics *metrics.TrafficMetrics
}

func NewTrafficMiddleware(trafficMetrics *metrics.TrafficMetrics) (*TrafficMiddleware, error) {
	return &TrafficMiddleware{trafficMetrics: trafficMetrics}, nil
}

func (tm *TrafficMiddleware) CalculateTraffic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		tm.trafficMetrics.RequestHandled(r.RequestURI)
	})
}

package http_server

import (
	"github.com/paulrozhkin/sport-tracker/internal/metrics"
	"net/http"
	"regexp"
)

type TrafficMiddleware struct {
	trafficMetrics *metrics.TrafficMetrics
	regex          *regexp.Regexp
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewTrafficMiddleware(trafficMetrics *metrics.TrafficMetrics) (*TrafficMiddleware, error) {
	re := regexp.MustCompile(`[\da-f]{8}-[\da-f]{4}-[\da-f]{4}-[\da-f]{4}-[\da-f]{12}`)
	return &TrafficMiddleware{trafficMetrics: trafficMetrics, regex: re}, nil
}

func (tm *TrafficMiddleware) CalculateTraffic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := tm.regex.ReplaceAllString(r.RequestURI, `{id}`)
		if uri == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}
		timer := tm.trafficMetrics.GetTimer(uri)
		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)
		timer.ObserveDuration()
		tm.trafficMetrics.RequestHandled(uri, rw.statusCode)
	})
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

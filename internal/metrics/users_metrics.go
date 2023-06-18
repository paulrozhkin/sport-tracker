package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type UsersMetrics struct {
	userCounter prometheus.Counter
}

func NewUsersMetrics() (*UsersMetrics, error) {
	result := new(UsersMetrics)
	result.userCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "registered_users_total",
			Help: "Number of registered users.",
		},
	)
	return result, nil
}

func (um *UsersMetrics) UserRegistered() {
	um.userCounter.Inc()
}

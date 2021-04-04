package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Durations API Duration histogram of all the apis by path
	Durations = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "notifica",
		Subsystem: "apis",
		Name:      "duration",
		Help:      "The durations of the individual api calls",
	}, []string{"path"})

	// UserLoadByAPI Load at any given moment by api path
	UserLoadByAPI = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "notifica",
		Subsystem: "apis",
		Name:      "api_load",
		Help:      "Number of active calls by api",
	}, []string{"path"})

	// UserLoadOverall Overall load for the sum of apis being handled
	UserLoadOverall = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "notifica",
		Subsystem: "apis",
		Name:      "user_load",
		Help:      "Number of active users",
	})

	// StatusErrors any non server error by path and status code
	StatusErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "notifica",
		Subsystem: "apis",
		Name:      "status_code",
		Help:      "status codes per api",
	}, []string{"path", "status"})

	// ServerErrors any server errors by path
	ServerErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "notifica",
		Subsystem: "apis",
		Name:      "errors",
		Help:      "The error count per api",
	}, []string{"path"})
)

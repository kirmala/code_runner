package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	HTTPRequestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "code_runner",
			Subsystem: "http_server",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests received.",
		},
	)
)

func Register(reg prometheus.Registerer) {
	reg.MustRegister(HTTPRequestsTotal)
	reg.MustRegister(collectors.NewGoCollector())
}

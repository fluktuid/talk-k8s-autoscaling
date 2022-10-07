package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	callsProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "server_processed_calls",
		Help: "The total number of processed calls",
	}, []string{"client", "host"})
)

func ServerRecordRequest(requestHost, targetHost string) {
	go func() {
		callsProcessed.WithLabelValues(requestHost, targetHost).Inc()
	}()
}

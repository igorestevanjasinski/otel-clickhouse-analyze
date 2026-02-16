package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	MessagesConsumed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "messages_consumed_total",
			Help: "Total messages consumed from Kafka",
		},
		[]string{"topic", "status"},
	)

	ProcessingDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "message_processing_duration_seconds",
			Help:    "Message processing duration in seconds",
			Buckets: []float64{.005, .01, .025, .05, .075, .1, .25, .5, .75, 1.0, 2.5, 5.0, 7.5, 10.0},
		},
		[]string{"topic"},
	)

	InFlightMessages = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "messages_in_flight",
			Help: "Messages currently being processed",
		},
	)

	DatabaseOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_operations_total",
			Help: "Total database operations",
		},
		[]string{"operation", "status"},
	)

	DatabaseOperationDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_operation_duration_seconds",
			Help:    "Database operation duration in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1.0},
		},
		[]string{"operation"},
	)
)

func Init() {
	prometheus.MustRegister(MessagesConsumed)
	prometheus.MustRegister(ProcessingDuration)
	prometheus.MustRegister(InFlightMessages)
	prometheus.MustRegister(DatabaseOperations)
	prometheus.MustRegister(DatabaseOperationDuration)
}

// GetHandler retorna o handler HTTP para o endpoint /metrics
// Para ser usado em um servidor HTTP compartilhado
func GetHandler() http.Handler {
	return promhttp.Handler()
}

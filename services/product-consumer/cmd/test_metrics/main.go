package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	messagesConsumed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "messages_consumed_total",
			Help: "Total messages consumed from Kafka",
		},
		[]string{"topic", "status"},
	)

	processingDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "message_processing_duration_seconds",
			Help:    "Message processing duration in seconds",
			Buckets: []float64{.005, .01, .025, .05, .075, .1, .25, .5, .75, 1.0, 2.5, 5.0, 7.5, 10.0},
		},
		[]string{"topic"},
	)

	inFlightMessages = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "messages_in_flight",
			Help: "Messages currently being processed",
		},
	)

	databaseOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_operations_total",
			Help: "Total database operations",
		},
		[]string{"operation", "status"},
	)
)

func init() {
	prometheus.MustRegister(messagesConsumed)
	prometheus.MustRegister(processingDuration)
	prometheus.MustRegister(inFlightMessages)
	prometheus.MustRegister(databaseOperations)
}

func main() {
	log.Println("Iniciando servidor de métricas Prometheus na porta 8081...")

	// Simular algumas métricas para teste
	go func() {
		for {
			messagesConsumed.WithLabelValues("products.events", "success").Inc()
			processingDuration.WithLabelValues("products.events").Observe(0.05)
			inFlightMessages.Set(5)
			databaseOperations.WithLabelValues("insert", "success").Inc()
			time.Sleep(2 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("✓ Servidor rodando em http://localhost:8081")
	log.Println("  - Métricas: http://localhost:8081/metrics")
	log.Println("  - Health: http://localhost:8081/health")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/igor-jasinski/product-consumer/pkg/metrics"
)

func main() {
	metrics.Init()
	metrics.StartServer("8081")
	log.Println("Prometheus metrics server started on :8081")
	log.Println("Test metrics at: http://localhost:8081/metrics")

	// Simular algumas métricas
	go func() {
		for {
			metrics.MessagesConsumed.WithLabelValues("products.events", "success").Inc()
			metrics.InFlightMessages.Set(5)
			metrics.ProcessingDuration.WithLabelValues("products.events").Observe(0.123)
			metrics.DatabaseOperations.WithLabelValues("insert", "success").Inc()
			metrics.DatabaseOperationDuration.WithLabelValues("insert").Observe(0.015)
			time.Sleep(2 * time.Second)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
}

package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/igor-jasinski/product-consumer/pkg/metrics"
)

func main() {
	metrics.Init()

	// Criar servidor HTTP com /metrics e /health
	mux := http.NewServeMux()
	mux.Handle("/metrics", metrics.GetHandler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	go func() {
		log.Println("HTTP server started on :8081")
		log.Println("- Metrics: http://localhost:8081/metrics")
		log.Println("- Health:  http://localhost:8081/health")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

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

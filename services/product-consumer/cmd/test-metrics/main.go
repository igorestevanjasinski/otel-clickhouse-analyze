package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/igor-jasinski/product-consumer/pkg/metrics"
)

func main() {
	ctx := context.Background()
	endpoint := os.Getenv("OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4317"
	}

	shutdown, err := metrics.SetupMetrics(ctx, metrics.MetricsConfig{
		ServiceName:    "product-consumer",
		ServiceVersion: "1.0.0",
		Environment:    "development",
		OTLPEndpoint:   endpoint,
	})
	if err != nil {
		log.Fatalf("SetupMetrics: %v", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdown(shutdownCtx)
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{Addr: ":8081", Handler: mux}
	go func() {
		log.Println("HTTP server on :8081 (health: /health); metrics exported via OTLP to", endpoint)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Simulate metrics (OTLP export)
	go func() {
		for {
			metrics.RecordMessagesConsumed("products.events", "success", 1)
			metrics.InFlightAdd(5)
			metrics.InFlightAdd(-5)
			metrics.RecordProcessingDuration("products.events", 0.123)
			metrics.RecordDatabaseOperations("insert", "success", 1)
			metrics.RecordDatabaseOperationDuration("insert", 0.015)
			time.Sleep(2 * time.Second)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down...")
}

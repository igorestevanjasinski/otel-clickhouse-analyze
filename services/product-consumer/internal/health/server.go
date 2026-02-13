package health

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/igor-jasinski/product-consumer/internal/repository"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type HealthServer struct {
	server  *http.Server
	repo    repository.ProductRepository
	brokers []string
	logger  *logrus.Logger
}

type HealthResponse struct {
	Status string `json:"status"`
}

type ReadinessResponse struct {
	Status   string `json:"status"`
	Postgres string `json:"postgres"`
	Kafka    string `json:"kafka"`
}

func NewHealthServer(port int, repo repository.ProductRepository, brokers []string, logger *logrus.Logger) *HealthServer {
	hs := &HealthServer{
		repo:    repo,
		brokers: brokers,
		logger:  logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", hs.healthHandler)
	mux.HandleFunc("/ready", hs.readyHandler)

	hs.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return hs
}

func (hs *HealthServer) Start() error {
	hs.logger.WithField("port", hs.server.Addr).Info("Starting health check server")
	if err := hs.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (hs *HealthServer) Shutdown(ctx context.Context) error {
	return hs.server.Shutdown(ctx)
}

func (hs *HealthServer) healthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{Status: "healthy"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (hs *HealthServer) readyHandler(w http.ResponseWriter, r *http.Request) {
	response := ReadinessResponse{
		Status:   "ready",
		Postgres: "unknown",
		Kafka:    "unknown",
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if err := hs.checkPostgres(ctx); err != nil {
		response.Status = "not ready"
		response.Postgres = fmt.Sprintf("error: %v", err)
		hs.respondWithError(w, http.StatusServiceUnavailable, response)
		return
	}
	response.Postgres = "connected"

	if err := hs.checkKafka(ctx); err != nil {
		response.Status = "not ready"
		response.Kafka = fmt.Sprintf("error: %v", err)
		hs.respondWithError(w, http.StatusServiceUnavailable, response)
		return
	}
	response.Kafka = "connected"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (hs *HealthServer) checkPostgres(ctx context.Context) error {
	stats := hs.repo.Stats()
	if stats.TotalConns() == 0 {
		return fmt.Errorf("no database connections")
	}
	return nil
}

func (hs *HealthServer) checkKafka(ctx context.Context) error {
	conn, err := kafka.DialContext(ctx, "tcp", hs.brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Brokers()
	return err
}

func (hs *HealthServer) respondWithError(w http.ResponseWriter, code int, response ReadinessResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

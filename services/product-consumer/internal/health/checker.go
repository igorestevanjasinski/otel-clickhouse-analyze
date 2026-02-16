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

type HealthChecker struct {
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

func NewHealthChecker(repo repository.ProductRepository, brokers []string, logger *logrus.Logger) *HealthChecker {
	return &HealthChecker{
		repo:    repo,
		brokers: brokers,
		logger:  logger,
	}
}

func (hc *HealthChecker) HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{Status: "healthy"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (hc *HealthChecker) ReadyHandler(w http.ResponseWriter, r *http.Request) {
	response := ReadinessResponse{
		Status:   "ready",
		Postgres: "unknown",
		Kafka:    "unknown",
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if err := hc.checkPostgres(ctx); err != nil {
		response.Status = "not ready"
		response.Postgres = fmt.Sprintf("error: %v", err)
		hc.respondWithError(w, http.StatusServiceUnavailable, response)
		return
	}
	response.Postgres = "connected"

	if err := hc.checkKafka(ctx); err != nil {
		response.Status = "not ready"
		response.Kafka = fmt.Sprintf("error: %v", err)
		hc.respondWithError(w, http.StatusServiceUnavailable, response)
		return
	}
	response.Kafka = "connected"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (hc *HealthChecker) checkPostgres(ctx context.Context) error {
	stats := hc.repo.Stats()
	if stats.TotalConns() == 0 {
		return fmt.Errorf("no database connections")
	}
	return nil
}

func (hc *HealthChecker) checkKafka(ctx context.Context) error {
	conn, err := kafka.DialContext(ctx, "tcp", hc.brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Brokers()
	return err
}

func (hc *HealthChecker) respondWithError(w http.ResponseWriter, code int, response ReadinessResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

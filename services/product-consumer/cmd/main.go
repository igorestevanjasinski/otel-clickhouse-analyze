package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/igor-jasinski/product-consumer/internal/config"
	"github.com/igor-jasinski/product-consumer/internal/consumer"
	"github.com/igor-jasinski/product-consumer/internal/health"
	"github.com/igor-jasinski/product-consumer/internal/repository"
	"github.com/igor-jasinski/product-consumer/pkg/metrics"
	"github.com/igor-jasinski/product-consumer/pkg/telemetry"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	setupLogger(cfg.App.LogLevel)

	log.WithFields(logrus.Fields{
		"kafka_brokers":  cfg.Kafka.Brokers,
		"kafka_topic":    cfg.Kafka.Topic,
		"kafka_group_id": cfg.Kafka.GroupID,
		"postgres_host":  cfg.Postgres.Host,
		"postgres_db":    cfg.Postgres.Database,
	}).Info("Application starting")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigChan
		log.WithField("signal", sig.String()).Info("Received shutdown signal")
		cancel()
	}()

	if err := run(ctx, cfg); err != nil {
		log.WithError(err).Error("Application failed")
		os.Exit(1)
	}

	log.Info("Application stopped gracefully")
}

func run(ctx context.Context, cfg *config.Config) error {
	metrics.Init()
	metrics.StartServer("8081")
	log.Info("Prometheus metrics server started on :8081")

	shutdownTracing, err := telemetry.SetupTracing(telemetry.TracingConfig{
		ServiceName:    "product-consumer",
		ServiceVersion: "1.0.0",
		Environment:    cfg.App.Environment,
		OTLPEndpoint:   cfg.OpenTelemetry.Endpoint,
	})
	if err != nil {
		log.WithError(err).Error("Failed to setup tracing")
	} else {
		defer func() {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := shutdownTracing(shutdownCtx); err != nil {
				log.WithError(err).Error("Failed to shutdown tracing")
			}
		}()
	}

	shutdownMetrics, err := telemetry.SetupMetrics(telemetry.MetricsConfig{
		ServiceName:    "product-consumer",
		ServiceVersion: "1.0.0",
		Environment:    cfg.App.Environment,
		OTLPEndpoint:   cfg.OpenTelemetry.Endpoint,
	})
	if err != nil {
		log.WithError(err).Error("Failed to setup metrics")
	} else {
		defer func() {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := shutdownMetrics(shutdownCtx); err != nil {
				log.WithError(err).Error("Failed to shutdown metrics")
			}
		}()
	}

	repo, err := repository.NewProductRepository(cfg, log)
	if err != nil {
		return err
	}
	defer func() {
		if err := repo.Close(); err != nil {
			log.WithError(err).Error("Failed to close repository")
		}
	}()

	dlqTopic := cfg.Kafka.Topic + "-dlq"

	chaosConfig := &consumer.ChaosConfig{
		Enabled:      cfg.Chaos.Enabled,
		ErrorRate:    cfg.Chaos.ErrorRate,
		LatencyMinMs: cfg.Chaos.LatencyMinMs,
		LatencyMaxMs: cfg.Chaos.LatencyMaxMs,
		Logger:       log,
	}

	if chaosConfig.Enabled {
		log.WithFields(logrus.Fields{
			"error_rate":     chaosConfig.ErrorRate,
			"latency_min_ms": chaosConfig.LatencyMinMs,
			"latency_max_ms": chaosConfig.LatencyMaxMs,
		}).Warn("Chaos Engineering enabled")
	}

	handler := consumer.NewMessageHandler(repo, log, cfg.Kafka.Brokers, dlqTopic, chaosConfig)
	defer func() {
		if err := handler.Close(); err != nil {
			log.WithError(err).Error("Failed to close message handler")
		}
	}()

	healthServer := health.NewHealthServer(8081, repo, cfg.Kafka.Brokers, log)
	go func() {
		if err := healthServer.Start(); err != nil {
			log.WithError(err).Error("Health server failed")
		}
	}()
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := healthServer.Shutdown(shutdownCtx); err != nil {
			log.WithError(err).Error("Failed to shutdown health server")
		}
	}()

	kafkaConsumer := consumer.NewKafkaConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.Topic,
		cfg.Kafka.GroupID,
		handler,
		log,
	)
	defer func() {
		if err := kafkaConsumer.Close(); err != nil {
			log.WithError(err).Error("Failed to close Kafka consumer")
		}
	}()

	consumerErr := make(chan error, 1)
	go func() {
		consumerErr <- kafkaConsumer.Start(ctx)
	}()

	log.Info("Consumer is ready to process messages")

	select {
	case <-ctx.Done():
		log.Info("Shutdown initiated")
	case err := <-consumerErr:
		if err != nil && err != context.Canceled {
			log.WithError(err).Error("Consumer failed")
			return err
		}
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		cfg.App.ShutdownTimeout,
	)
	defer shutdownCancel()

	shutdownComplete := make(chan struct{})
	go func() {
		log.Info("Performing graceful shutdown")
		close(shutdownComplete)
	}()

	select {
	case <-shutdownComplete:
		log.Info("Graceful shutdown completed")
	case <-shutdownCtx.Done():
		log.Warn("Shutdown timeout exceeded, forcing exit")
	}

	return nil
}

func setupLogger(level string) {
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	log.SetOutput(os.Stdout)

	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		log.Warnf("Invalid log level '%s', using INFO", level)
		logLevel = logrus.InfoLevel
	}
	log.SetLevel(logLevel)
}

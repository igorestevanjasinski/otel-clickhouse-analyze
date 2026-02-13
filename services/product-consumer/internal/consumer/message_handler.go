package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/igor-jasinski/product-consumer/internal/models"
	"github.com/igor-jasinski/product-consumer/internal/repository"
	"github.com/igor-jasinski/product-consumer/pkg/telemetry"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	maxRetries     = 3
	baseRetryDelay = 100 * time.Millisecond
)

type MessageHandler struct {
	repo        repository.ProductRepository
	logger      *logrus.Logger
	dlqWriter   *kafka.Writer
	chaosConfig *ChaosConfig
}

func NewMessageHandler(repo repository.ProductRepository, logger *logrus.Logger, brokers []string, dlqTopic string, chaosConfig *ChaosConfig) *MessageHandler {
	dlqWriter := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        dlqTopic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		Async:        false,
	}

	return &MessageHandler{
		repo:        repo,
		logger:      logger,
		dlqWriter:   dlqWriter,
		chaosConfig: chaosConfig,
	}
}

func (h *MessageHandler) ProcessMessage(ctx context.Context, msg kafka.Message) error {
	carrier := make(propagation.HeaderCarrier)
	for _, header := range msg.Headers {
		carrier.Set(header.Key, string(header.Value))
	}

	propagator := propagation.TraceContext{}
	ctx = propagator.Extract(ctx, carrier)

	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "kafka.consume",
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.destination", msg.Topic),
			attribute.Int("messaging.partition", msg.Partition),
			attribute.Int64("messaging.offset", msg.Offset),
		),
	)
	defer span.End()

	correlationID := h.extractCorrelationID(msg)
	traceID := span.SpanContext().TraceID().String()

	log := h.logger.WithFields(logrus.Fields{
		"correlation_id": correlationID,
		"trace_id":       traceID,
	})

	if telemetry.AppMetrics != nil {
		telemetry.AppMetrics.MessagesInFlight.Add(ctx, 1)
		defer telemetry.AppMetrics.MessagesInFlight.Add(ctx, -1)

		start := time.Now()
		defer func() {
			duration := time.Since(start).Seconds()
			telemetry.AppMetrics.ProcessingDuration.Record(ctx, duration)
			telemetry.AppMetrics.MessagesConsumed.Add(ctx, 1)
		}()
	}

	InjectLatency(h.chaosConfig)

	if err := InjectRandomError(h.chaosConfig); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Chaos error injected")
		log.WithError(err).Error("Chaos error injected, skipping message processing")
		if telemetry.AppMetrics != nil {
			telemetry.AppMetrics.ProcessingErrors.Add(ctx, 1)
		}
		return err
	}

	var product models.Product
	if err := json.Unmarshal(msg.Value, &product); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Deserialization failed")
		log.WithError(err).Error("Failed to deserialize message")
		h.handleDeserializationError(ctx, msg, correlationID)
		if telemetry.AppMetrics != nil {
			telemetry.AppMetrics.ProcessingErrors.Add(ctx, 1)
		}
		return nil
	}

	if err := product.Validate(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Validation failed")
		log.WithError(err).Error("Invalid product data")
		h.handleDeserializationError(ctx, msg, correlationID)
		if telemetry.AppMetrics != nil {
			telemetry.AppMetrics.ProcessingErrors.Add(ctx, 1)
		}
		return nil
	}

	span.SetAttributes(
		attribute.String("product.id", product.ID),
		attribute.String("product.name", product.Name),
		attribute.Float64("product.price", product.Price),
	)

	log.WithFields(logrus.Fields{
		"product_id":   product.ID,
		"product_name": product.Name,
		"price":        product.Price,
	}).Info("Processing product message")

	if err := h.createProductWithRetry(ctx, &product, log); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to persist product")
		if telemetry.AppMetrics != nil {
			telemetry.AppMetrics.ProcessingErrors.Add(ctx, 1)
		}
		return fmt.Errorf("failed to persist product after retries: %w", err)
	}

	span.SetStatus(codes.Ok, "Product persisted successfully")
	log.WithField("product_id", product.ID).Info("Product persisted successfully")
	return nil
}

func (h *MessageHandler) createProductWithRetry(ctx context.Context, product *models.Product, log *logrus.Entry) error {
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		start := time.Now()
		err := h.repo.CreateProduct(ctx, product)
		duration := time.Since(start).Milliseconds()

		if err == nil {
			log.WithFields(logrus.Fields{
				"attempt":  attempt,
				"duration": duration,
			}).Info("Product created successfully")
			return nil
		}

		lastErr = err

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			log.WithError(err).Warn("Context cancelled, stopping retries")
			return err
		}

		log.WithFields(logrus.Fields{
			"attempt":  attempt,
			"duration": duration,
			"error":    err.Error(),
		}).Warn("Failed to create product, will retry")

		if attempt < maxRetries {
			delay := baseRetryDelay * time.Duration(1<<uint(attempt-1))
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return fmt.Errorf("exhausted retries (%d attempts): %w", maxRetries, lastErr)
}

func (h *MessageHandler) handleDeserializationError(ctx context.Context, msg kafka.Message, correlationID string) {
	dlqMessage := kafka.Message{
		Key:   msg.Key,
		Value: msg.Value,
		Headers: []kafka.Header{
			{Key: "original-topic", Value: []byte(msg.Topic)},
			{Key: "original-partition", Value: []byte(fmt.Sprintf("%d", msg.Partition))},
			{Key: "original-offset", Value: []byte(fmt.Sprintf("%d", msg.Offset))},
			{Key: "correlation-id", Value: []byte(correlationID)},
			{Key: "error-timestamp", Value: []byte(time.Now().Format(time.RFC3339))},
		},
	}

	if err := h.dlqWriter.WriteMessages(ctx, dlqMessage); err != nil {
		h.logger.WithFields(logrus.Fields{
			"correlation_id": correlationID,
			"error":          err.Error(),
		}).Error("Failed to send message to DLQ")
		return
	}

	h.logger.WithField("correlation_id", correlationID).Warn("Invalid message sent to DLQ")
}

func (h *MessageHandler) extractCorrelationID(msg kafka.Message) string {
	for _, header := range msg.Headers {
		if header.Key == "correlation_id" {
			return string(header.Value)
		}
	}
	return "unknown"
}

func (h *MessageHandler) Close() error {
	if h.dlqWriter != nil {
		return h.dlqWriter.Close()
	}
	return nil
}

package consumer

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	reader  *kafka.Reader
	handler *MessageHandler
	logger  *logrus.Logger
}

func NewKafkaConsumer(brokers []string, topic string, groupID string, handler *MessageHandler, logger *logrus.Logger) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: 0,
		StartOffset:    kafka.LastOffset,
		MaxWait:        500 * time.Millisecond,
		Logger:         kafka.LoggerFunc(func(msg string, args ...interface{}) {}),
		ErrorLogger:    kafka.LoggerFunc(func(msg string, args ...interface{}) {}),
	})

	return &KafkaConsumer{
		reader:  reader,
		handler: handler,
		logger:  logger,
	}
}

func (c *KafkaConsumer) Start(ctx context.Context) error {
	c.logger.WithFields(logrus.Fields{
		"topic":    c.reader.Config().Topic,
		"group_id": c.reader.Config().GroupID,
	}).Info("Starting Kafka consumer")

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Context cancelled, stopping consumer")
			return ctx.Err()
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if err == context.Canceled {
					return nil
				}
				c.logger.WithError(err).Error("Failed to fetch message")
				continue
			}

			c.processMessage(ctx, msg)
		}
	}
}

func (c *KafkaConsumer) processMessage(ctx context.Context, msg kafka.Message) {
	start := time.Now()

	log := c.logger.WithFields(logrus.Fields{
		"topic":     msg.Topic,
		"partition": msg.Partition,
		"offset":    msg.Offset,
	})

	log.Info("Received message")

	if err := c.handler.ProcessMessage(ctx, msg); err != nil {
		log.WithFields(logrus.Fields{
			"error":    err.Error(),
			"duration": time.Since(start).Milliseconds(),
		}).Error("Failed to process message")
		return
	}

	if err := c.reader.CommitMessages(ctx, msg); err != nil {
		log.WithError(err).Error("Failed to commit message")
		return
	}

	log.WithField("duration", time.Since(start).Milliseconds()).Info("Message committed successfully")
}

func (c *KafkaConsumer) Close() error {
	c.logger.Info("Closing Kafka consumer")
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("failed to close reader: %w", err)
	}
	return nil
}

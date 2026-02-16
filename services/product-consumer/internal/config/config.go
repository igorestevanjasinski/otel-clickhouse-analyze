package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

type AppConfig struct {
	LogLevel        string
	Environment     string
	ShutdownTimeout time.Duration
}

type ChaosConfig struct {
	Enabled      bool
	ErrorRate    float64
	LatencyMinMs int
	LatencyMaxMs int
}

type OpenTelemetryConfig struct {
	Endpoint string
}

type Config struct {
	Kafka         KafkaConfig
	Postgres      PostgresConfig
	App           AppConfig
	Chaos         ChaosConfig
	OpenTelemetry OpenTelemetryConfig
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{
		Kafka: KafkaConfig{
			Brokers: []string{getEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")},
			Topic:   getEnv("KAFKA_TOPIC_PRODUCTS", "products.events"),
			GroupID: getEnv("KAFKA_GROUP_ID", "product-consumer-group"),
		},
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnvInt("POSTGRES_PORT", 5432),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			Database: getEnv("POSTGRES_DB", "products_db"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},
		App: AppConfig{
			LogLevel:        getEnv("LOG_LEVEL", "INFO"),
			Environment:     getEnv("ENVIRONMENT", "development"),
			ShutdownTimeout: getEnvDuration("SHUTDOWN_TIMEOUT", 30*time.Second),
		},
		Chaos: ChaosConfig{
			Enabled:      getEnvBool("ENABLE_CHAOS", false),
			ErrorRate:    getEnvFloat("CHAOS_ERROR_RATE", 0.1),
			LatencyMinMs: getEnvInt("CHAOS_LATENCY_MIN_MS", 100),
			LatencyMaxMs: getEnvInt("CHAOS_LATENCY_MAX_MS", 2000),
		},
		OpenTelemetry: OpenTelemetryConfig{
			Endpoint: getEnv("OTLP_ENDPOINT", "localhost:4317"),
		},
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

func (c *Config) Validate() error {
	if len(c.Kafka.Brokers) == 0 {
		return errors.New("kafka brokers are required")
	}

	if c.Kafka.Topic == "" {
		return errors.New("kafka topic is required")
	}

	if c.Kafka.GroupID == "" {
		return errors.New("kafka group ID is required")
	}

	if c.Postgres.Host == "" {
		return errors.New("postgres host is required")
	}

	if c.Postgres.User == "" {
		return errors.New("postgres user is required")
	}

	if c.Postgres.Database == "" {
		return errors.New("postgres database is required")
	}

	return nil
}

func (c *PostgresConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

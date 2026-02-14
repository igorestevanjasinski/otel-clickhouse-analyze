package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/igor-jasinski/product-consumer/internal/config"
	"github.com/igor-jasinski/product-consumer/internal/models"
	"github.com/igor-jasinski/product-consumer/pkg/metrics"
	"github.com/igor-jasinski/product-consumer/pkg/telemetry"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrDuplicateKey    = errors.New("product already exists")
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *models.Product) error
	GetProductByID(ctx context.Context, id string) (*models.Product, error)
	Stats() *pgxpool.Stat
	Close() error
}

type productRepository struct {
	pool *pgxpool.Pool
	log  *logrus.Logger
}

func NewProductRepository(cfg *config.Config, logger *logrus.Logger) (ProductRepository, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.Postgres.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = time.Minute
	poolConfig.ConnConfig.ConnectTimeout = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"host":      cfg.Postgres.Host,
		"database":  cfg.Postgres.Database,
		"max_conns": poolConfig.MaxConns,
		"min_conns": poolConfig.MinConns,
	}).Info("Database connection pool initialized")

	return &productRepository{
		pool: pool,
		log:  logger,
	}, nil
}

func (r *productRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "database.insert",
		trace.WithAttributes(
			attribute.String("db.system", "postgresql"),
			attribute.String("db.operation", "INSERT"),
			attribute.String("db.table", "products"),
			attribute.String("product.id", product.ID),
		),
	)
	defer span.End()

	query := `
		INSERT INTO products (id, name, price, description, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO NOTHING
	`

	start := time.Now()
	result, err := r.pool.Exec(
		ctx,
		query,
		product.ID,
		product.Name,
		product.Price,
		product.Description,
		product.CreatedAt,
	)

	duration := time.Since(start)

	metrics.DatabaseOperationDuration.WithLabelValues("insert").Observe(duration.Seconds())

	if telemetry.AppMetrics != nil {
		telemetry.AppMetrics.DatabaseOperations.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("operation", "insert"),
				attribute.String("table", "products"),
			),
		)
	}

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to insert product")
		metrics.DatabaseOperations.WithLabelValues("insert", "error").Inc()
		r.log.WithFields(logrus.Fields{
			"product_id": product.ID,
			"error":      err.Error(),
			"duration":   duration.Milliseconds(),
		}).Error("Failed to insert product")
		return fmt.Errorf("failed to insert product: %w", err)
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		span.AddEvent("Product already exists")
		metrics.DatabaseOperations.WithLabelValues("insert", "duplicate").Inc()
		r.log.WithFields(logrus.Fields{
			"product_id": product.ID,
			"duration":   duration.Milliseconds(),
		}).Warn("Product already exists, skipping insert")
		return nil
	}

	span.SetStatus(codes.Ok, "Product created successfully")
	span.SetAttributes(attribute.Int64("db.rows_affected", rowsAffected))
	metrics.DatabaseOperations.WithLabelValues("insert", "success").Inc()

	r.log.WithFields(logrus.Fields{
		"product_id":   product.ID,
		"product_name": product.Name,
		"price":        product.Price,
		"duration":     duration.Milliseconds(),
	}).Info("Product created successfully")

	return nil
}

func (r *productRepository) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	query := `
		SELECT id, name, price, description, created_at
		FROM products
		WHERE id = $1
	`

	var product models.Product
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Description,
		&product.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &product, nil
}

func (r *productRepository) Close() error {
	r.log.Info("Closing database connection pool")
	r.pool.Close()
	return nil
}

func (r *productRepository) Stats() *pgxpool.Stat {
	return r.pool.Stat()
}

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/igor-jasinski/product-consumer/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupPostgresContainer(t *testing.T) (testcontainers.Container, string) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	connString := "postgres://testuser:testpass@" + host + ":" + port.Port() + "/testdb?sslmode=disable"

	return container, connString
}

func TestProductRepository_CreateProduct_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	container, connString := setupPostgresContainer(t)
	defer container.Terminate(context.Background())

	time.Sleep(2 * time.Second)

	repo, err := NewProductRepository(connString)
	require.NoError(t, err)
	defer repo.Close()

	ctx := context.Background()

	_, err = repo.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS products (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			price DECIMAL(10,2) NOT NULL,
			description TEXT,
			created_at TIMESTAMP NOT NULL
		)
	`)
	require.NoError(t, err)

	product := &models.Product{
		ID:          "test-123",
		Name:        "Test Product",
		Price:       99.99,
		Description: "Test description",
		CreatedAt:   time.Now(),
	}

	err = repo.CreateProduct(ctx, product)
	assert.NoError(t, err)
}

func TestProductRepository_CreateProduct_Duplicate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	container, connString := setupPostgresContainer(t)
	defer container.Terminate(context.Background())

	time.Sleep(2 * time.Second)

	repo, err := NewProductRepository(connString)
	require.NoError(t, err)
	defer repo.Close()

	ctx := context.Background()

	_, err = repo.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS products (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			price DECIMAL(10,2) NOT NULL,
			description TEXT,
			created_at TIMESTAMP NOT NULL
		)
	`)
	require.NoError(t, err)

	product := &models.Product{
		ID:          "duplicate-123",
		Name:        "Test Product",
		Price:       99.99,
		Description: "Test",
		CreatedAt:   time.Now(),
	}

	err = repo.CreateProduct(ctx, product)
	require.NoError(t, err)

	err = repo.CreateProduct(ctx, product)
	assert.NoError(t, err)
}

func TestProductRepository_Stats(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	container, connString := setupPostgresContainer(t)
	defer container.Terminate(context.Background())

	time.Sleep(2 * time.Second)

	repo, err := NewProductRepository(connString)
	require.NoError(t, err)
	defer repo.Close()

	stats := repo.Stats()
	assert.NotNil(t, stats)
}

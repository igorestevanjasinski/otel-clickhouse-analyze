package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/igor-jasinski/product-consumer/internal/config"
	"github.com/igor-jasinski/product-consumer/internal/models"
	"github.com/igor-jasinski/product-consumer/internal/repository"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	repo, err := repository.NewProductRepository(cfg, log)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	ctx := context.Background()

	productID1 := uuid.New().String()
	product1 := &models.Product{
		ID:          productID1,
		Name:        "Test Product 1",
		Price:       99.99,
		Description: "First test product",
		CreatedAt:   time.Now(),
	}

	log.Info("Creating first product")
	if err := repo.CreateProduct(ctx, product1); err != nil {
		log.Fatalf("Failed to create product: %v", err)
	}

	log.Info("Creating duplicate product (should be idempotent)")
	if err := repo.CreateProduct(ctx, product1); err != nil {
		log.Fatalf("Failed on duplicate: %v", err)
	}

	log.Info("Retrieving product")
	retrieved, err := repo.GetProductByID(ctx, product1.ID)
	if err != nil {
		log.Fatalf("Failed to retrieve product: %v", err)
	}

	log.WithFields(logrus.Fields{
		"id":    retrieved.ID,
		"name":  retrieved.Name,
		"price": retrieved.Price,
	}).Info("Product retrieved successfully")

	log.Info("Testing non-existent product")
	nonExistentID := uuid.New().String()
	_, err = repo.GetProductByID(ctx, nonExistentID)
	if err == repository.ErrProductNotFound {
		log.WithField("product_id", nonExistentID).Info("Correctly returned ErrProductNotFound")
	} else {
		log.Fatalf("Expected ErrProductNotFound, got: %v", err)
	}

	fmt.Println("\nAll tests passed!")
}

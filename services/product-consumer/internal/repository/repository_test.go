package repository

import (
	"context"
	"testing"
	"time"

	"github.com/igor-jasinski/product-consumer/internal/models"
)

func TestProductValidation(t *testing.T) {
	tests := []struct {
		name    string
		product *models.Product
		wantErr bool
	}{
		{
			name: "valid product",
			product: &models.Product{
				ID:          "123e4567-e89b-12d3-a456-426614174000",
				Name:        "Test Product",
				Price:       99.99,
				Description: "Test description",
				CreatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "invalid product - empty ID",
			product: &models.Product{
				Name:      "Test Product",
				Price:     99.99,
				CreatedAt: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	select {
	case <-ctx.Done():
		if ctx.Err() != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", ctx.Err())
		}
	default:
		t.Error("context should be cancelled")
	}
}

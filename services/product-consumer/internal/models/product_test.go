package models

import (
	"testing"
	"time"
)

func TestProduct_Validate(t *testing.T) {
	tests := []struct {
		name    string
		product Product
		wantErr bool
	}{
		{
			name: "valid product",
			product: Product{
				ID:          "123e4567-e89b-12d3-a456-426614174000",
				Name:        "Test Product",
				Price:       99.99,
				Description: "Test description",
				CreatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			product: Product{
				Name:      "Test Product",
				Price:     99.99,
				CreatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "missing name",
			product: Product{
				ID:        "123e4567-e89b-12d3-a456-426614174000",
				Price:     99.99,
				CreatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "name too short",
			product: Product{
				ID:        "123e4567-e89b-12d3-a456-426614174000",
				Name:      "AB",
				Price:     99.99,
				CreatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "negative price",
			product: Product{
				ID:        "123e4567-e89b-12d3-a456-426614174000",
				Name:      "Test Product",
				Price:     -10.00,
				CreatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "zero price",
			product: Product{
				ID:        "123e4567-e89b-12d3-a456-426614174000",
				Name:      "Test Product",
				Price:     0,
				CreatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "missing created_at",
			product: Product{
				ID:    "123e4567-e89b-12d3-a456-426614174000",
				Name:  "Test Product",
				Price: 99.99,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Product.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

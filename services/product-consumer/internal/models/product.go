package models

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type Product struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Price       float64   `json:"price" db:"price"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

func (p *Product) UnmarshalJSON(data []byte) error {
	type Alias Product
	aux := &struct {
		CreatedAt string `json:"created_at"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05.999999",
	}

	var parseErr error
	for _, format := range formats {
		createdAt := strings.TrimSpace(aux.CreatedAt)
		if t, err := time.Parse(format, createdAt); err == nil {
			p.CreatedAt = t
			return nil
		} else {
			parseErr = err
		}
	}

	return parseErr
}

func (p *Product) Validate() error {
	if p.ID == "" {
		return errors.New("product ID is required")
	}

	if p.Name == "" {
		return errors.New("product name is required")
	}

	if len(p.Name) < 3 || len(p.Name) > 255 {
		return errors.New("product name must be between 3 and 255 characters")
	}

	if p.Price <= 0 {
		return errors.New("product price must be greater than zero")
	}

	if p.CreatedAt.IsZero() {
		return errors.New("product created_at is required")
	}

	return nil
}

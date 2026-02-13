package consumer

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/igor-jasinski/product-consumer/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockRepository) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockRepository) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRepository) Stats() *pgxpool.Stat {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*pgxpool.Stat)
}

func TestProcessMessage_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	logger := logrus.New()
	logger.SetOutput(&testWriter{})

	product := &models.Product{
		ID:          "test-123",
		Name:        "Test Product",
		Price:       99.99,
		Description: "Test description",
		CreatedAt:   time.Now(),
	}

	mockRepo.On("CreateProduct", mock.Anything, mock.AnythingOfType("*models.Product")).Return(nil)

	handler := NewMessageHandler(mockRepo, logger, []string{"localhost:9092"}, "dlq-topic", &ChaosConfig{Enabled: false})

	msgBytes, _ := json.Marshal(product)
	msg := kafka.Message{
		Value: msgBytes,
		Headers: []kafka.Header{
			{Key: "correlation-id", Value: []byte("test-correlation")},
		},
	}

	err := handler.ProcessMessage(context.Background(), msg)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "CreateProduct", mock.Anything, mock.AnythingOfType("*models.Product"))
}

func TestProcessMessage_InvalidJSON(t *testing.T) {
	mockRepo := new(MockRepository)
	logger := logrus.New()
	logger.SetOutput(&testWriter{})

	handler := NewMessageHandler(mockRepo, logger, []string{"localhost:9092"}, "dlq-topic", &ChaosConfig{Enabled: false})

	msg := kafka.Message{
		Value: []byte("invalid json"),
	}

	err := handler.ProcessMessage(context.Background(), msg)

	assert.NoError(t, err)
}

func TestShouldFail(t *testing.T) {
	result1 := ShouldFail(0.0)
	assert.False(t, result1)

	result2 := ShouldFail(1.0)
	assert.True(t, result2)
}

type testWriter struct{}

func (tw *testWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

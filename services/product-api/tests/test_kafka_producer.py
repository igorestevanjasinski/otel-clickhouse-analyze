import pytest
from unittest.mock import Mock, patch, MagicMock
from confluent_kafka import KafkaException
from app.services.kafka_producer import KafkaProducerService


class TestKafkaProducerService:
    
    @pytest.fixture
    def mock_producer(self):
        with patch('app.services.kafka_producer.Producer') as mock:
            producer_instance = MagicMock()
            mock.return_value = producer_instance
            yield producer_instance
    
    @pytest.fixture
    def kafka_service(self, mock_producer):
        service = KafkaProducerService(
            bootstrap_servers="localhost:9092",
            topic="test-topic",
            retries=3,
            retry_backoff_ms=100,
            acks="all"
        )
        return service
    
    def test_send_product_event_success(self, kafka_service, mock_producer):
        product_data = {
            "id": "test-id",
            "name": "Test Product",
            "price": 99.99,
            "description": "Test",
            "created_at": "2024-01-01T00:00:00"
        }
        
        mock_producer.flush.return_value = 0
        
        result = kafka_service.send_product_event(product_data, "correlation-123")
        
        assert result is True
        mock_producer.produce.assert_called_once()
        mock_producer.flush.assert_called_once()
    
    def test_send_product_event_kafka_exception(self, kafka_service, mock_producer):
        product_data = {
            "id": "test-id",
            "name": "Test Product",
            "price": 99.99
        }
        
        mock_producer.produce.side_effect = KafkaException("Connection error")
        
        result = kafka_service.send_product_event(product_data, "correlation-123")
        
        assert result is False
    
    def test_send_product_event_with_retry(self, kafka_service, mock_producer):
        product_data = {
            "id": "test-id",
            "name": "Test Product",
            "price": 99.99
        }
        
        mock_producer.produce.side_effect = [
            KafkaException("Temp error"),
            KafkaException("Temp error"),
            None
        ]
        mock_producer.flush.return_value = 0
        
        result = kafka_service.send_product_event(product_data, "correlation-123")
        
        assert result is True
        assert mock_producer.produce.call_count == 3
    
    def test_close(self, kafka_service, mock_producer):
        kafka_service.close()
        
        mock_producer.flush.assert_called()

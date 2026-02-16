import pytest
from fastapi.testclient import TestClient
from unittest.mock import Mock
from app.main import create_application
from app.services import KafkaProducerService
from app.api.routes import set_kafka_producer


@pytest.fixture
def mock_kafka_producer():
    mock_producer = Mock(spec=KafkaProducerService)
    mock_producer.send_product_event.return_value = True
    return mock_producer


@pytest.fixture
def client(mock_kafka_producer):
    app = create_application()
    set_kafka_producer(mock_kafka_producer)
    with TestClient(app) as test_client:
        yield test_client

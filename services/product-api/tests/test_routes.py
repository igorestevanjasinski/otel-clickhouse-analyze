import pytest
from fastapi import status


class TestProductRoutes:
    
    def test_create_product_missing_name(self, client):
        response = client.post(
            "/products",
            json={
                "price": 99.99,
                "description": "Test"
            }
        )
        
        assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY
    
    def test_create_product_missing_price(self, client):
        response = client.post(
            "/products",
            json={
                "name": "Test Product",
                "description": "Test"
            }
        )
        
        assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY
    
    def test_create_product_invalid_price_type(self, client):
        response = client.post(
            "/products",
            json={
                "name": "Test Product",
                "price": "invalid",
                "description": "Test"
            }
        )
        
        assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY
    
    def test_create_product_negative_price(self, client):
        response = client.post(
            "/products",
            json={
                "name": "Test Product",
                "price": -10.00,
                "description": "Test"
            }
        )
        
        assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY
    
    def test_create_product_zero_price(self, client):
        response = client.post(
            "/products",
            json={
                "name": "Test Product",
                "price": 0,
                "description": "Test"
            }
        )
        
        assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY
    
    def test_create_product_name_too_short(self, client):
        response = client.post(
            "/products",
            json={
                "name": "AB",
                "price": 99.99,
                "description": "Test"
            }
        )
        
        assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY


class TestHealthRoutes:
    
    def test_health_check(self, client):
        response = client.get("/health")
        
        assert response.status_code == status.HTTP_200_OK
        data = response.json()
        assert data["status"] == "healthy"
    
    def test_ready_check(self, client, mock_kafka_producer):
        response = client.get("/ready")
        
        assert response.status_code == status.HTTP_200_OK

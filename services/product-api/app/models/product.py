from datetime import datetime
from decimal import Decimal
from uuid import UUID
from pydantic import BaseModel, Field, field_validator


class ProductCreate(BaseModel):
    name: str = Field(..., min_length=3, max_length=255, description="Product name")
    price: Decimal = Field(..., gt=0, decimal_places=2, description="Product price")
    description: str | None = Field(None, max_length=1000, description="Product description")

    @field_validator("name")
    @classmethod
    def validate_name(cls, v: str) -> str:
        if not v.strip():
            raise ValueError("Name cannot be empty or whitespace")
        return v.strip()

    @field_validator("price")
    @classmethod
    def validate_price(cls, v: Decimal) -> Decimal:
        if v <= 0:
            raise ValueError("Price must be greater than 0")
        return round(v, 2)

    model_config = {
        "json_schema_extra": {
            "examples": [
                {
                    "name": "Laptop Dell XPS 13",
                    "price": 1299.99,
                    "description": "High-performance ultrabook"
                }
            ]
        }
    }


class ProductResponse(BaseModel):
    id: UUID
    name: str
    price: Decimal
    description: str | None
    created_at: datetime

    model_config = {
        "json_schema_extra": {
            "examples": [
                {
                    "id": "123e4567-e89b-12d3-a456-426614174000",
                    "name": "Laptop Dell XPS 13",
                    "price": 1299.99,
                    "description": "High-performance ultrabook",
                    "created_at": "2024-01-01T12:00:00Z"
                }
            ]
        }
    }

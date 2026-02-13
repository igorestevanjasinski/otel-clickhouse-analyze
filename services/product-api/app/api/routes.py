from datetime import datetime
from uuid import uuid4
from fastapi import APIRouter, Depends, HTTPException, status
from app.models import ProductCreate, ProductResponse
from app.services import KafkaProducerService
import logging

router = APIRouter()
logger = logging.getLogger(__name__)

_kafka_producer: KafkaProducerService | None = None


def set_kafka_producer(producer: KafkaProducerService):
    global _kafka_producer
    _kafka_producer = producer


def get_kafka_producer() -> KafkaProducerService:
    if _kafka_producer is None:
        raise RuntimeError("Kafka producer not initialized")
    return _kafka_producer


@router.post(
    "/products",
    response_model=ProductResponse,
    status_code=status.HTTP_202_ACCEPTED,
    summary="Create a new product",
    description="Create a new product and publish event to Kafka"
)
async def create_product(
    product: ProductCreate,
    kafka_producer: KafkaProducerService = Depends(get_kafka_producer)
) -> ProductResponse:
    correlation_id = str(uuid4())
    product_id = uuid4()
    
    try:
        logger.info(
            "Product creation initiated",
            extra={
                "correlation_id": correlation_id,
                "product_id": str(product_id),
                "product_name": product.name,
                "product_price": float(product.price)
            }
        )
        
        response = ProductResponse(
            id=product_id,
            name=product.name,
            price=product.price,
            description=product.description,
            created_at=datetime.utcnow()
        )
        
        product_event = {
            "id": str(response.id),
            "name": response.name,
            "price": float(response.price),
            "description": response.description,
            "created_at": response.created_at.isoformat()
        }
        
        success = kafka_producer.send_product_event(product_event, correlation_id)
        
        if not success:
            logger.error(
                "Failed to send product event to Kafka after retries",
                extra={
                    "correlation_id": correlation_id,
                    "product_id": str(product_id)
                }
            )
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail="Failed to publish product event"
            )
        
        logger.info(
            "Product created successfully",
            extra={
                "correlation_id": correlation_id,
                "product_id": str(product_id)
            }
        )
        
        return response
        
    except HTTPException:
        raise
    except Exception as e:
        logger.error(
            "Failed to create product",
            extra={
                "correlation_id": correlation_id,
                "error": str(e),
                "product_name": product.name
            },
            exc_info=True
        )
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to create product"
        )


@router.get(
    "/health",
    status_code=status.HTTP_200_OK,
    summary="Health check",
    description="Check if the service is running"
)
async def health_check() -> dict[str, str]:
    return {"status": "healthy"}


@router.get(
    "/ready",
    status_code=status.HTTP_200_OK,
    summary="Readiness check",
    description="Check if the service and its dependencies are ready"
)
async def readiness_check(
    kafka_producer: KafkaProducerService = Depends(get_kafka_producer)
) -> dict[str, str]:
    try:
        if kafka_producer.producer is None:
            raise HTTPException(
                status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
                detail="Kafka producer not initialized"
            )
        return {"status": "ready", "kafka": "connected"}
    except Exception as e:
        logger.error(f"Readiness check failed: {str(e)}")
        raise HTTPException(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
            detail="Service dependencies not ready"
        )

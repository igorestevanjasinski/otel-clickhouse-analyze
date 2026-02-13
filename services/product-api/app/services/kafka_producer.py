import json
import logging
import time
from typing import Any, Dict, Optional
from confluent_kafka import Producer
from confluent_kafka import KafkaException, KafkaError as ConfluentKafkaError
from opentelemetry import trace
from opentelemetry.trace.propagation.tracecontext import TraceContextTextMapPropagator


logger = logging.getLogger(__name__)


class KafkaProducerService:
    def __init__(
        self,
        bootstrap_servers: str,
        topic: str,
        retries: int = 3,
        retry_backoff_ms: int = 100,
        acks: str = "all"
    ):
        self.topic = topic
        self.max_retries = retries
        self.retry_backoff_ms = retry_backoff_ms
        
        conf = {
            'bootstrap.servers': bootstrap_servers,
            'acks': acks,
            'retries': retries,
            'compression.type': 'gzip',
            'linger.ms': 10,
            'batch.size': 16384,
            'max.in.flight.requests.per.connection': 5,
            'client.id': 'product-api-producer'
        }
        
        self.producer = Producer(conf)
        
        logger.info(
            "Kafka producer initialized",
            extra={
                "bootstrap_servers": bootstrap_servers,
                "topic": topic,
                "acks": acks
            }
        )

    def _delivery_callback(self, err, msg, correlation_id: str):
        if err:
            logger.error(
                "Message delivery failed",
                extra={
                    "correlation_id": correlation_id,
                    "error": str(err),
                    "topic": msg.topic() if msg else None
                }
            )
        else:
            logger.debug(
                "Message delivered",
                extra={
                    "correlation_id": correlation_id,
                    "topic": msg.topic(),
                    "partition": msg.partition(),
                    "offset": msg.offset()
                }
            )

    def send_product_event(
        self,
        product_data: Dict[str, Any],
        correlation_id: Optional[str] = None
    ) -> bool:
        attempt = 0
        last_error = None
        
        tracer = trace.get_tracer(__name__)
        
        with tracer.start_as_current_span("kafka.publish") as span:
            span.set_attribute("messaging.system", "kafka")
            span.set_attribute("messaging.destination", self.topic)
            span.set_attribute("messaging.operation", "publish")
            span.set_attribute("product.id", product_data.get("id"))
            
            headers = {}
            TraceContextTextMapPropagator().inject(headers)
            
            kafka_headers = [(k, v.encode('utf-8')) for k, v in headers.items()]
            
            while attempt < self.max_retries:
                try:
                    value = json.dumps(product_data).encode('utf-8')
                    
                    self.producer.produce(
                        self.topic,
                        value=value,
                        headers=kafka_headers,
                        callback=lambda err, msg: self._delivery_callback(err, msg, correlation_id)
                    )
                    
                    self.producer.poll(0)
                    self.producer.flush(timeout=10)
                    
                    logger.info(
                        "Product event sent successfully",
                        extra={
                            "correlation_id": correlation_id,
                            "product_id": product_data.get("id"),
                            "topic": self.topic,
                            "attempt": attempt + 1
                        }
                    )
                    return True
                    
                except BufferError as e:
                    attempt += 1
                    last_error = e
                    if attempt < self.max_retries:
                        backoff_time = (self.retry_backoff_ms / 1000) * (2 ** attempt)
                        logger.warning(
                            "Kafka buffer full, retrying",
                            extra={
                                "correlation_id": correlation_id,
                                "attempt": attempt,
                                "max_retries": self.max_retries,
                                "backoff_seconds": backoff_time,
                                "error": str(e)
                            }
                        )
                        time.sleep(backoff_time)
                        
                except KafkaException as e:
                    attempt += 1
                    last_error = e
                    if attempt < self.max_retries:
                        backoff_time = (self.retry_backoff_ms / 1000) * (2 ** attempt)
                        logger.warning(
                            "Kafka error, retrying",
                            extra={
                                "correlation_id": correlation_id,
                                "attempt": attempt,
                                "max_retries": self.max_retries,
                                "backoff_seconds": backoff_time,
                                "error": str(e)
                            }
                        )
                        time.sleep(backoff_time)
                        
                except Exception as e:
                    span.record_exception(e)
                    logger.error(
                        "Unexpected error sending product event",
                        extra={
                            "correlation_id": correlation_id,
                            "product_id": product_data.get("id"),
                            "error": str(e),
                            "error_type": type(e).__name__
                        },
                        exc_info=True
                    )
                    raise
            
            logger.error(
                "Failed to send product event after retries",
                extra={
                    "correlation_id": correlation_id,
                    "product_id": product_data.get("id"),
                    "attempts": attempt,
                    "last_error": str(last_error)
                }
            )
            return False

    def close(self):
        try:
            self.producer.flush(timeout=10)
            logger.info("Kafka producer closed successfully")
        except Exception as e:
            logger.error(
                "Error closing Kafka producer",
                extra={"error": str(e)},
                exc_info=True
            )

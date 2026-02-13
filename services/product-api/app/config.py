from functools import lru_cache
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    app_name: str = "Product API"
    app_version: str = "1.0.0"
    environment: str = "development"
    
    kafka_bootstrap_servers: str = "localhost:9092"
    kafka_topic_products: str = "products.events"
    kafka_retries: int = 3
    kafka_retry_backoff_ms: int = 100
    kafka_acks: str = "all"
    
    otlp_endpoint: str = "http://localhost:4317"
    
    log_level: str = "INFO"
    
    cors_origins: list[str] = ["*"]
    
    enable_chaos: bool = False
    error_rate: float = 0.1
    latency_ms_min: int = 100
    latency_ms_max: int = 2000
    
    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        case_sensitive=False,
        extra="allow"
    )


@lru_cache
def get_settings() -> Settings:
    return Settings()

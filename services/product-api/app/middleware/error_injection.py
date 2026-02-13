import random
import logging
from typing import Callable
from fastapi import Request, Response, HTTPException
from starlette.middleware.base import BaseHTTPMiddleware
from starlette.status import (
    HTTP_500_INTERNAL_SERVER_ERROR,
    HTTP_503_SERVICE_UNAVAILABLE,
    HTTP_429_TOO_MANY_REQUESTS
)

logger = logging.getLogger(__name__)

ERROR_TYPES = [
    (HTTP_500_INTERNAL_SERVER_ERROR, "Internal Server Error", 0.5),
    (HTTP_503_SERVICE_UNAVAILABLE, "Service Temporarily Unavailable", 0.3),
    (HTTP_429_TOO_MANY_REQUESTS, "Too Many Requests", 0.2)
]


class ErrorInjectionMiddleware(BaseHTTPMiddleware):
    def __init__(self, app, error_rate: float = 0.0):
        super().__init__(app)
        self.error_rate = max(0.0, min(1.0, error_rate))
        
    async def dispatch(self, request: Request, call_next: Callable) -> Response:
        if request.url.path in ["/health", "/ready", "/docs", "/openapi.json"]:
            return await call_next(request)
        
        if random.random() < self.error_rate:
            error_choice = random.random()
            cumulative = 0.0
            
            for status_code, message, probability in ERROR_TYPES:
                cumulative += probability
                if error_choice <= cumulative:
                    logger.warning(
                        "Chaos: Error injected",
                        extra={
                            "path": request.url.path,
                            "method": request.method,
                            "status_code": status_code,
                            "error_message": message
                        }
                    )
                    raise HTTPException(
                        status_code=status_code,
                        detail=f"Chaos Engineering: {message}"
                    )
        
        return await call_next(request)

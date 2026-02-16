import asyncio
import random
import logging
from typing import Callable
from fastapi import Request, Response
from starlette.middleware.base import BaseHTTPMiddleware

logger = logging.getLogger(__name__)


class LatencyInjectionMiddleware(BaseHTTPMiddleware):
    def __init__(self, app, min_ms: int = 0, max_ms: int = 0):
        super().__init__(app)
        self.min_ms = max(0, min_ms)
        self.max_ms = max(self.min_ms, max_ms)
        
    async def dispatch(self, request: Request, call_next: Callable) -> Response:
        if request.url.path in ["/health", "/ready", "/docs", "/openapi.json"]:
            return await call_next(request)
        
        if self.max_ms > 0:
            latency_ms = random.randint(self.min_ms, self.max_ms)
            
            logger.info(
                "Chaos: Latency injected",
                extra={
                    "path": request.url.path,
                    "method": request.method,
                    "latency_ms": latency_ms
                }
            )
            
            await asyncio.sleep(latency_ms / 1000.0)
        
        return await call_next(request)

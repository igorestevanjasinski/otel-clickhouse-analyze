#!/usr/bin/env python3
"""
Smoke test for OpenTelemetry metrics (OTLP). Verifies module import and recording.
Metrics are exported to the collector when the app runs; this script only checks
that the module loads and record_* functions can be called (no live OTLP needed).
"""
import sys
import os

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

try:
    print("Testing OpenTelemetry metrics module...")
    from app.telemetry import otel_metrics

    print("✓ app.telemetry.otel_metrics imported")

    # Call record functions; without setup_metrics() they no-op (instruments are None)
    otel_metrics.record_request_count("success", "/products")
    otel_metrics.record_request_latency("/products", 0.1)
    otel_metrics.active_requests_inc("/products")
    otel_metrics.active_requests_dec("/products")
    otel_metrics.record_kafka_messages_published("products.events", "success")
    print("✓ record_* functions run without error (no-op when metrics not initialized)")

    print("\n✅ OpenTelemetry metrics module OK (use the running API + collector for full export)")
except ImportError as e:
    print(f"✗ Import error: {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)
except Exception as e:
    print(f"✗ Error: {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

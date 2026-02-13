FROM alpine:latest

RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

COPY bin/consumer .

RUN chown -R appuser:appuser /app

USER appuser

ENV TZ=UTC

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD pgrep -f consumer || exit 1

EXPOSE 8080

CMD ["./consumer"]

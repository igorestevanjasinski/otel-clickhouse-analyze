# Contributing Guidelines

Thank you for contributing to the Product Portfolio Observability project!

## Development Workflow

### 1. Setup Environment

```bash
# Clone and setup
git clone <repository-url>
cd otel-clickhouse-analyze

# Run automated setup
./scripts/setup.sh
```

### 2. Make Changes

- Create a feature branch: `git checkout -b feature/your-feature`
- Follow language-specific style guides (see below)
- Write tests for new functionality
- Update documentation as needed

### 3. Testing

```bash
# Run all tests
./scripts/run-tests.sh

# Run specific tests
cd services/product-api && pytest tests/ -v
cd services/product-consumer && go test ./... -v

# E2E test
./scripts/test-e2e.sh
```

### 4. Submit Changes

- Ensure all tests pass
- Coverage should be >= 70%
- Commit with clear messages
- Push and create Pull Request

## Code Style

### Python (FastAPI)

- Follow **PEP 8**
- Use **black** for formatting: `black app/`
- Use **flake8** for linting: `flake8 app/`
- Type hints preferred: `mypy app/`
- Docstrings for public functions

**Example:**
```python
def process_product(product: Product) -> bool:
    """
    Process a product and publish to Kafka.
    
    Args:
        product: Product model instance
        
    Returns:
        bool: True if successful, False otherwise
    """
    pass
```

### Go (Consumer)

- Follow **gofmt** standard: `go fmt ./...`
- Run **golint**: `golint ./...`
- Use **go vet**: `go vet ./...`
- Clear error handling (no naked returns)
- Interface-based design

**Example:**
```go
// ProductRepository defines product persistence operations
type ProductRepository interface {
    CreateProduct(ctx context.Context, product *models.Product) error
    Close() error
}
```

### Logging Standards

- **NO EMOJIS** in logs
- **JSON format** only
- Include **correlation IDs** in all log entries
- Use appropriate log levels:
  - `ERROR`: Failures requiring attention
  - `WARN`: Degraded performance or recoverable issues
  - `INFO`: Important business events
  - `DEBUG`: Detailed diagnostic information

**Example:**
```python
logger.info(
    "Product created",
    extra={
        "correlation_id": correlation_id,
        "product_id": product.id,
        "event": "product_created"
    }
)
```

## Testing Standards

### Unit Tests

- Test coverage >= 70%
- Test happy path AND edge cases
- Use mocks for external dependencies
- Clear test names: `test_create_product_success`, `test_create_product_missing_name`

### Integration Tests

- Use **testcontainers** for Go repository tests
- Clean up resources in teardown
- Skip with `-short` flag for fast feedback

### E2E Tests

- Test complete flow: API → Kafka → Consumer → Database
- Verify data persistence
- Clean up test data

## Commit Messages

Follow conventional commits:

```
feat: add chaos engineering middleware
fix: correct database connection pool settings
docs: update API endpoint documentation
test: add integration tests for repository
refactor: extract kafka producer to service layer
```

## Pull Request Checklist

- [ ] Tests pass (`./scripts/run-tests.sh`)
- [ ] Coverage >= 70%
- [ ] E2E test passes (`./scripts/test-e2e.sh`)
- [ ] Code follows style guide
- [ ] Documentation updated
- [ ] No sensitive data in code or logs
- [ ] Changelog updated (if applicable)

## Project Structure

```
otel-clickhouse-analyze/
├── services/
│   ├── product-api/          # Python FastAPI service
│   └── product-consumer/     # Go Kafka consumer
├── infrastructure/           # Docker Compose configs
├── scripts/                  # Automation scripts
├── examples/                 # Usage examples
├── sre/                      # SRE artifacts (runbooks, SLOs)
└── docs/                     # Additional documentation
```

## Common Tasks

### Add New API Endpoint

1. Define route in `services/product-api/app/api/routes.py`
2. Add Pydantic model in `app/models/`
3. Write tests in `tests/test_routes.py`
4. Update Postman collection
5. Update service README

### Add New Configuration

1. Add to `.env.example` with description
2. Update `config.py` or `config.go`
3. Document in service README
4. Update setup script if needed

### Add New Middleware

1. Create file in `app/middleware/`
2. Register in `app/main.py`
3. Add tests
4. Document in README

## Getting Help

- Check existing documentation in READMEs
- Review similar implementations in codebase
- Ask questions in Pull Request comments

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

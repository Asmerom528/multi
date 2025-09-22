all: lint test bench

lint:
	@golangci-lint run ./...
	@echo "✓ lint"

test:
	@go test ./...
	@echo "✓ test"

bench:
	@go run github.com/nalgeon/multi/internal/benchmark

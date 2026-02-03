# go-goog-cli

A Golang CLI application built with Clean Architecture and TDD practices.

## Quick Start

```bash
# Build
go build -o bin/go-goog-cli ./cmd/go-goog-cli

# Run
./bin/go-goog-cli

# Run tests
go test ./...
```

## Project Structure

- `cmd/` - Application entry points
- `internal/` - Private application code (domain, usecase, adapter, infrastructure layers)
- `documentation/` - Product and technical documentation

## Development

See [AGENTS.md](AGENTS.md) for development guidelines, architecture details, and workflows.

## Documentation

- [Product Summary](documentation/product-summary.md) - What this product does
- [Product Details](documentation/product-details.md) - Features and user workflows
- [Technical Details](documentation/technical-details.md) - Architecture and implementation

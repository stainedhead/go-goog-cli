# Technical Details

## Architecture

This application follows Clean Architecture with four layers:

1. **Domain** (`internal/domain/`) - Business entities and rules
2. **Use Case** (`internal/usecase/`) - Application-specific business logic
3. **Adapter** (`internal/adapter/`) - Interface adapters for CLI and data access
4. **Infrastructure** (`internal/infrastructure/`) - External service implementations

See [AGENTS.md](../AGENTS.md) for detailed architectural guidelines.

## Technology Stack

- **Language**: Go 1.21+
- **Architecture**: Clean Architecture
- **Testing**: Go standard testing package
- **Linting**: golangci-lint

## Data Flow

*To be documented as components are implemented.*

## API Documentation

*To be documented as APIs are defined.*

## Configuration

*To be documented as configuration options are added.*

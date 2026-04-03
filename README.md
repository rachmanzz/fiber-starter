# Fiber V3 Boilerplate

A production-ready Go backend boilerplate using Fiber v3, following clean architecture principles with repository and service patterns.

## Directory Layout

```
project-root/
├── cmd/                          # Application entrypoints
│   ├── api/
│   │   └── main.go              # API server entrypoint
│   └── cron/
│       └── main.go              # Cron  entrypoint
│
├── bootstrap/
│   ├── app.go              # Main application bootstrap
│   ├── database.go         # Database initialization
│   ├── server.go           # HTTP server initialization
│   ├── providers.go        # Provider registration
│   └── health.go           # Health check aggregation
│
├── providers/                     # Custom providers directory
│   ├── provider.go               # Provider interface and types
│   ├── cache/                    # Cache provider example
│   │   ├── cache.go             # Cache provider implementation
│   │   └── config.go            # Cache configuration
│   ├── queue/                    # Queue provider example
│   │   ├── queue.go             # Queue provider implementation
│   │   └── config.go            # Queue configuration
│   └── logger/                   # Logger provider example
│       ├── logger.go            # Logger provider implementation
│       └── config.go            # Logger configuration
│
├── config/                       # Configuration
│   ├── config.go                # Config loader
│   ├── database.go              # Database configuration
│   └── env.go                   # Environment variables
│
├── app/                          # Application layer
│   ├── delivery/                 # HTTP handlers (delivery layer)
│   ├── dto/                      # Data Transfer Objects
│   │   ├── request/             # Request DTOs
│   │   └── response/            # Response DTOs
│   ├── repository/               # Data access layer
│   │   ├── repository.go        # Repository interface
│   │   ├── user_repository.go   # User data access
│   │   └── {entity}_repository.go
│   │
│   ├── services/                 # Business logic layer
│   │   ├── service.go           # Service interface
│   │   ├── user_service.go      # User business logic
│   │   └── {entity}_service.go
│   │
│   ├── route/                    # Route definitions
│   │   └── api.go               # API routes registration
│   │
│   └── middleware/               # HTTP middleware
│       ├── auth.go              # Authentication middleware
│       ├── logger.go            # Logging middleware
│       └── cors.go              # CORS middleware
│
├── cores/                        # Core features
│   ├── database/              # Database core (non-editable)
│   │   ├── pool.go           # Connection pool management
│   │   ├── monitoring.go     # Database monitoring
│   │   ├── transaction.go    # Transaction handling
│   │   └── context.go        # Database context utilities
│   ├── logging/              # Logging core (non-editable)
│   │   ├── logger.go         # Zap logger initialization
│   │   ├── config.go         # Logger configuration
│   │   ├── level.go          # Log level definitions
│   │   └── output.go         # Output destinations (terminal, file, custom)
│   ├── response/             # Response core (non-editable)
│   │   ├── response.go       # Core response implementation
│   │   ├── success.go        # Success response builders
│   │   ├── error.go          # Error response builders
│   │   └── codes.go          # HTTP status code mappings
│   ├── errors/              # Core error definitions (non-editable)
│   │   ├── errors.go       # Core error types and constructors
│   │   ├── codes.go        # Error code constants
│   │   ├── sql.go          # SQL error mapping
│   │   ├── validation.go   # Validation error builders
│   │   └── checker.go      # Error checking utilities
│   └── config/              # Core config structures (non-editable)
│       └── config.go        # Config struct definitions
│
├── utils/                        # Utility functions
│   ├── response.go              # Response helpers
│   ├── validator.go             # Validation helpers
│   └── hash.go                  # Hashing utilities
│
├── db/                           # Generated database code (sqlc)
│   ├── models.go                # Generated models
│   ├── queries.sql.go           # Generated queries
│   └── db.go                    # DB connection
│
├── queries/                      # SQL queries for sqlc
│   ├── users.sql                # User queries
│   └── {entity}.sql             # Entity queries
│
├── migrations/                   # Database migrations (tern)
│   ├── 001_init.sql             # Initial schema
│   ├── 002_create_users.sql     # Users table
│   └── tern.conf                # Tern configuration
│
├── scripts/                      # Build and deployment scripts
│   ├── build.sh                 # Build script
│
├── test/                         # Integration tests
│   ├── integration/             # Integration tests
│   └── e2e/                     # End-to-end tests
│
├── go.mod                        # Go module definition
├── go.sum                        # Go dependencies checksum
├── sqlc.yaml                     # SQLC configuration
├── Makefile                      # Make commands
├── Dockerfile                    # Docker configuration
├── docker-compose.yaml          # Docker compose
└── README.md                     # Project documentation
```

## Features

- **Fiber v3** - Fast and lightweight HTTP framework
- **Clean Architecture** - Repository and service patterns
- **Database** - PostgreSQL with sqlc for type-safe queries
- **Migrations** - Tern for database migrations
- **Logging** - Zap logger with configurable outputs
- **Error Handling** - Centralized error handling with custom error types
- **Response Format** - Consistent API response structure
- **Middleware** - Pre-built auth, CORS, and logging middleware
- **Provider System** - Extensible provider architecture for cache, queue, and logging
- **Health Checks** - Aggregated health check system
- **Docker Support** - Dockerfile and docker-compose configuration

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd boilerplate-fiber-v3
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Generate database code:
```bash
make sqlc
```

5. Run migrations:
```bash
make migrate
```

## Usage

### Development

Run the API server:
```bash
make run
```

Run with Docker:
```bash
docker-compose up --build
```

### Available Commands (Makefile)

- `make run` - Run the API server
- `make build` - Build the application
- `make sqlc` - Generate database code
- `make migrate` - Run database migrations
- `make test` - Run tests
- `make lint` - Run linter
- `make docker` - Build Docker image

### API Endpoints

The boilerplate includes example routes. Check `app/route/api.go` for the default route configuration.

## License

MIT License
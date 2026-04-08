# Fiber V3 Boilerplate

A production-ready Go backend boilerplate using Fiber v3, following clean architecture principles with repository and service patterns.

## Directory Layout

```
project-root/
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
├── bootstrap/                    # Application bootstrap
│   └── app.go                    # App lifecycle & graceful shutdown logic
│
├── cmd/                          # Application entrypoints
│   └── server/
│       └── main.go               # Main API server entrypoint
│
├── config/                       # Configuration files
│   ├── app.go                    # Application specific config
│   ├── config.go                 # Main config loader
│   ├── database.go               # Database connection settings
│   └── logger.go                 # Logging settings
│
├── cores/                        # Core framework components
│   ├── config.go                 # Core configuration structures
│   ├── contract.go               # App contract & hook registration
│   ├── database.go               # Database connection pool (pgx)
│   └── logger.go                 # Zap logger core initialization
│
├── spark-cli/                    # Spark CLI source code
│   └── main.go                   # CLI implementation
│
├── spark                         # Spark CLI binary
├── logs/                         # Application logs
├── tmp/                          # Temporary build files
├── go.mod                        # Go module definition
├── go.sum                        # Go dependencies checksum
└── README.md                     # Project documentation
```

## Features

- **Fiber v3** - Fast and lightweight HTTP framework
- **Spark CLI** - Custom CLI for dev server, migrations, and initialization
- **Clean Architecture** - Repository and service patterns
- **Database** - PostgreSQL with pgx/v5 for high-performance pooling
- **Migrations** - Integrated Tern support for database migrations
- **Lifecycle Management** - Built-in Graceful Shutdown & Hook system (Before/After)
- **Logging** - Production-grade Zap logger with configurable outputs
- **Environment Driven** - Configuration via environment variables
- **Response Format** - Consistent API response structure

## Installation

1. Clone the repository using degit (recommended for a clean start) or git:
```bash
# Using degit
npx degit rachmanzz/fiber-starter my-project
cd my-project

# Or using git
git clone https://github.com/rachmanzz/fiber-starter.git my-project
cd my-project
```

2. Initialize the project (change module name):
```bash
# This will automatically update the module name in all files and go.mod
./spark init
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

## Usage

### Development

Run the API server with live reloading (Air):
```bash
./spark dev
```

### Database Migrations

Spark uses **Tern** for migrations. It will automatically install tern and initialize `tern.conf` if needed.

Run migrations:
```bash
./spark migrate
```

Create a new migration:
```bash
./spark migrate new create_users_table
```

### Database Queries (SQLC)

This project uses **SQLC** for type-safe database access.

1.  **Define Queries**: Place your SQL query files in the `queries/` directory at the root (e.g., `queries/users.sql`).
2.  **Configuration**: To link your queries with the repository layer, ensure you have a `sqlc.yaml` in the root with the following configuration:

```yaml
version: "2"
sql:
  - schema: "migrations" # or your schema path
    queries: "queries"
    engine: "postgresql"
    gen:
      go:
        package: "repository"
        out: "app/repository"
        sql_package: "pgx/v5"
```

3.  **Generate Code**: Run the following command to generate the Go code:
```bash
sqlc generate
```
The generated files will be placed in `app/repository/`, making them ready to be used by your services.

### Testing

Fiber v3 is designed to be easily testable. This boilerplate follows Go's standard testing patterns combined with Fiber's built-in testing utilities.

#### Unit Testing

Fiber v3 provides the `app.Test` method to simulate HTTP requests without starting a network server. This makes tests extremely fast.

**Best Practices:**
- Use `github.com/stretchr/testify/assert` for idiomatic assertions.
- Use **Dependency Injection** to pass mock repositories or services into your handlers.
- Use **Table-Driven Tests** for comprehensive coverage.

**Example Unit Test:**
```go
func TestUserHandler(t *testing.T) {
    // 1. Setup
    app := fiber.New()
    mockRepo := new(MockUserRepository)
    app.Get("/users/:id", handlers.GetUser(mockRepo))

    // 2. Create Request
    req, _ := http.NewRequest("GET", "/users/1", nil)

    // 3. Perform Test
    // In v3, Test takes (req, timeout_ms). Use -1 for no timeout.
    resp, err := app.Test(req, -1)

    // 4. Assertions
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
}
```

#### Integration Testing

For integration tests that require a real database, it is recommended to:
1.  Use a dedicated **test database**.
2.  Run migrations before tests using `./spark migrate`.
3.  Optionally use **Testcontainers** to spin up temporary PostgreSQL instances in Docker.

To run all tests:
```bash
go test ./...
```

### Available Commands (Spark CLI)

The Spark CLI is provided as a pre-built binary (`spark`). To audit the source code or build it yourself, please refer to the [Spark Resource (with-spark-cli branch)](https://github.com/rachmanzz/fiber-starter/tree/with-spark-cli), as the `spark-cli/` directory is removed from the main branch to keep the boilerplate lean.

- `spark init` - Initialize the project with a new module name.
- `spark dev` - Run the application with live reloading using Air.
- `spark migrate` - Run all pending database migrations.
- `spark migrate --to [version]` - Migrate to a specific version.
- `spark migrate new [name]` - Create a new migration file.
- `spark version` - Print the version of Spark.

### Docker Support

Run with Docker:
```bash
docker-compose up --build
```

### API Endpoints

The boilerplate includes example routes. Check `app/route/api.go` for the default route configuration.

## License

MIT License

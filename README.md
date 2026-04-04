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
- **Error Handling** - Centralized error handling structure
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

### Available Commands (Spark CLI)

- `./spark init` - Initialize the project with a new module name
- `./spark dev` - Run the application with live reloading using Air
- `./spark migrate` - Run database migrations
- `./spark migrate new [name]` - Create a new migration file
- `./spark version` - Print the version number of Spark

### Docker Support

Run with Docker:
```bash
docker-compose up --build
```

### API Endpoints

The boilerplate includes example routes. Check `app/route/api.go` for the default route configuration.

## License

MIT License
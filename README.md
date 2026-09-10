# Fiber V3 Boilerplate

A minimalist and high-performance Go backend boilerplate using **Fiber v3**. Designed to be lightweight, environment-driven, and easy to extend.

## Directory Layout

```text
.
├── app/
│   ├── repository/       # Data access layer (PostgreSQL with pgx/v5)
│   └── routes/           # Route definitions and registrations
├── bootstrap/            # Application lifecycle (Bootstrap, Middleware & Graceful Shutdown)
├── cmd/server/           # Application entrypoint (main.go)
├── config/               # Configuration loaders (App, DB, Logger)
├── cores/                # Core Framework components
│   ├── config.go         # Config structures
│   ├── contract.go       # Fiber instance & Hook management
│   ├── database.go       # DB connection pool (pgx)
│   ├── logger.go         # Zap logger initialization
│   └── response.go       # Standardized API response helpers
├── docs/                 # Layer & tooling documentation
├── spark-cli/            # Source code of the Spark CLI
└── .env.example          # Environment template
```

## Features

- **Fiber v3** - Leveraging the latest features of the Fiber framework.
- **Spark CLI** - Custom tool for project initialization, migrations, and live-reloading.
- **Dual Response Format** - Built-in support for **JSON** and **MessagePack** (via `Accept` header).
- **Graceful Shutdown** - Handles OS signals to close DB connections and stop the server safely.
- **Lifecycle Hooks** - Register "Before" and "After" hooks for setup/teardown logic.
- **Structured Logging** - High-performance logging using **Uber Zap**.
- **PostgreSQL Ready** - Pre-configured connection pooling using `pgx/v5`.

## Getting Started

### 1. Installation

Clone the repository and enter the directory:

```bash
git clone https://github.com/rachmanzz/fiber-starter.git my-project
cd my-project
```

### 2. Build the Spark CLI

Build the Spark CLI binary from source (the result is already executable, no `chmod` needed):

```bash
# Linux / macOS (runs on your current platform)
go build -o spark ./spark-cli
```

Cross-compiling for another OS:

```bash
# Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o spark ./spark-cli

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o spark ./spark-cli

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o spark ./spark-cli

# Windows
GOOS=windows GOARCH=amd64 go build -o spark.exe ./spark-cli
```

> **Note (Windows):** the binary is `spark.exe` — run it as `spark init`, without the `./` prefix.

### 3. Initialization

Use the Spark CLI to rename the module to your own:

```bash
./spark init
```
*This will interactively ask for your module name and update all imports automatically.*

### 4. Environment Setup

```bash
cp .env.example .env
# Edit .env with your database credentials and app port
```

### 5. Running the App

For development with **live-reloading** (requires [Air](https://github.com/air-verse/air)):

```bash
./spark dev
```

To run normally:

```bash
go run cmd/server/main.go
```

### Adding New Routes

Define your routes in `app/routes/api.go`. They are automatically loaded during the bootstrap process in `bootstrap/app.go`.

## Spark CLI Commands

The `spark` binary is a helper tool for common tasks. Build it first with `go build -o spark ./spark-cli`, or use `go run ./spark-cli <command>` directly.

- `spark init` - Initialize project and rename module.
- `spark dev` - Run development server with live-reloading.
- `spark migrate` - Run database migrations using [Goose](https://github.com/pressly/goose).
- `spark migrate down` - Roll back the most recently applied migration.
- `spark migrate new [name]` - Create a new migration file in `/migrations`.
- `spark version` - Show CLI version.

## License

MIT License

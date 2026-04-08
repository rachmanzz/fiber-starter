# Fiber V3 Boilerplate

A minimalist and high-performance Go backend boilerplate using **Fiber v3**. Designed to be lightweight, environment-driven, and easy to extend.

## Directory Layout

```text
.
├── app/
│   ├── repository/       # Data access layer (PostgreSQL with pgx/v5)
│   └── routes/           # Route definitions and registrations
├── bootstrap/            # Application lifecycle (Bootstrap & Graceful Shutdown)
├── cmd/server/           # Application entrypoint (main.go)
├── config/               # Configuration loaders (App, DB, Logger)
├── cores/                # Core Framework components
│   ├── config.go         # Config structures
│   ├── contract.go       # Fiber instance & Hook management
│   ├── database.go       # DB connection pool (pgx)
│   ├── logger.go         # Zap logger initialization
│   └── response.go       # Standardized API response helpers
├── spark                 # Spark CLI binary
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

### 2. Initialization

Use the Spark CLI to rename the module to your own:

```bash
./spark init
```
*This will interactively ask for your module name and update all imports automatically.*

### 3. Environment Setup

```bash
cp .env.example .env
# Edit .env with your database credentials and app port
```

### 4. Running the App

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

The `spark` binary is a helper tool for common tasks. 
*(Source code for the CLI is maintained in the `with-spark-cli` branch for auditing/custom builds).*

- `spark init` - Initialize project and rename module.
- `spark dev` - Run development server with live-reloading.
- `spark migrate` - Run database migrations using [Tern](https://github.com/jackc/tern).
- `spark migrate new [name]` - Create a new migration file in `/migrations`.
- `spark version` - Show CLI version.

## License

MIT License

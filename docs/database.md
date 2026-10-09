# Database Documentation

This document explains how the database and repository layers are structured in **fiber-starter**.

## Overview

The database connection is managed centrally in `cores/database.go` using **pgxpool (v5)**. The core database module supports both single default connection pools and **named multi-database pools** (e.g. for separating transactional and analytical/event data).

Access to the database is abstracted through the **Repository** pattern to ensure type safety and clean architecture.

---

## Environment Configuration

Configure the database connection in `.env` (copied from `.env.example`):

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=your_database
DB_SSLMODE=disable
DB_ENABLE=true
DB_MAX_CONNS=10
DB_MIN_CONNS=2
DB_MAX_CONN_LIFETIME=1h
DB_MAX_CONN_IDLE_TIME=30m
```

The `DatabaseConfig` struct in `config/database.go` exposes a helper method `.DSN()` that automatically constructs the PostgreSQL connection string.

---

## Connection Management (`cores/database.go`)

> **⚠️ Strict Contract Rule**: `cores` does **not** export `GetDB()`. Direct database pool access from handlers or services is strictly prohibited. Database pools must only be accessed through the **Repository Contract** (`contract.GetQueries()`).

### 1. Default Connection Pool
For standard single-database setups:

- **`cores.ConnectDB() error`**: Connects to the default PostgreSQL database using `Config().Database`. The default pool is registered under the identifier `"default"` (`cores.DefaultDatabaseName`).

### 2. Multi-Database (Named Pools)
For applications connecting to multiple databases (e.g. secondary event/analytics store):

- **`cores.ConnectNamedDB(name string, cfg config.DatabaseConfig) (*pgxpool.Pool, error)`**: Initializes and registers a named connection pool.
- **`cores.RegisterDBPool(name string, pool *pgxpool.Pool)`**: Manually registers an existing pool pointer under a specific name.
- **`cores.HasDB(name string) bool`**: Checks if a named database pool exists and is active.
- **`cores.CloseNamedDB(name string) error`**: Safely closes a single named connection pool.
- **`cores.CloseDB()`**: Safely closes **all** active connection pools with panic recovery during application shutdown.


---

## Lifecycle & Hook Integration (`bootstrap/hook.go`)

Database connections and repository contracts are wired inside `bootstrap/hook.go` using application lifecycle hooks:

### 1. Pre-Start Hook (`RegisterBeforeStart`)
Database connections and repository contracts are initialized **before** the server begins listening for traffic:

```go
// bootstrap/hook.go
core.RegisterBeforeStart(func() error {
    if cores.Config().Database.Enable {
        cores.SetNamedDatabaseContract(func(name string, pool *pgxpool.Pool) {
            // contract.RegisterNamedDatabase(name, pool) // Automatically registers default and any named pool
        })

        if err := cores.ConnectDB(); err != nil {
            return fmt.Errorf("failed to connect to database: %w", err)
        }
    }

    return nil
})
```

> **Contract Backfill & Deadlock Safety**: `SetDatabaseContract` and `SetNamedDatabaseContract` invoke their callbacks **after releasing internal mutex locks**, completely eliminating re-entrant deadlocks. If a contract is registered after a database pool has already connected, it immediately backfills the existing pool(s).

### 2. Teardown Hook (`OnPostShutdown`)
When the application receives termination signals (`SIGINT`, `SIGTERM`), Fiber v3 triggers `OnPostShutdown` to gracefully close all active connection pools:

```go
// bootstrap/hook.go
core.App.Hooks().OnPostShutdown(func(err error) error {
    if cores.Config().Database.Enable {
        cores.CloseDB()
        zap.L().Info("Database connection pool closed successfully")
    }
    return nil
})
```

---

## Accessing the Database via Repository Contract

Handlers and services should never interact with database connection pools directly. Instead, use the registered repository contract located in `app/repository/contract/registry.go`.

### Contract API Reference

| Function | Description |
| :--- | :--- |
| `contract.GetQueries()` | Returns the default `*repository.Queries` instance (panics if not initialized). |
| `contract.Use("pointer_name")` | Retrieves queries for a named database (e.g. `contract.Use("analytics")`). Defaults to `"default"` if omitted. |
| `contract.GetQueriesNamed(name)` | Retrieves queries for a specific named pointer with an error return. |
| `contract.Has(name)` | Checks if a named database repository is registered. |
| `contract.GetPool(name...)` | Returns raw `*pgxpool.Pool` for transactions (`pool.Begin(ctx)`) or third-party drivers. |

### Example Service Usage

```go
// app/service/user_service.go
func (s *UserService) GetUser(ctx context.Context, id int64) (*repository.User, error) {
    // 1. Default database queries
    return contract.GetQueries().GetUserByID(ctx, id)
}

func (s *UserService) LogAuditEvent(ctx context.Context, event repository.CreateEventParams) error {
    // 2. Secondary named database queries (e.g. analytics/events)
    return contract.Use("events").CreateEvent(ctx, event)
}
```

---

## Database Queries (SQLC)

This project uses **SQLC** for compile-time type-safe database queries.

### 1. Installation

**Go tool:**
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

**macOS (Homebrew):**
```bash
brew install sqlc
```

### 2. How to use SQLC

1. **Define Queries**: Place your SQL query files in the `queries/` directory at the project root (e.g., `queries/users.sql`).
   
2. **Configuration (`sqlc.yaml`)**:
   ```yaml
   version: "2"
   sql:
     - schema: "migrations"
       queries: "queries"
       engine: "postgresql"
       gen:
         go:
           package: "repository"
           out: "app/repository"
           sql_package: "pgx/v5"
   ```

3. **Generate Code**:
   ```bash
   sqlc generate
   ```

The generated code will populate `app/repository/`. Once generated, uncomment `contract.DatabaseContract(pool)` in `bootstrap/hook.go`.

# Database Documentation

This document explains how the database and repository layers are structured in this boilerplate.

## Overview

The database connection is managed centrally in `cores/database.go` using **pgxpool (v5)**. Access to the database is abstracted through the **Repository** pattern to ensure type safety and clean architecture.

## Environment Requirement

Make sure these variables exist in your .env file:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=your_database
DB_SSLMODE=disable
DB_ENABLE=true
```
## Accessing the Database

You should never access the database connection directly in your handlers or services. Instead, use the registered repository contract.

### The Repository Contract

The contract is located in `app/repository/contract/registry.go`. It provides a centralized point to get the initialized queries.


### Lifecycle & Registration

The database initialization follows this flow:
- **Registration**: In `bootstrap/db.go`, we register how the database pool should be "contracted" to the repository layer.

> **Note**: By default, the database contract registration in `bootstrap/db.go` is commented out. This is because the boilerplate doesn't come with pre-generated SQLC code. Once you have generated your repository code, you should uncomment it:
> ```go
> // bootstrap/db.go
> func RegisterDatabaseContract() {
>     cores.SetDatabaseContract(func(pool *pgxpool.Pool) {
>         contract.DatabaseContract(pool) // Uncomment this
>     })
> }
> ```

## Database Queries (SQLC)

This project uses **SQLC** for type-safe database access. To keep things simple and avoid "magic" that hides the implementation, SQLC is integrated manually.

### Installation

Before using SQLC, you need to install it on your system. 

**Using Go:**
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

**macOS (Homebrew):**
```bash
brew install sqlc
```

For other platforms, refer to the [official SQLC installation guide](https://docs.sqlc.dev/en/latest/overview/install.html).

### How to use SQLC

1. **Define Queries**: Place your SQL query files in the `queries/` directory at the root (e.g., `queries/users.sql`).
   
2. **Configuration**: Ensure you have a `sqlc.yaml` file in the root with the following configuration:

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

3. **Generate Code**: Run the following command to generate the Go code:
   ```bash
   sqlc generate
   ```

The generated files will be placed in `app/repository/`. Once generated, they will be automatically available through `contract.GetQueries()`.

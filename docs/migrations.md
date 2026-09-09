# Database Migrations

This project uses [Goose](https://github.com/pressly/goose) for managing database migrations. Goose is a widely used, actively maintained migration tool that supports plain SQL migrations with `-- +goose Up` / `-- +goose Down` annotations.

Migration files live in the `migrations/` directory at the project root.

## Using Spark CLI for Migrations

The `spark` CLI wraps goose and builds the database DSN automatically from your `.env` file (`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`), so you never have to pass connection strings by hand.

- `./spark migrate` - Apply all pending migrations (`goose up`).
- `./spark migrate --to [version]` - Migrate up to a specific version.
- `./spark migrate down` - Roll back the most recently applied migration.
- `./spark migrate new [name]` - Create a new SQL migration file in `/migrations`.

If goose is not installed, spark will install it automatically via `go install github.com/pressly/goose/v3/cmd/goose@latest`.

## Using Goose Directly

You can also run goose yourself. Build the DSN from your `.env` values:

```bash
goose -dir migrations postgres "postgres://user:password@localhost:5432/your_database?sslmode=disable" up
goose -dir migrations postgres "postgres://user:password@localhost:5432/your_database?sslmode=disable" down
goose -dir migrations create add_users_table sql
```

## Migration File Format

Each migration is a SQL file with goose annotations:

```sql
-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE users;
```

For more commands and options, refer to the official [goose documentation](https://github.com/pressly/goose#usage).

# Database Migrations

This project uses `tern` for managing database migrations. `tern` is a standalone migration tool that provides a simple and effective way to evolve your database schema.

## Using Tern

To use `tern` for migrations, you need to ensure that the `TERN_MIGRATIONS` environment variable is set to the directory where your migration files are located. In this project, migration files are expected to be in a `migrations` directory.

You can set this environment variable in your shell before running `tern` commands:

```bash
export TERN_MIGRATIONS=migrations
tern migrate
```

This command will apply any pending migrations found in the `migrations` directory.

For more information on `tern` commands, please refer to the official `tern` documentation.

## Using Spark CLI for Migrations

The project's `spark` CLI tool also provides a convenient wrapper for running migrations using `tern`. This simplifies the process by abstracting away the need to manually set environment variables.

To run migrations using `spark`, you can use the following command:

```bash
./spark migrate
```

This command will execute the `tern migrate` command with the appropriate `TERN_MIGRATIONS` environment variable already configured by `spark`.

Here's a list of `spark` migration commands:

- `./spark migrate` - Run database migrations.
- `./spark migrate --to [version]` - Rollback to a specific migration version.
- `./spark migrate new [name]` - Create a new migration file.
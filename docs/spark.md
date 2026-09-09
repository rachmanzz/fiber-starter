# Spark CLI Documentation

The Spark CLI is an optional feature designed to streamline common development tasks within this project.

## Building Your Own Spark CLI

The source code of the Spark CLI lives in the [`spark-cli/`](../spark-cli) directory of this repository. You have the flexibility to build and customize your own Spark CLI from it:

```bash
# Linux / macOS (runs on your current platform)
go build -o spark ./spark-cli

# Windows
GOOS=windows GOARCH=amd64 go build -o spark.exe ./spark-cli
```

The build output is already executable (`spark.exe` on Windows) — no `chmod` required. Use `GOOS`/`GOARCH` to cross-compile for other platforms.

This allows you to tailor the CLI to your specific project needs or audit its functionality.

## Spark CLI Commands

The `spark` binary is a helper tool for common tasks within this project. Here's a list of available commands:

-   `spark init`: Initializes the project and renames the module. This command will interactively ask for your module name and update all imports automatically.
-   `spark dev`: Runs the development server with live-reloading. This requires [Air](https://github.com/air-verse/air) to be installed.
-   `spark migrate`: Runs database migrations using [Goose](https://github.com/pressly/goose).
-   `spark migrate down`: Rolls back the most recently applied migration.
-   `spark migrate new [name]`: Creates a new migration file in the `/migrations` directory with the specified name.
-   `spark version`: Displays the current version of the Spark CLI.
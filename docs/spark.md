# Spark CLI Documentation

The Spark CLI is an optional feature designed to streamline common development tasks within this project.

## Building Your Own Spark CLI

You have the flexibility to build and customize your own Spark CLI by following the codebase in this repository branch:
[https://github.com/rachmanzz/fiber-starter/tree/with-spark-cli](https://github.com/rachmanzz/fiber-starter/tree/with-spark-cli)

This allows you to tailor the CLI to your specific project needs or audit its functionality.

## Spark CLI Commands

The `spark` binary is a helper tool for common tasks within this project. Here's a list of available commands:

-   `spark init`: Initializes the project and renames the module. This command will interactively ask for your module name and update all imports automatically.
-   `spark dev`: Runs the development server with live-reloading. This requires [Air](https://github.com/air-verse/air) to be installed.
-   `spark migrate`: Runs database migrations using [Tern](https://github.com/jackc/tern).
-   `spark migrate new [name]`: Creates a new migration file in the `/migrations` directory with the specified name.
-   `spark version`: Displays the current version of the Spark CLI.
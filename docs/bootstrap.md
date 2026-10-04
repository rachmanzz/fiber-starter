# Application Lifecycle & Bootstrap Documentation

The `bootstrap` package manages the initialization, middleware registration, lifecycle hooks, and server lifecycle of the application.

## Lifecycle Hooks Overview

Application hooks are configured in `bootstrap/hook.go` using the `core` contract instance:

### 1. Pre-Start Hooks (`RegisterBeforeStart`)
Pre-start hooks run **before** the server opens the network socket listener (`app.Listen()`). If any pre-start hook returns an `error`, the startup process is aborted immediately, preventing traffic on an unready server.

```go
core.RegisterBeforeStart(func() error {
    // Perform pre-flight validation (e.g., Cache Warmup, External Service Health Check)
    if err := validateDependencies(); err != nil {
        return err // Aborts server startup
    }
    return nil
})
```

### 2. Teardown Hooks (`OnPostShutdown`)
Teardown hooks utilize Fiber v3's native lifecycle hooks (`app.Hooks().OnPostShutdown`) to gracefully close background connections (e.g., PostgreSQL pool) and flush log buffers when the server receives termination signals (`SIGINT`/`SIGTERM`).

```go
core.App.Hooks().OnPostShutdown(func(err error) error {
    if cores.Config().Database.Enable {
        cores.CloseDB()
        zap.L().Info("Database connection pool closed successfully")
    }
    return nil
})
```

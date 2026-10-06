# External Provider Architecture Documentation

This document explains the guidelines and conventions for writing and registering **External Feature Providers** in `fiber-starter`.

---

## 1. Architectural Philosophy

In this boilerplate structure:
* **`cores/`**: Strictly reserved for internal, generic framework utilities (Logger, Response formatter, Validator, Global Error Handler). It remains free of external business or infrastructure dependencies.
* **`app/`**: Contains application-specific business logic (Services, Repositories, Routes, Handlers, Middlewares).
* **`providers/`**: Placed at the project root to house integrations with external systems and third-party infrastructure:
  * `providers/redis/` (Cache, distributed locks, pub/sub)
  * `providers/nats/` (NATS JetStream event streaming)
  * `providers/mail/` (Email delivery: SMTP, SES, Sendgrid)
  * `providers/payment/` (Payment gateways: Duitku, Stripe)
  * `providers/storage/` (Object storage: S3, MinIO)
* **`bootstrap/`**: Wires the lifecycle of each provider to application start and graceful shutdown hooks (`RegisterBeforeStart` and `OnPostShutdown`).

---

## 2. How to Write a Provider

Each provider lives in its own dedicated package under `providers/<provider_name>/`.

### Standard Provider Structure
A provider package should expose three core components:
1. **`Connect()` / `Init()`**: Thread-safe initialization using `sync.Once` or `sync.Mutex`.
2. **Getter Function (`Client()`)**: Exposes the initialized client instance.
3. **`Close()`**: Performs graceful cleanup and connection teardown.

### Example Implementation: `providers/redis/redis.go`

```go
package redisprovider

import (
	"context"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/rachmanzz/fiber-starter/cores"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	client     *redis.Client
	clientOnce sync.Once
)

// Connect initializes the global Redis client using cores.Config().Redis.
func Connect() *redis.Client {
	clientOnce.Do(func() {
		cfg := cores.Config().Redis
		addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))

		rdb := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: cfg.Password,
			DB:       cfg.DB,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := rdb.Ping(ctx).Err(); err != nil {
			zap.L().Fatal("Redis connection failed", zap.Error(err), zap.String("addr", addr))
		}

		zap.L().Info("Redis provider connected", zap.String("addr", addr), zap.Int("db", cfg.DB), zap.String("prefix", cfg.Prefix))
		client = rdb
	})
	return client
}

// Client returns the active Redis client instance, or nil if not connected.
func Client() *redis.Client {
	return client
}

// Close gracefully closes the Redis client connection.
func Close() error {
	if client != nil {
		err := client.Close()
		if err != nil {
			zap.L().Warn("Error closing Redis connection", zap.Error(err))
			return err
		}
		zap.L().Info("Redis provider connection closed")
		client = nil
	}
	return nil
}
```

---

## 3. Configuring the Provider

Define the provider's environment variables in `config/`:

### 1. Create `config/redis.go`
```go
package config

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" envDefault:"localhost"`
	Port     int    `env:"REDIS_PORT" envDefault:"6379"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB" envDefault:"0"`
	Prefix   string `env:"REDIS_PREFIX" envDefault:"app:"`
	Enable   bool   `env:"REDIS_ENABLE" envDefault:"false"`
}
```

### 2. Register in `config/config.go`
```go
package config

type ConfigRegistry struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig // Add here
	Log      LoggerConfig
}
```

### 3. Add to `.env.example`
```env
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_PREFIX=app:
REDIS_ENABLE=false
```

---

## 4. How to Register a Provider via Lifecycle Hooks

Registration is handled cleanly in **`bootstrap/hook.go`** without polluting `app.go`:

* **`RegisterBeforeStart` (Pre-start Hook)**: Establishes the provider connection **before** the HTTP listener opens (`app.Listen()`). If the connection fails, startup aborts immediately, preventing the server from accepting traffic with broken dependencies.
* **`OnPostShutdown` (Teardown Hook)**: Gracefully closes connections when the application receives termination signals (`SIGINT`/`SIGTERM`).

### Implementation in `bootstrap/hook.go`

```go
package bootstrap

import (
	"github.com/rachmanzz/fiber-starter/cores"
	redisprovider "github.com/rachmanzz/fiber-starter/providers/redis"
	"go.uber.org/zap"
)

func RegisterHook(core *cores.AppContracts) {
	if core.App == nil {
		return
	}

	// 1. PRE-START: Connect provider before opening HTTP listener
	core.RegisterBeforeStart(func() error {
		if cores.Config().Redis.Enable {
			redisprovider.Connect()
		}
		return nil
	})

	// 2. POST-SHUTDOWN: Gracefully disconnect on termination
	core.App.Hooks().OnPostShutdown(func(err error) error {
		// Close database connection
		if cores.Config().Database.Enable {
			cores.CloseDB()
			zap.L().Info("Database connection pool closed successfully")
		}

		// Close Redis provider connection
		if cores.Config().Redis.Enable {
			_ = redisprovider.Close()
		}

		return nil
	})
}
```

---

## 5. Consuming Providers in the Service Layer

Application services in `app/services/` import the provider package and call `Client()`:

```go
package services

import (
	"context"
	redisprovider "github.com/rachmanzz/fiber-starter/providers/redis"
)

type ExampleService struct{}

func (s *ExampleService) SetCache(ctx context.Context, key, val string) error {
	rdb := redisprovider.Client()
	if rdb != nil {
		return rdb.Set(ctx, key, val, 0).Err()
	}
	return nil
}
```

---

## 6. Quick Checklist for New Providers

1. [ ] Create `config/<name>.go` and register in `config/config.go`.
2. [ ] Add environment variable keys and defaults to `.env.example`.
3. [ ] Create `providers/<name>/<name>.go` with `Connect()`, `Client()`, and `Close()`.
4. [ ] Register the connection in `core.RegisterBeforeStart` in `bootstrap/hook.go`.
5. [ ] Register the teardown in `core.App.Hooks().OnPostShutdown` in `bootstrap/hook.go`.

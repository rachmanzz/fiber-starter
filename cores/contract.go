package cores

import (
	"context"
	"sync"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type HookFunc func(ctx context.Context, app *AppContracts) error
type RouteFunc func(app *AppContracts) error
type AppContracts struct {
	App         *fiber.App
	beforeHooks []HookFunc
	once        sync.Once
}

func CreateContract() *AppContracts {
	return &AppContracts{}
}

func (app *AppContracts) Initialize() *AppContracts {
	NewLogger()
	return app
}

func (app *AppContracts) CreateApp(ctx context.Context) *AppContracts {
	app.once.Do(func() {
		if err := app.runBeforeHooks(ctx); err != nil {
			zap.L().Fatal("hook failed to run", zap.Error(err))
		}

		app.App = fiber.New()
	})

	return app
}

func (app *AppContracts) RegisterBefore(hook HookFunc) {
	app.beforeHooks = append(app.beforeHooks, hook)
}

func (app *AppContracts) runBeforeHooks(ctx context.Context) error {
	for _, hook := range app.beforeHooks {
		if err := hook(ctx, app); err != nil {
			return err
		}
	}
	return nil
}

func (app *AppContracts) RegisterRoute(route RouteFunc) {
	route(app)
}

func (app *AppContracts) Start() error {
	return app.App.Listen(Config().App.Port)
}

func (app *AppContracts) Shutdown() error {
	zap.L().Info("shutting down application...")

	return nil
}

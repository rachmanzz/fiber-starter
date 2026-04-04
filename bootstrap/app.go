package bootstrap

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rachmanzz/fiber-starter/app/routes"
	"github.com/rachmanzz/fiber-starter/cores"
	"go.uber.org/zap"
)

type Application struct {
	contract *cores.AppContracts
}

func NewApplication() *Application {
	core := cores.CreateContract().Initialize()

	core.RegisterBefore(func(ctx context.Context, app *cores.AppContracts) error {
		//cores.ConnectDB()
		return nil
	})

	// 4. Register After Hooks
	core.RegisterAfter(func(ctx context.Context, app *cores.AppContracts) error {
		//cores.CloseDB()
		return nil
	})

	return &Application{
		contract: core,
	}
}

func (app *Application) Bootstrap() *Application {
	ctx := context.Background()
	app.contract.CreateApp(ctx).RegisterRoute(func(c *cores.AppContracts) error {
		routes.ApiRoute(c.App)
		return nil
	})

	return app
}

func (app *Application) Run() {
	go func() {
		if err := app.contract.Start(); err != nil {
			zap.L().Fatal("Server failed to start", zap.Error(err))
		}
	}()
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	sig := <-stop
	zap.L().Info("Signal received, starting graceful shutdown", zap.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.contract.Shutdown(ctx); err != nil {
		zap.L().Error("Graceful shutdown failed", zap.Error(err))
	}

	zap.L().Info("Application stopped safely")
}

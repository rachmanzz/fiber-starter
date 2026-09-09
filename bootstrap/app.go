package bootstrap

import (
	"context"

	"github.com/rachmanzz/fiber-starter/app/routes"
	"github.com/rachmanzz/fiber-starter/cores"
	"go.uber.org/zap"
)

type Application struct {
	contract *cores.AppContracts
}

func NewApplication() *Application {
	core := cores.CreateContract().Initialize()
	InitializedHooks(core)
	RegisterDatabaseContract()

	if cores.Config().Database.Enable {
		cores.ConnectDB()
	}
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
	app.contract.SetupShutdownHook()

	if err := app.contract.Start(); err != nil {
		zap.L().Fatal("Server failed to start", zap.Error(err))
	}
}

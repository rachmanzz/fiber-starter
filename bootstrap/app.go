package bootstrap

import (
	"github.com/rachmanzz/fiber-starter/app/routes"
	"github.com/rachmanzz/fiber-starter/cores"
	"go.uber.org/zap"
)

type Application struct {
	contract *cores.AppContracts
}

func NewApplication() *Application {
	core := cores.CreateContract().Initialize()
	RegisterDatabaseContract()

	if cores.Config().Database.Enable {
		cores.ConnectDB()
	}
	return &Application{
		contract: core,
	}
}

func (app *Application) Bootstrap() *Application {
	core := app.contract.CreateApp()
	RegisterHook(core)
	RegisterMiddleware(core.App)
	core.RegisterRoute(func(c *cores.AppContracts) error {
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


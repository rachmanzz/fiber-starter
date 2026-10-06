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
	return &Application{
		contract: core,
	}
}

func (app *Application) Bootstrap() *Application {
	app.contract.
		CreateApp().
		RegisterHook(RegisterHook).
		RegisterMiddleware(RegisterMiddleware).
		RegisterRoute(
			routes.ApiRoute,
		)

	return app
}

func (app *Application) Run() {
	app.contract.SetupShutdownHook()

	if err := app.contract.Start(); err != nil {
		zap.L().Fatal("Server failed to start", zap.Error(err))
	}
}


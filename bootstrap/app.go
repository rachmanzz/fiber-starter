package bootstrap

import (
	"github.com/gofiber/fiber/v3"
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

// Contract returns the underlying *cores.AppContracts instance.
func (app *Application) Contract() *cores.AppContracts {
	return app.contract
}

// App returns the underlying *fiber.App instance.
func (app *Application) App() *fiber.App {
	if app.contract == nil {
		return nil
	}
	return app.contract.App
}

func (app *Application) Run() {
	if err := app.contract.Start(); err != nil {
		zap.L().Fatal("Server failed to start", zap.Error(err))
	}
}


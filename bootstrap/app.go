package bootstrap

import (
	"context"

	"github.com/rachmanzz/fiber-starter/app/routes"
	"github.com/rachmanzz/fiber-starter/cores"
)

type Application struct {
	contract *cores.AppContracts
}

func NewApplication() *Application {
	core := cores.CreateContract()
	return &Application{
		core,
	}
}

func (app *Application) Bootstrap() error {
	ctx := context.Background()
	app.contract.Initialize()
	app.contract.CreateApp(ctx).RegisterRoute(func(app *cores.AppContracts) error {
		routes.ApiRoute(app.App)
		return nil
	})
	return app.contract.Start()
}

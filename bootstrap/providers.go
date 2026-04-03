package bootstrap

import (
	"context"

	"github.com/rachmanzz/fiber-starter/cores"
)

func (app *Application) RegisterProviders() {
	app.contract.RegisterBefore(func(ctx context.Context, app *cores.AppContracts) error {
		return nil
	})
}

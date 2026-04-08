package bootstrap

import (
	"context"

	"github.com/rachmanzz/fiber-starter/cores"
)

func InitializedHooks(core *cores.AppContracts) {
	core.RegisterBefore(func(ctx context.Context, app *cores.AppContracts) error {
		return nil
	})

	core.RegisterAfter(func(ctx context.Context, app *cores.AppContracts) error {
		if cores.Config().Database.Enable {
			cores.CloseDB()
		}
		return nil
	})

}

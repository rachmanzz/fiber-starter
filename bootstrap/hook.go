package bootstrap

import (
	"github.com/rachmanzz/fiber-starter/cores"
	"go.uber.org/zap"
)

func RegisterHook(core *cores.AppContracts) {
	if core.App == nil {
		return
	}

	core.RegisterBeforeStart(func() error {
		return nil
	})

	core.App.Hooks().OnPostShutdown(func(err error) error {
		if cores.Config().Database.Enable {
			cores.CloseDB()
			zap.L().Info("Database connection pool closed successfully")
		}
		return nil
	})
}

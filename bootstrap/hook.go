package bootstrap

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	//"github.com/rachmanzz/fiber-starter/app/repository/contract"
	"github.com/rachmanzz/fiber-starter/cores"
	"go.uber.org/zap"
)

func RegisterHook(core *cores.AppContracts) {
	if core.App == nil {
		return
	}

	core.RegisterBeforeStart(func() error {
		// Database connection & contract registration
		if cores.Config().Database.Enable {
			cores.SetNamedDatabaseContract(func(name string, pool *pgxpool.Pool) {
				// contract.RegisterNamedDatabase(name, pool)
			})
			if err := cores.ConnectDB(); err != nil {
				return fmt.Errorf("failed to connect to database: %w", err)
			}
		}

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

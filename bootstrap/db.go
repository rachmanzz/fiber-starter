package bootstrap

import (
	"github.com/jackc/pgx/v5/pgxpool"
	//"github.com/rachmanzz/fiber-starter/app/repository/contract"
	"github.com/rachmanzz/fiber-starter/cores"
)

func RegisterDatabaseContract() {
	cores.SetDatabaseContract(func(pool *pgxpool.Pool) {
		//contract.DatabaseContract(pool)
	})
}

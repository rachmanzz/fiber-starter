package contract

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rachmanzz/fiber-starter/app/repository"
)

var (
	queries *repository.Queries
)

func DatabaseContract(db *pgxpool.Pool) {
	queries = repository.New(db)
}

func GetQueries() *repository.Queries {
	if queries == nil {
		panic("Repository Error: database contract must be initialized before accessing GetQueries")
	}
	return queries
}

package sqlc

import "github.com/jackc/pgx/v5/pgxpool"

type Store interface {
	Querier
}

type SQLStore struct {
	conn *pgxpool.Pool
	*Queries
}

func NewStore(conn *pgxpool.Pool) Store {
	return &SQLStore{
		conn:    conn,
		Queries: New(conn),
	}
}

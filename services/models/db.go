package models

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DAO struct {
	Conn *sqlx.DB
}

func InitDAO(db *sqlx.DB) *DAO {
	return &DAO{
		Conn: db,
	}
}

func Connect(connString string) *DAO {
	c, err := sqlx.Open("postgres", connString)
	if err != nil {
		panic(err)
	}

	return &DAO{
		Conn: c,
	}
}

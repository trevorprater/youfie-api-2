package postgres

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var (
	dbConnStr string
	defaultDB *sqlx.DB
)

func dbConnectionString() (string, error) {
	dbAddr := os.Getenv("YOUFIE_DB_ADDR")
	dbUser := os.Getenv("YOUFIE_DB_USER")
	dbPassword := os.Getenv("YOUFIE_DB_PASSWORD")

	host, port, err := net.SplitHostPort(dbAddr)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse DB host and port: %v", dbAddr)
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=youfie sslmode=disable", host, port, dbUser, dbPassword), nil
}

func connectDB() (*sqlx.DB, error) {
	dbConnStr, err := dbConnectionString()
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Connect("postgres", dbConnStr)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to DB: %v", dbConnStr)
	}

	//Like unused variables, columns which you ignore are a waste of network and database resources,
	//and detecting things like an incompatible mapping or a typo in a struct tag early
	//can be difficult without the mapper letting you know something wasn't found.

	//Despite this, there are some cases (DB MIGRATIONS) where ignoring columns with no destination might be desired.
	//For this, there is the Unsafe method on each Handle type which
	//returns a new copy of that handle whith this safety turned off:
	return db.Unsafe(), nil
}

func Load() error {
	var err error
	defaultDB, err = connectDB()
	if err != nil {
		return err
	}

	return nil
}

func init() {
	err := Load()
	if err != nil {
		log.Fatal(err)
	}
}

func DB() *sqlx.DB {
	return defaultDB
}

func main() {
	fmt.Println("vim-go")
}

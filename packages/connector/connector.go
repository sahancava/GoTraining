package connector

import (
	"database/sql"
	"fmt"
)

const (
	dbhost = "localhost"
	dbport = 5432
	dbuser = "postgres"
	dbpass = "Asdf1234"
	dbname = "dvdrental"
)

func ConnectDB() *sql.DB {
	var err error
	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", dbuser, dbpass, dbhost, dbport, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db
}
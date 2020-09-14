package connector

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
)

func ConnectDB() *sql.DB {
	var err error
	err = godotenv.Load("./packages/connector/.env")
	psqlInfo := os.ExpandEnv("$POSTGRES_URL")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db
}

func RunDB() {
	ConnectDB()
}

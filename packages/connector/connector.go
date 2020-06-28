package connector

import (
	"database/sql"
	"github.com/joho/godotenv"
	"os"
)

func ConnectDB() *sql.DB {
	var err error
	err = godotenv.Load(".env")
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

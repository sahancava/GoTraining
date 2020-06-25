package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)
var db22 *sql.DB
// repository contains the details of a repository
type CitySummary2 struct {
	CityID int `json:"city_id" db:"city_id"`
	CityName string `json:"city_name" db:"city"`
	InsertDate time.Time `json:"insert_date" db:"last_update"`
	CountryID int `json:"country_id" db:"country_id"`
}
type repositories2 struct {
	Repositories []CitySummary2
}

const (
	dbhost22 = "localhost"
	dbport22 = 5432
	dbuser22 = "postgres"
	dbpass22 = "Asdf1234"
	dbname22 = "dvdrental"
)

func main() {
	initDb()
	defer db22.Close()
	http.HandleFunc("/api/index", indexHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
func initDb() {
	//config := dbConfig()
	var err error
	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", dbuser22, dbpass22, dbhost22, dbport22, dbname22)

	db22, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db22.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	repos := repositories2{}

	err := queryRepos(&repos)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out, err := json.Marshal(repos)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(out))
}

// queryRepos first fetches the repositories data from the db
func queryRepos(repos *repositories2) error {
	rows, err := db22.Query(`
		SELECT
			city_id,
			city,
			country_id,
			last_update
		FROM city
		LIMIT 5`)
	if err != nil {
		return err
	}
	defer rows.Close()
	//snbs := make([]CitySummary2, 10)

	for rows.Next() {
		repo := CitySummary2{}
		err = rows.Scan(
			&repo.CityID,
			&repo.CityName,
			&repo.CountryID,
			&repo.InsertDate,
		)
		if err != nil {
			return err
		}
		//snbs = append(snbs,repo)
		repos.Repositories = append(repos.Repositories, repo)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}
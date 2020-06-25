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

var db *sql.DB

type City struct {
	CityID     int       `json:"city_id" db:"city_id"`
	CityName   string    `json:"city_name" db:"city"`
	InsertDate time.Time `json:"insert_date" db:"last_update"`
	CountryID  int       `json:"country_id" db:"country_id"`
}
type repository struct {
	cityRepository []City
}

const (
	dbhost = "localhost"
	dbport = 5432
	dbuser = "postgres"
	dbpass = "Asdf1234"
	dbname = "dvdrental"
)

func main() {
	initDb()
	defer db.Close()
	http.HandleFunc("/api/index", indexHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
func initDb() {

	var err error
	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", dbuser, dbpass, dbhost, dbport, dbname)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	repos := repository{}

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

func queryRepos(repos *repository) error {
	rows, err := db.Query(`
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

	for rows.Next() {
		repo := City{}
		err = rows.Scan(
			&repo.CityID,
			&repo.CityName,
			&repo.CountryID,
			&repo.InsertDate,
		)
		if err != nil {
			return err
		}
		repos.cityRepository = append(repos.cityRepository, repo)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

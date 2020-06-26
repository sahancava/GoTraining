package main

import (
	"encoding/json"
	"fmt"
	"github/GoTraining/packages/connector"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type Cities struct {
	CityID     int       `json:"city_id" db:"city_id"`
	CityName   string    `json:"city_name" db:"city"`
	InsertDate time.Time `json:"insert_date" db:"last_update"`
	CountryID  int       `json:"country_id" db:"country_id"`
	Country    string    `json:"country" db:"country"`
}
type repository struct {
	Data  []Cities
	limit int
}

var db = connector.ConnectDB()

func main() {
	connector.ConnectDB()
	defer db.Close()
	http.HandleFunc("/api/index", indexHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	repos := repository{}
	q := req.URL.Query().Get("limit")
	repos.limit, _ = strconv.Atoi(q)
	if q == "" {
		repos.limit = 10
	} else {
		if repos.limit < 1 {
			repos.limit = 10
		} else {
			repos.limit, _ = strconv.Atoi(q)
		}
	}
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
	sql := fmt.Sprintf(`
						SELECT
							city_id,
							city,
							ct.country_id,
							ct.last_update,
							co.country
						FROM city as ct
						JOIN country as co
						ON co.country_id=ct.country_id
						ORDER by city_id ASC
						LIMIT %d`,
		repos.limit)
	rows, err := db.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		repo := Cities{}
		err = rows.Scan(
			&repo.CityID,
			&repo.CityName,
			&repo.CountryID,
			&repo.InsertDate,
			&repo.Country,
		)
		if err != nil {
			return err
		}
		repos.Data = append(repos.Data, repo)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

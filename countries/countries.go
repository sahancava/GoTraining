package main

import (
	"encoding/json"
	"fmt"
	"github/GoTraining/packages/connector"
	"github/GoTraining/packages/token"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

type Countries struct {
	CountryID     int       `json:"country_id" db:"country_id"`
	Name   string    `json:"country" db:"country"`
}
type repository struct {
	Data  []Countries
	limit int
}

var db = connector.ConnectDB()

func main() {
	connector.RunDB()
	defer db.Close()
	http.HandleFunc("/GetCountries", indexHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	key := token.TokenChecker(w, req)
	for key != true {
		return
	}
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
							country_id,
							country
						FROM country
						ORDER by country_id ASC
						LIMIT %d`,
		repos.limit)
	rows, err := db.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		repo := Countries{}
		err = rows.Scan(
			&repo.CountryID,
			&repo.Name,
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

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github/GoTraining/packages/connector"
	"github/GoTraining/packages/token"
	"log"
	"net/http"
	"strconv"
	"time"
	"reflect"

	_ "github.com/lib/pq"
)

type City struct {
	CityID     int       `json:"city_id" db:"city_id"`
	CityName   string    `json:"city_name" db:"city"`
	InsertDate time.Time `json:"insert_date" db:"last_update"`
	CountryID  int       `json:"country_id" db:"country_id"`
	Country    string    `json:"country"`
}

type addCity struct {
	CityName   string    `json:"city_name" db:"city"`
	InsertDate time.Time `json:"insert_date" db:"last_update"`
	CountryID  int       `json:"country_id" db:"country_id"`
	Data       []addCity `json:"-"`
	Success    bool      `json:"success"`
	Message    string    `json:"Message"`
}

type repository struct {
	Data    []City
	ID      int    `json:"-"`
	Message string `json:"Message"`
	Success bool   `json:"is_Success"`
}

var db = connector.ConnectDB()

func main() {
	connector.RunDB()
	defer db.Close()
	handler := mux.NewRouter()
	handler.HandleFunc("/GetCity/{id}", indexHandler).Methods("GET")
	handler.HandleFunc("/AddCity/city={city}&country_id={country_id}", addHandler).Methods("GET")
	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func addHandler(w http.ResponseWriter, req *http.Request) {
	key := token.TokenChecker(w, req)
	if key != true {
		return
	}
	repos := addCity{}
	params := mux.Vars(req)
	repos.CityName, _ = params["city"]
	repos.CountryID, _ = strconv.Atoi(params["country_id"])
	repos.InsertDate = time.Now()

	err := addRepos(&repos)

	if err != nil {
		repos.Success = false
		http.Error(w, err.Error(), 500)
		return
	}
	repos.Success = true
	repos.Message = "Success"
	out, erro := json.Marshal(repos)

	if erro != nil {
		repos.Success = false
		out, _ = json.Marshal(repos)
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(out))
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	key := token.TokenChecker(w, req)
	if key != true {
		return
	}
	repos := repository{}
	params := mux.Vars(req)
	id:=reflect.TypeOf(params["id"]).Kind()
	repos.ID, _ = strconv.Atoi(params["id"])
	err := queryRepos(&repos)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	repos.Success = true
	repos.Message = "Success"
	out, erro := json.Marshal(repos)

	if id!=reflect.Int {
		fmt.Println("int değil")
		repos.Message = fmt.Sprintf("ID can only take integer values.")
		out, _ = json.Marshal(repos)
	}

	/*if len(repos.Data) < 1 {
		repos.Message = fmt.Sprintf("City with ID %d cannot be found.", repos.ID)
		repos.Success = false
		out, _ = json.Marshal(repos)
	}*/

	if erro != nil {
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
						WHERE city_id=%d`,
		repos.ID)
	rows, err := db.Query(sql)
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

func addRepos(repos *addCity) error {
	sql := fmt.Sprintf(`INSERT INTO city (city,country_id,last_update) values ('%v',%d,Now())`, repos.CityName, repos.CountryID)
	rows, err := db.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		repo := addCity{}
		err = rows.Scan(
			&repo.CityName,
			&repo.CountryID,
			&repo.InsertDate,
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

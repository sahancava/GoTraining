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

	_ "github.com/lib/pq"
)

type Actor struct {
	ActorID     int       `json:"actor_id" db:"actor_id"`
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName   string    `json:"last_name" db:"last_name"`
	InsertDate time.Time `json:"insert_date" db:"last_update"`
}

type addActor struct {
	ActorID     int       `json:"-"`
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName   string    `json:"last_name" db:"last_name"`
	InsertDate time.Time `json:"insert_date" db:"last_update"`
	Data       []addActor `json:"-"`
	Success    bool      `json:"success"`
	Message    string    `json:"Message"`
}

type repository struct {
	Data    []Actor
	ID      int    `json:"-"`
	Message string `json:"Message"`
	Success bool   `json:"is_Success"`
}

var db = connector.ConnectDB()

func main() {
	connector.RunDB()
	defer db.Close()
	handler := mux.NewRouter()
	handler.HandleFunc("/GetActor/{id}", indexHandler).Methods("GET")
	handler.HandleFunc("/AddActor/fname={first_name}&lname={last_name}", addHandler).Methods("GET")
	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func addHandler(w http.ResponseWriter, req *http.Request) {
	key := token.TokenChecker(w, req)
	if key != true {
		return
	}
	repos := addActor{}
	params := mux.Vars(req)
	repos.FirstName, _ = params["first_name"]
	repos.LastName, _ = params["first_name"]
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
	repos.ID, _ = strconv.Atoi(params["id"])
	err := queryRepos(&repos)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	repos.Success = true
	repos.Message = "Success"
	out, erro := json.Marshal(repos)

	if len(repos.Data) < 1 {
		repos.Message = fmt.Sprintf("Actor with ID %d cannot be found.", repos.ID)
		repos.Success = false
		out, _ = json.Marshal(repos)
	}

	if erro != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(out))
}

func queryRepos(repos *repository) error {
	sql := fmt.Sprintf(`
						SELECT
							actor_id,
							first_name,
							last_name,
							last_update
						FROM actor
						WHERE actor_id=%d`,
		repos.ID)
	rows, err := db.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		repo := Actor{}
		err = rows.Scan(
			&repo.ActorID,
			&repo.FirstName,
			&repo.LastName,
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

func addRepos(repos *addActor) error {
	sql := fmt.Sprintf(`INSERT INTO actor (first_name,last_name,last_update) values ('%v','%v',Now())`, repos.FirstName, repos.LastName)
	rows, err := db.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		repo := addActor{}
		err = rows.Scan(
			&repo.FirstName,
			&repo.LastName,
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

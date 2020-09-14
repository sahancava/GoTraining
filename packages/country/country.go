package country

import (
	"encoding/json"
	"fmt"
	"github/GoTraining/packages/connector"
	"github/GoTraining/packages/token"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

type Country struct {
	CountryID  int       `json:"country_id" db:"country_id"`
	Country    string    `json:"country" db:"country"`
	InsertDate time.Time `json:"insert_date" db:"last_update"`
}

type addCountry struct {
	CountryID  int          `json:"-"`
	Country    string       `json:"country" db:"country"`
	InsertDate time.Time    `json:"insert_date" db:"last_update"`
	Data       []addCountry `json:"-"`
	Success    bool         `json:"success"`
	Message    string       `json:"Message"`
}

type repository struct {
	Data    []Country
	ID      int    `json:"-"`
	Message string `json:"Message"`
	Success bool   `json:"is_Success"`
}

var db = connector.ConnectDB()

func AddHandler(w http.ResponseWriter, req *http.Request) {
	key := token.TokenChecker(w, req)
	if key != true {
		return
	}
	repos := addCountry{}
	params := mux.Vars(req)
	repos.Country, _ = params["country"]
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

func IndexHandler(w http.ResponseWriter, req *http.Request) {
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
		repos.Message = fmt.Sprintf("Country with ID %d cannot be found.", repos.ID)
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
							country_id,
							country,
							last_update
						FROM country
						WHERE country_id=%d`,
		repos.ID)
	rows, err := db.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		repo := Country{}
		err = rows.Scan(
			&repo.CountryID,
			&repo.Country,
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

func addRepos(repos *addCountry) error {
	sql := fmt.Sprintf(`INSERT INTO country (country,last_update) values ('%v',Now())`, repos.Country)
	rows, err := db.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		repo := addCountry{}
		err = rows.Scan(
			&repo.Country,
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

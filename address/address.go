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

type Address struct {
	AddressID   int       `json:"address_id" db:"address_id"`
	_Address    string    `json:"address" db:"address"`
	_Address2   string    `json:"address2" db:"address2"`
	District    string    `json:"district" db:"district"`
	CityID      int       `json:"city_id" db:"city_id"`
	City        string    `json:"city"`
	Postal_Code string    `json:"postal_code" db:"postal_code"`
	Phone       string    `json:"phone" db:"phone"`
	InsertDate  time.Time `json:"insert_date" db:"last_update"`
}

type addAddress struct {
	AddressID   int          `json:"address_id" db:"address_id"`
	_Address    string       `json:"address" db:"address"`
	_Address2   string       `json:"address2" db:"address2"`
	District    string       `json:"district" db:"district"`
	CityID      int          `json:"city_id" db:"city_id"`
	Postal_Code string       `json:"postal_code" db:"postal_code"`
	Phone       string       `json:"phone" db:"phone"`
	InsertDate  time.Time    `json:"insert_date" db:"last_update"`
	Data        []addAddress `json:"-"`
	Success     bool         `json:"success"`
	Message     string       `json:"Message"`
}

type repository struct {
	Data    []Address
	ID      int    `json:"-"`
	Message string `json:"Message"`
	Success bool   `json:"is_Success"`
}

var db = connector.ConnectDB()

func main() {
	connector.RunDB()
	defer db.Close()
	handler := mux.NewRouter()
	handler.HandleFunc("/GetAddress/{id}", indexHandler).Methods("GET")
	handler.HandleFunc("/AddAddress/"+
		"address={address}&"+
		"address2={address2}&"+
		"district={district}&"+
		"city_id={city_id}&"+
		"postal_code={postal_code}&"+
		"phone={phone}", addHandler).Methods("GET")
	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func addHandler(w http.ResponseWriter, req *http.Request) {
	key := token.TokenChecker(w, req)
	if key != true {
		return
	}
	repos := addAddress{}
	params := mux.Vars(req)
	repos._Address, _ = params["address"]
	repos._Address2, _ = params["address2"]
	repos.District, _ = params["district"]
	repos.CityID, _ = strconv.Atoi(params["city_id"])
	repos.Postal_Code, _ = params["postal_code"]
	repos.Phone, _ = params["phone"]
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
		repos.Message = fmt.Sprintf("Address with ID %d cannot be found.", repos.ID)
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
							ad.address_id,
							ad.address,
							ad.address2,
							ad.district,
							ad.city_id,
							ci.city,
							ad.postal_code,
							ad.phone,
							ad.last_update
						FROM address as ad
						join city as ci
						on ci.city_id=ad.city_id
						WHERE address_id=%d`,
		repos.ID)
	rows, err := db.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		repo := Address{}
		err = rows.Scan(
			&repo.AddressID,
			&repo._Address,
			&repo._Address2,
			&repo.District,
			&repo.CityID,
			&repo.City,
			&repo.Postal_Code,
			&repo.Phone,
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

func addRepos(repos *addAddress) error {
	sql := fmt.Sprintf(`INSERT INTO address (address,address2,district,city_id,postal_code,phone,last_update) values 
	('%v','%v','%v',%d,'%v','%v',Now()) --returning address_id`,
		repos._Address, repos._Address2, repos.District, repos.CityID, repos.Postal_Code, repos.Phone)

	/*sqlcheck := fmt.Sprintf(`select city_id from city where city_id=%d`, 606)

	denem, check_rows := db.Query(sqlcheck)

	if check_rows == nil {
		for denem.Next(){
			repo := addAddress{}
			check_rows = denem.Scan(
				&repo.AddressID,
			)
		}
	} else {*/
		rows, err := db.Query(sql)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			repo := addAddress{}
			err = rows.Scan(
				&repo._Address,
				&repo._Address2,
				&repo.District,
				&repo.CityID,
				&repo.Postal_Code,
				&repo.Phone,
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
	//}
	return nil
}

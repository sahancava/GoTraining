package actors

import (
	"encoding/json"
	"fmt"
	"github/GoTraining/packages/connector"
	"github/GoTraining/packages/token"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type Actors struct {
	ActorID    int       `json:"actor_id" db:"actor_id"`
	FirstName  string    `json:"first_name" db:"first_name"`
	LastName   string    `json:"last_name" db:"last_name"`
	InsertDate time.Time `json:"insert_date" db:"last_update"`
}

type actorRepository struct {
	Data  []Actors
	limit int
}

var db = connector.ConnectDB()

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	key := token.TokenChecker(w, req)
	for key != true {
		return
	}
	repos := actorRepository{}
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

func queryRepos(repos *actorRepository) error {
	sql := fmt.Sprintf(`
						SELECT
							actor_id,
							first_name,
							last_name,
							last_update
						FROM actor
						LIMIT %d`,
		repos.limit)
	rows, err := db.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		repo := Actors{}
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

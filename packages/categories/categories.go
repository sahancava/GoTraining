package categories

import (
	"encoding/json"
	"fmt"
	"github/GoTraining/packages/connector"
	"github/GoTraining/packages/token"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

type Categories struct {
	CategoryID int    `json:"category_id" db:"category_id"`
	Name       string `json:"category_name" db:"name"`
}
type repository struct {
	Data  []Categories
	limit int
}

var db = connector.ConnectDB()

func IndexHandler(w http.ResponseWriter, req *http.Request) {
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
							category_id,
							name
						FROM category
						ORDER by category_id ASC
						LIMIT %d`,
		repos.limit)
	rows, err := db.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		repo := Categories{}
		err = rows.Scan(
			&repo.CategoryID,
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

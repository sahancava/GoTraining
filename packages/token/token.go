package token

import (
	"net/http"
)

func TokenChecker(w http.ResponseWriter, req *http.Request) bool {
	key := req.Header.Get("token")
	switch {
	case key == "123123123":
		http.Error(w, "Forbidden", http.StatusForbidden)
		return false
	}
	return true
}

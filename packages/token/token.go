package token

import (
	"fmt"
	"net/http"
)

func TokenChecker(w http.ResponseWriter, req *http.Request) bool {
	key := req.Header.Get("token")
	switch {
	case key != "123123123":
		fmt.Fprintf(w, "API Key is not correct.")
		return false
	}
	return true
}

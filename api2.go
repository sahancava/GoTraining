package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
)

var db2 *gorm.DB
var err2 error

const (
	dbhost21 = "localhost"
	dbport21 = 5432
	dbuser21 = "postgres"
	dbpass21 = "Asdf1234"
	dbname21 = "dvdrental"
)

type city struct {
	cityname int `json:"city_name" db:"city_id"`
}
func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to HomePage!")
	fmt.Println("Endpoint Hit: HomePage")
}
func returnAllCities(w http.ResponseWriter, r *http.Request){
	city := city{}
	db2.Find(&city)
	fmt.Println("Endpoint Hit: returnAllCities")
	json.NewEncoder(w).Encode(city)
}
func handleRequests(){
	log.Println("Starting development server at http://127.0.0.1:10000/")
	log.Println("Quit the server with CONTROL-C.")
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/Cities", returnAllCities)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}
func main() {
	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", dbuser21, dbpass21, dbhost21, dbport21, dbname21)
	db2, err2 = gorm.Open("postgres",   psqlInfo)
	if err2 != nil {
		panic(err2.Error())
	}
	handleRequests()
}
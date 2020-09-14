package main

import (
	"github/GoTraining/packages/actor"
	"github/GoTraining/packages/actors"
	"github/GoTraining/packages/address"
	"github/GoTraining/packages/categories"
	"github/GoTraining/packages/cities"
	"github/GoTraining/packages/city"
	"github/GoTraining/packages/connector"
	"github/GoTraining/packages/countries"
	"github/GoTraining/packages/country"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var db = connector.ConnectDB()

func main() {
	connector.RunDB()
	defer db.Close()
	handler := mux.NewRouter()
	handler.HandleFunc("/GetActor/{id}", actor.IndexHandler).Methods("GET")
	handler.HandleFunc("/AddActor/fname={first_name}&lname={last_name}", actor.AddHandler).Methods("GET")
	handler.HandleFunc("/GetCity/{id}", city.IndexHandler).Methods("GET")
	handler.HandleFunc("/AddCity/city={city}&country_id={country_id}", city.AddHandler).Methods("GET")
	handler.HandleFunc("/GetActors", actors.IndexHandler).Methods("GET")
	handler.HandleFunc("/GetAddress/{id}", address.IndexHandler).Methods("GET")
	handler.HandleFunc("/AddAddress/"+
		"address={address}&"+
		"address2={address2}&"+
		"district={district}&"+
		"city_id={city_id}&"+
		"postal_code={postal_code}&"+
		"phone={phone}", address.AddHandler).Methods("GET")
	handler.HandleFunc("/GetCategories", categories.IndexHandler).Methods("GET")
	handler.HandleFunc("/GetCities", cities.IndexHandler).Methods("GET")
	handler.HandleFunc("/GetCountry/{id}", country.IndexHandler).Methods("GET")
	handler.HandleFunc("/AddCountry/country={country}", country.AddHandler).Methods("GET")
	handler.HandleFunc("/GetCountries", countries.IndexHandler)
	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(":80002", nil))
}

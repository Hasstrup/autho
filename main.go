package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/authenticate/middlewares"
	"github.com/authenticate/models"
	utils "github.com/authenticate/utilities"
	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/authenticate/controllers"
	"github.com/gorilla/mux"
)

var mongoConnectionString = flag.String("mongostring", "mongodb://localhost:27017", "The mongo connection string")

func main() {
	client := models.RegisterDatabase(*mongoConnectionString)
	start(registerRoutes(client))
}

func registerRoutes(c *mongo.Client) *mux.Router {
	r := mux.NewRouter()
	r.Use(middlewares.RequestLogger)
	s := r.PathPrefix("/api/v1").Subrouter()
	a := controllers.NewApplicationController(c)

	s.HandleFunc("/register", a.RegisterApplication).Methods("POST")
	s.HandleFunc("/application/update/{id}", a.UpdateApplicationDetails).Methods("PUT")
	s.HandleFunc("/application/{name}/{appKey}", a.GetApplicationDetails).Methods("Get")
	s.HandleFunc("/available/{name}", a.CheckAvailability).Methods("GET")
	s.HandleFunc("/applications", a.GetAllApplications).Methods("GET")
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.RespondWithJSON(w, 200, map[string]interface{}{"up": true})
	}).Methods("Get")

	return r
}

func start(router *mux.Router) {
	log.Printf("The server is listening on port 4600")
	err := http.ListenAndServe(":4600", router)
	if err != nil {
		log.Fatalf(err.Error())
	}

}

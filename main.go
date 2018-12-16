package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	utils "github.com/authenticate/utilities"

	"github.com/authenticate/controllers"
	"github.com/gorilla/mux"

	"github.com/mongodb/mongo-go-driver/mongo"
)

var mongoConnectionString = flag.String("mongostring", "Placeholder", "The mongo connection string")

func main() {
	registerDatabase()
	start(registerRoutes())
}

func registerDatabase() context.Context {
	flag.Parse()
	client, err := mongo.NewClient(*mongoConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to the mongo database: Error %v", err.Error())
		panic("failed to connect")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	return ctx
}

func registerRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.RespondWithJSON(w, 200, map[string]interface{}{"up": true})
	}).Methods("Get")
	s := r.PathPrefix("/api/v1").Subrouter()
	a := controllers.NewApplicationController()
	s.HandleFunc("/register", a.RegisterApplication).Methods("POST")
	s.HandleFunc("/application/update/{id}", a.UpdateApplicationDetails).Methods("PUT")
	s.HandleFunc("/application/{id}", a.GetApplicationDetails).Methods("Get")
	s.HandleFunc("/application/available", a.CheckAvailability).Methods("POST")
	return r
}

func start(router *mux.Router) {
	err := http.ListenAndServe(":4600", router)
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println("The server is listening on port 4600")

}

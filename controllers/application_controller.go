package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/authenticate/services"
	utils "github.com/authenticate/utilities"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type ApplicationController struct {
	service *services.ApplicationService
	client  *mongo.Client
}

func (ctr *ApplicationController) RegisterApplication(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()
	app, err := services.RegisterApplication(dec, ctr.client)
	if err != nil {
		utils.RespondWithJSON(w, 400, map[string]interface{}{"error": err.Error()})
	}
	utils.RespondWithJSON(w, 200, app)
}

func (ctr *ApplicationController) GetApplicationDetails(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	query := map[string]string{"name": name}
	// TODO: check that the application key matches the result before using the guy
	result := services.FindOneApplication(query, ctr.client)
	utils.RespondWithJSON(w, 200, result)
}

func (ctr *ApplicationController) GetAllApplications(w http.ResponseWriter, r *http.Request) {
	results := services.FindAllApplications(ctr.client)
	utils.RespondWithJSON(w, 200, results)
}

func (ctr *ApplicationController) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	result := services.FindOneApplication(map[string]string{"name": name}, ctr.client)
	if _, present := result["_id"]; present {
		utils.RespondWithJSON(w, 200, map[string]bool{"available": false})
		return
	}
	utils.RespondWithJSON(w, 200, map[string]bool{"available": true})
}

func (ctr *ApplicationController) UpdateApplicationDetails(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, 200, map[string]string{"here": "now"})
}

func NewApplicationController(client *mongo.Client) *ApplicationController {
	return &ApplicationController{service: &services.ApplicationService{}, client: client}
}

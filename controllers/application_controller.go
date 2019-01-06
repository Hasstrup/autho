package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/authenticate/models"
	"github.com/authenticate/services"
	utils "github.com/authenticate/utilities"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type ApplicationController struct {
	service *services.ApplicationService
	client  *mongo.Client
}

func (ctr *ApplicationController) RegisterApplication(w http.ResponseWriter, r *http.Request) {
	var m models.ApplicationModel
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		utils.RespondWithJSON(w, 400, map[string]interface{}{"error": err.Error()})
		return
	}
	app, err := services.RegisterApplication(&m, ctr.client)
	if err != nil {
		utils.RespondWithJSON(w, 400, map[string]interface{}{"error": err.Error()})
		return
	}
	utils.RespondWithJSON(w, 200, app)
}

func (ctr *ApplicationController) GetApplicationDetails(w http.ResponseWriter, r *http.Request) {
	// check that the request contains the secret key -
	// A more interesting feature will be to show some fields based on the presence/absence
	// of the pass(secret key)
	pass := r.Header.Get("x-access-token")
	if pass == "" {
		utils.RespondWithJSON(w, 403, map[string]string{"error": "Hey you need to send in you pass key"})
		return
	}
	name := mux.Vars(r)["name"]
	query := map[string]string{"name": name}
	result := services.FindOneApplication(query, ctr.client)
	if result == nil {
		utils.RespondWithJSON(w, 422, map[string]string{"error": "Oops looks like there is no application matching that record"})
		return
	}
	if services.CompareWithBcrypt(result["app_key"].(string), pass) {
		utils.RespondWithJSON(w, 200, map[string]interface{}{"result": utils.Transform(result)})
		return
	}
	utils.RespondWithJSON(w, 401, map[string]interface{}{"error": "You do not have the permission to view this application"})
}

func (ctr *ApplicationController) GetAllApplications(w http.ResponseWriter, r *http.Request) {
	if err := services.RootUserOnly(r.Header.Get("Authorization")); err != nil {
		utils.RespondWithJSON(w, 200, map[string]string{"error": err.Error()})
		return
	}
	results, err := services.FindAllApplications(ctr.client)
	if err != nil {
		utils.RespondWithJSON(w, 422, map[string]string{"error": err.Error()})
		return
	}
	utils.RespondWithJSON(w, 200, map[string]interface{}{"results": results})
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

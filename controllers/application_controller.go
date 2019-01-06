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

// check that the request contains the secret key -
// A more interesting feature will be to show some fields based on the presence/absence
// of the pass(secret key)
func (ctr *ApplicationController) GetApplicationDetails(w http.ResponseWriter, r *http.Request) {
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
	// TODO: Repeating this block in the method following this one, makes a case for abstraction into
	// a whole new function.
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.RespondWithJSON(w, 422, map[string]string{"error": "Sorry we could not parse the input sent"})
		return
	}
	application := body["main_application"].(map[string]interface{})
	pass := r.Header.Get("x-access-token")
	if !services.CompareWithBcrypt(application["app_key"].(string), pass) {
		utils.RespondWithJSON(w, 401, map[string]string{"error": "Hey you do not have permissions to do that"})
		return
	}

	utils.RespondWithJSON(w, 200, map[string]string{"here": "now"})
}

func (ctr *ApplicationController) RemoveApplication(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.RespondWithJSON(w, 422, map[string]string{"error": "Sorry we could not parse the input sent"})
		return
	}
	application := body["main_application"].(map[string]interface{})
	pass := r.Header.Get("x-access-token")
	if !services.CompareWithBcrypt(application["app_key"].(string), pass) {
		utils.RespondWithJSON(w, 401, map[string]string{"error": "Hey you do not have permissions to do that"})
		return
	}
	if err := services.RemoveApplication(application["name"].(string), ctr.client); err != nil {
		utils.RespondWithJSON(w, 400, map[string]string{"error": "Couldn't delete it at this point try again later"})
	}
	utils.RespondWithJSON(w, 209, map[string]interface{}{})
}

func NewApplicationController(client *mongo.Client) *ApplicationController {
	return &ApplicationController{service: &services.ApplicationService{}, client: client}
}

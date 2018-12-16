package controllers

import (
	"net/http"

	"github.com/authenticate/services"
	utils "github.com/authenticate/utilities"
)

type ApplicationController struct {
	service *services.ApplicationService
}

func (ctr *ApplicationController) RegisterApplication(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, 200, []interface{}{map[string]string{"status": "success"}})
}

func (ctr *ApplicationController) GetApplicationDetails(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, 200, map[string]string{"here": "now"})
}

func (ctr *ApplicationController) UpdateApplicationDetails(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, 200, map[string]string{"here": "now"})
}

func (ctr *ApplicationController) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, 200, map[string]string{"here": "now"})
}

func NewApplicationController() *ApplicationController {
	return &ApplicationController{service: &services.ApplicationService{}}
}

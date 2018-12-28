package controllers

import (
	"encoding/json"
	"net/http"

	utils "github.com/authenticate/utilities"
)

type AuthController struct {
	*ApplicationController
}

/*
	The logic is to get the the apiKey from the request and check if it exists in the db.
	then decrypt the key - and fetch the address - ping the address to make sure it's okay
	use the schema provided to validate the fields, write to the db
	then encrypt the payload using the email as the payload of the token and then we return
	We are assuming that the email is always going to be unique
*/
func (AuthController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		utils.RespondWithJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	utils.RespondWithJSON(w, 200, map[string]interface{}{"application": body})
}

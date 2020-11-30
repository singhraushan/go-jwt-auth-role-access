package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/singhraushan/go-jwt-auth-role-access/entities"
	"github.com/singhraushan/go-jwt-auth-role-access/service"
)

//AddNewEmployee will parse request body to get UserDetails & validate then invoke service layer to add entry in db
func AddNewEmployee(w http.ResponseWriter, r *http.Request) {
	log.Println("conroller--AddNewEmployee start.")
	userDetails, err := parseValidateUser(w, r)
	if err != nil {
		return
	}
	w.Write([]byte(service.AddNewEmployee(userDetails)))
	log.Println("conroller--AddNewEmployee end.")
}

//DeleteEmployee will parse request body to get UserDetails & validate then invoke service layer to delete entry
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	log.Println("conroller--DeleteEmployee start.")
	userDetails, err := parseValidateUser(w, r)
	if err != nil {
		return
	}
	w.Write([]byte(service.DeleteEmployee(userDetails)))
	log.Println("conroller--DeleteEmployee end.")
}

func parseValidateUser(w http.ResponseWriter, r *http.Request) (entities.UserDetails, error) {
	var userDetails entities.UserDetails
	err := json.NewDecoder(r.Body).Decode(&userDetails)
	if err != nil {
		log.Println("controller:Error while parsing userDetails", err)
		w.WriteHeader(http.StatusBadRequest)
		return userDetails, err
	}
	if validErrs := userDetails.Validate(); len(validErrs) > 0 {
		log.Println("controller:Validation error:", validErrs)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validErrs)
	}
	return userDetails, nil
}

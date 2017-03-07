package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/trevorprater/youfie-api-2/core/authentication"
	"github.com/trevorprater/youfie-api-2/core/postgres"
	"github.com/trevorprater/youfie-api-2/services"
	"github.com/trevorprater/youfie-api-2/services/models"
)

func GetUserByDisplayName(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	user, err := models.GetUserByDisplayName(vars["display_name"], postgres.DB())
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusNotFound)
	}
	userJson, err := json.MarshalIndent(&user, "", "    ")
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(userJson)
	}
}

func CreateUser(rw http.ResponseWriter, r *http.Request) {
	requestUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	resp, statusCode := requestUser.Insert(postgres.DB())
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	rw.Write(resp)
}

func UpdateUser(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	currentUser, err := authentication.GetUserByToken(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	updateUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&updateUser)

	resp, statusCode := currentUser.Update(postgres.DB(), updateUser)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	rw.Write(resp)
}

func DeleteUser(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	currentUser, err := authentication.GetUserByToken(r)
	rw.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))
	}
	resp, status := currentUser.Delete(postgres.DB())
	if status != http.StatusCreated {
		rw.WriteHeader(status)
		rw.Write(resp)
	} else {
		next(rw, r)
	}
}

func Login(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)
	requestUser.DisplayName = vars["display_name"]

	responseStatus, token := services.Login(requestUser)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(responseStatus)
	rw.Write(token)
}

func RefreshToken(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	requestUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(services.RefreshToken(requestUser.ID))
}

func Logout(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	err := services.Logout(r)
	rw.Header().Set("Content-Type", "application/json")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.WriteHeader(http.StatusOK)
	}
}

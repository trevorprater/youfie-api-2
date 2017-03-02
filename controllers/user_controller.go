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

func GetUserByDisplayName(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	user, err := models.GetUserByDisplayName(vars["display_name"], postgres.DB())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
	}
	userJson, err := json.MarshalIndent(&user, "", "    ")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(userJson)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	requestUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	resp, statusCode := requestUser.Insert(postgres.DB())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(resp)
}

func UpdateUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	currentUser, err := authentication.GetUserByToken(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	updateUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&updateUser)

	resp, statusCode := currentUser.Update(postgres.DB(), updateUser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(resp)
}

func DeleteUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	currentUser, err := authentication.GetUserByToken(r)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}
	resp, status := currentUser.Delete(postgres.DB())
	if status != http.StatusCreated {
		w.WriteHeader(status)
		w.Write(resp)
	} else {
		next(w, r)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)
	requestUser.DisplayName = vars["display_name"]

	responseStatus, token := services.Login(requestUser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(token)
}

func RefreshToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	requestUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	w.Header().Set("Content-Type", "application/json")
	w.Write(services.RefreshToken(requestUser.ID))
}

func Logout(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	err := services.Logout(r)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

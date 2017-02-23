package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/trevorprater/youfie-api-2/services/models"
	"github.com/trevorprater/youfie-api-2/core/postgres"
)


func CreateUser(w http.ResponseWriter, r *http.Request) {
	requestUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	responseStatus, user := requestUser.InsertUser(postgres.DB())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(user)
}


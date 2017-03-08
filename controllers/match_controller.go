package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/trevorprater/youfie-api-2/core/authentication"
	"github.com/trevorprater/youfie-api-2/core/postgres"
	"github.com/trevorprater/youfie-api-2/services/models"
)

func GetMatches(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rw.Header().Set("Content-Type", "application/json")
	user, err := authentication.GetUserByToken(r)
	if err != nil {
		log.Println("could not find user by token")
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("could not find user by token"))
		return
	}
	matches, err := models.GetMatchesForUser(user.ID, postgres.DB())
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("could not find user"))
		return
	}
	matchesJson, err := json.MarshalIndent(&matches, "", "    ")
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))
	} else {
		rw.WriteHeader(http.StatusOK)
		rw.Write(matchesJson)
	}
}

func GetMatch(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	match, err := models.GetMatchByID(vars["match_id"], postgres.DB())
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusNotFound)
	}
	matchJson, err := json.MarshalIndent(&match, "", "    ")
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(matchJson)
	}
}

func CreateMatch(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	requestMatch := new(models.Match)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestMatch)

	resp, statusCode := requestMatch.Insert(postgres.DB())
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	rw.Write(resp)
}

func UpdateMatch(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	rw.Header().Set("Content-Type", "application/json")
	updateMatch := new(models.Match)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&updateMatch)
	db := postgres.DB()
	dbMatch, err := models.GetMatchByID(vars["match_id"], db)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("cannot find match"))
	} else {
		resp, statusCode := dbMatch.Update(db, updateMatch)
		rw.WriteHeader(statusCode)
		rw.Write(resp)
	}
}

func DeleteMatch(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	match, err := models.GetMatchByID(vars["match_id"], postgres.DB())
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusNotFound)
	}
	resp, status := match.Delete(postgres.DB())
	rw.WriteHeader(status)
	rw.Write(resp)
}

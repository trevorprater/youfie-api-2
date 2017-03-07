package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/trevorprater/youfie-api-2/core/postgres"
	"github.com/trevorprater/youfie-api-2/services/models"
)

func GetMatches(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	matches, err := models.GetMatchesForUser(vars["display_name"], postgres.DB())
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusNotFound)
	}
	matchesJson, err := json.MarshalIndent(&matches, "", "    ")
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.Header().Set("Content-Type", "application/json")
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

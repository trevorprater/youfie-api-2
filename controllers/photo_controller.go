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

func GetPhotos(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rw.Header().Set("Content-Type", "application/json")
	user, err := authentication.GetUserByToken(r)
	if err != nil {
		log.Println("could not get user from token")
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("could not get user from token"))
		return
	}
	photos, err := models.GetPhotosForUser(user.ID, postgres.DB())
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("could not find user"))
	}
	if len(photos) == 0 {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("[]"))
		return
	}
	photosJson, err := json.MarshalIndent(&photos, "", "    ")
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))
	} else {
		rw.WriteHeader(http.StatusOK)
		rw.Write(photosJson)
	}
}

func CreatePhoto(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	requestPhoto := new(models.Photo)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestPhoto)

	rw.Header().Set("Content-Type", "application/json")
	user, err := authentication.GetUserByToken(r)
	if err != nil {
		log.Println("could not get user via token")
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("could not get user via supplied token"))
	} else {
		requestPhoto.OwnerID = user.ID
		resp, statusCode := requestPhoto.Insert(postgres.DB())
		rw.WriteHeader(statusCode)
		rw.Write(resp)
	}
}

func GetPhoto(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	photo, err := models.GetPhotoByID(vars["photo_id"], postgres.DB())
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	photoJson, err := json.MarshalIndent(&photo, "", "    ")
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(photoJson)
	}
}

func DeletePhoto(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	photo, err := models.GetPhotoByID(vars["photo_id"], postgres.DB())
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusNotFound)
	}
	resp, status := photo.Delete(postgres.DB())
	rw.WriteHeader(status)
	rw.Write(resp)
}

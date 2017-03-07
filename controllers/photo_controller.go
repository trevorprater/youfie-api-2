package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/trevorprater/youfie-api-2/core/postgres"
	"github.com/trevorprater/youfie-api-2/services/models"
)

func GetPhotos(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	photos, err := models.GetPhotosForUser(vars["display_name"], postgres.DB())
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusNotFound)
	}
	photosJson, err := json.MarshalIndent(&photos, "", "    ")
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(photosJson)
	}
}

func CreatePhoto(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	requestPhoto := new(models.Photo)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestPhoto)

	resp, statusCode := requestPhoto.Insert(postgres.DB())
	rw.Header.Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	rw.Write(resp)
}

func GetPhoto(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	photo, err := models.GetPhotoByID(vars["photo_id"], postgres.DB())
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusNotFound)
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

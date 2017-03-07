package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/trevorprater/youfie-api-2/core/postgres"
	"github.com/trevorprater/youfie-api-2/services/models"
)

func GetFaces(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	faces, err := models.GetFacesForPhoto(vars["photo_id"], postgres.DB())
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusNotFound)
	}
	facesJson, err := json.MarshalIndent(&faces, "", "    ")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(facesJson)
	}
}

func GetFace(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	face, err := models.GetFaceByID(vars["face_id"], postgres.DB())
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusNotFound)
	}
	faceJson, err := json.MarshalIndent(&face, "", "    ")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(faceJson)
	}
}

func CreateFace(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	requestFace := new(models.Face)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestPhoto)

	resp, statusCode := requestPhoto.Insert(postgres.DB())
	w.Header.Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(resp)
}

func DeleteFace(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	face, err := models.GetFaceByID(vars["face_id"], postgres.DB())
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusNotFound)
	}
	resp, status := face.Delete(postgres.DB())
	w.WriteHeader(status)
	w.Write(resp)
}

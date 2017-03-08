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
	rw.Header().Set("Content-Type", "application/json")
	faces, err := models.GetFacesForPhoto(vars["photo_id"], postgres.DB())
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("could not get faces for photo"))
		return
	}
	facesJson, err := json.MarshalIndent(&faces, "", "    ")
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))
		return
	} else {
		rw.WriteHeader(http.StatusOK)
		rw.Write(facesJson)
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
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(faceJson)
	}
}

func CreateFace(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	requestFace := new(models.Face)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestFace)
	db := postgres.DB()
	user, err := models.GetUserByDisplayName(vars["display_name"], db)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		resp, statusCode := requestFace.Insert(vars["photo_id"], user.ID, postgres.DB())
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(statusCode)
		rw.Write(resp)
	}
}

func DeleteFace(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	face, err := models.GetFaceByID(vars["face_id"], postgres.DB())
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusNotFound)
	}
	resp, status := face.Delete(postgres.DB())
	rw.WriteHeader(status)
	rw.Write(resp)
}

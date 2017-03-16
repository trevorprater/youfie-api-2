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

func CreateConversation(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	requestConversation := new(models.Conversation)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestConversation)

	rw.Header().Set("Content-Type", "application/json")
	user, err := authentication.GetUserByToken(r)
	if err != nil {
		log.Println("could not get user via token")
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("could not get user via supplied token")
	} else {
		resp, statusCode, := requestConversation.Insert(postgres.DB())
		rw.WriteHeader(statusCode)
		rw.Write(resp)
	}
}

func GetConversations(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rw.Header().Set("Content-Type", "application/json")
	user, err := authentication.GetUserByToken(r)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("could not get user from token"))
		return
	}
	conversations, err := models.GetConversationsForUser(user.ID, postgres.DB())
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("could not find user")
	}
	if len(conversations) == 0 {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("[]"))
		return
	}
	conversationsJson, err := json.MarshalIndent(&conversations, "", "    ")
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))
	} else {
		rw.WriteHeader(http.StatusOK)
		rw.Write(conversationsJson)
	}
}

func GetConversation(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	conversation, err := models.GetConversationByID(vars["conversation_id"], postgres.DB())
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	conversationJson, err := json.MarshalIndent(&conversation, "", "    ")
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("could not find conversation by id"))
	} else {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(conversationJson)
	}
}

func UpdateConversation(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	rw.Header().Set("Content-Type", "application/json")
	updateConversation := new(models.ConversationHttpResp)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&updateConversation)
	db := postgres.DB()
	dbConversation, err := models.GetConversationByID(vars["match_id"], db)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("cannot find conversation"))
	} else {
		resp, statusCode := dbConversation.Update(db, updateConversation.MessageText)
		rw.WriteHeader(statusCode)
		rw.Write(resp)
	}
}

func DeleteConversation(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	rw.Header().Set("Content-Type", "application/json")
	conversation, err := models.GetConversationByID(vars["conversation_id"}, postgres.DB())
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusNotFound)
	}
	resp, status := conversation.Delete(postgres.DB())
	rw.WriteHeader(status)
	rw.Write(resp)
}

package authentication

import (
	"fmt"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	request "github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/mux"
	"github.com/trevorprater/youfie-api-2/core/postgres"
	"github.com/trevorprater/youfie-api-2/services/models"
)

func GetUserByToken(req *http.Request) (*models.User, error) {
	token, err := parseTokenFromRequest(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if cleanToken, ok := token.(*jwt.Token); ok {
		userID, err := GetTokenSubject(cleanToken)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return models.GetUserByID(userID, postgres.DB())
	} else {
		return nil, err
	}
}

func parseTokenFromRequest(req *http.Request) (interface{}, error) {
	authBackend := InitJWTAuthenticationBackend()
	return request.ParseFromRequest(req, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		} else {
			return authBackend.PublicKey, nil
		}
	})
}

func RequireTokenAuthentication(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	authBackend := InitJWTAuthenticationBackend()
	token, err := parseTokenFromRequest(req)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusUnauthorized)
	} else {
		if cleanToken, ok := token.(*jwt.Token); ok {
			if cleanToken.Valid && !authBackend.IsInBlacklist(req.Header.Get("Authorization")) {
				next(rw, req)
			} else {
				rw.WriteHeader(http.StatusUnauthorized)
			}
		}
	}
}

func RequireUserReadPermission(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// TODO assert that user has 'follow' acccess to requested user. For now, simply authenticate if user.id == requested_user_id
	vars := mux.Vars(req)
	user, err := GetUserByToken(req)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusUnauthorized)
	}
	if user.DisplayName == vars["display_name"] {
		next(rw, req)
	} else {
		log.Println("the current user, " + user.Email + ", is not the user being requested.")
		rw.WriteHeader(http.StatusUnauthorized)
	}
}

func RequireUserWritePermission(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(req)
	user, err := GetUserByToken(req)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusUnauthorized)
	}
	if user.DisplayName == vars["display_name"] {
		next(rw, req)
	} else {
		log.Println("the current user, " + user.Email + ", is not the user being requested.")
		rw.WriteHeader(http.StatusUnauthorized)
	}
}

func RequireUserDeletePermission(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(req)
	user, err := GetUserByToken(req)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusUnauthorized)
	}
	if user.DisplayName == vars["display_name"] {
		next(rw, req)
	} else {
		log.Println("the current user, " + user.Email + ", is not the user being requested.")
		rw.WriteHeader(http.StatusUnauthorized)
	}
}

func RequireUserLogoutPermission(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(req)
	user, err := GetUserByToken(req)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusUnauthorized)
	}
	if user.DisplayName == vars["display_name"] {
		next(rw, req)
	} else {
		log.Println("the current user, " + user.Email + ", is not the user being requested.")
		rw.WriteHeader(http.StatusUnauthorized)
	}
}

func RequireUserConversationPermission(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	vars := mux.vars(req)
	user, err := GetUserByToken(req)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusUnauthorized)
	}
	conversation, err := models.GetConversationByID(vars["conversation_id"], postgres.DB())
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusNotFound)
	}
	if conversation.OwnerID == user.ID {
		return
	}

	matches, err := models.GetMatchesForUser(user.ID, postgres.DB())
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusUnauthorized)
	}

	for _, conversationParticipant := range conversation.Participants {
		for _, match := range matches {
			if match.FaceID == conversationParticipant.FaceID {
				return
			}
		}
	}

	rw.WriteHeader(http.StatusUnauthorized)
}

package services

import (
	"encoding/json"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	request "github.com/dgrijalva/jwt-go/request"
	"github.com/trevorprater/youfie-api-2/api/parameters"
	"github.com/trevorprater/youfie-api-2/core/authentication"
	"github.com/trevorprater/youfie-api-2/core/postgres"
	"github.com/trevorprater/youfie-api-2/services/models"
)

func Login(requestUser *models.User) (int, []byte) {
	authBackend := authentication.InitJWTAuthenticationBackend()

	authedUser, err := authBackend.Authenticate(requestUser)
	if err != nil || authedUser == nil {
		return http.StatusUnauthorized, []byte("login failed")
	}

	token, err := authBackend.GenerateToken(authedUser.ID)
	if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	} else {
		response, _ := json.Marshal(parameters.TokenAuthentication{token})
		err = authedUser.UpdateLastLogin(postgres.DB())
		if err != nil {
			log.Println(err)
		}
		return http.StatusOK, response
	}
}

func RefreshToken(userID string) []byte {
	authBackend := authentication.InitJWTAuthenticationBackend()
	token, err := authBackend.GenerateToken(userID)
	if err != nil {
		panic(err)
	}
	response, err := json.Marshal(parameters.TokenAuthentication{token})
	if err != nil {
		panic(err)
	}
	return response
}

func Logout(req *http.Request) error {
	authBackend := authentication.InitJWTAuthenticationBackend()
	tokenRequest, err := request.ParseFromRequest(req, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})
	if err != nil {
		return err
	}
	tokenString := req.Header.Get("Authorization")
	return authBackend.Logout(tokenString, tokenRequest)
}

package authentication

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/trevorprater/youfie-api-2/core/redis"
	"github.com/trevorprater/youfie-api-2/services/models"
	"github.com/trevorprater/youfie-api-2/settings"
	"golang.org/x/crypto/bcrypt"
)

type JWTAuthenticationBackend struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

const (
	tokenDuration = 72
	expireOffset  = 3600
)

var authBackendInstance *JWTAuthenticationBackend = nil

func InitJWTAuthenticationBackend() *JWTAuthenticationBackend {
	if authBackendInstance == nil {
		authBackendInstance = &JWTAuthenticationBackend{
			privateKey: getPrivateKey(),
			PublicKey:  getPublicKey(),
		}
	}

	return authBackendInstance
}

func (backend *JWTAuthenticationBackend) GenerateToken(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	token.Claims.(jwt.MapClaims)["exp"] = time.Now().Add(time.Hour * time.Duration(settings.Get().JWTExpirationDelta)).Unix()
	token.Claims.(jwt.MapClaims)["iat"] = time.Now().Unix()
	token.Claims.(jwt.MapClaims)["sub"] = userID
	tokenString, err := token.SignedString(backend.privateKey)
	if err != nil {
		panic(err)
		return "", err
	}
	return tokenString, nil
}

func (backend *JWTAuthenticationBackend) Authenticate(user *models.User, db sqlx.Ext) (*models.User, error, int) {
	dbUser, err := models.GetUserByDisplayName(user.DisplayName, db)
	if err != nil {
		log.Printf("Cannot get user by email: %v, error: %v", user.Email, err.Error())
		return nil, err, http.StatusUnauthorized
	}

	if dbUser.Disabled {
		log.Printf("user has been disabled until %t", dbUser.DisabledUntil)
		return nil, err, http.StatusUnauthorized
	}

	if bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(user.Password)) != nil {
		log.Println("Incorrect password supplied!")
		return nil, errors.New("incorrect password"), http.StatusUnprocessableEntity
	}
	if dbUser.DisplayName == user.DisplayName {
		return dbUser, nil, http.StatusOK
	}
	return nil, errors.New("login failed"), http.StatusInternalServerError
}

func (backend *JWTAuthenticationBackend) getTokenRemainingValidity(timestamp interface{}) int {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainder := tm.Sub(time.Now())
		if remainder > 0 {
			return int(remainder.Seconds() + expireOffset)
		}
	}
	return expireOffset
}

func (backend *JWTAuthenticationBackend) IsTokenValid(token *jwt.Token) bool {
	if claimsMap, ok := token.Claims.(jwt.MapClaims); ok {
		if backend.getTokenRemainingValidity(claimsMap["exp"]) > 0 {
			return true
		}
	}
	return false
}

func GetTokenSubject(token *jwt.Token) (string, error) {
	if claimsMap, ok := token.Claims.(jwt.MapClaims); ok {
		if subject, _ok := claimsMap["sub"].(string); _ok {
			return subject, nil
		}
	}
	return "", errors.New("Invalid Key: subject not found")
}

func (backend *JWTAuthenticationBackend) Logout(tokenString string, token *jwt.Token) error {
	redisConn := redis.Connect()
	return redisConn.SetValue(tokenString, tokenString, backend.getTokenRemainingValidity(token.Claims.(jwt.MapClaims)["exp"]))
}

func (backend *JWTAuthenticationBackend) IsInBlacklist(token string) bool {
	redisConn := redis.Connect()
	redisToken, _ := redisConn.GetValue(token)

	if redisToken == nil {
		return false
	}

	return true
}

func getPrivateKey() *rsa.PrivateKey {
	privateKeyFile, err := os.Open(os.Getenv("GOPATH") + settings.Get().PrivateKeyPath)
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	privateKeyFile.Close()

	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	return privateKeyImported
}

func getPublicKey() *rsa.PublicKey {
	publicKeyFile, err := os.Open(os.Getenv("GOPATH") + settings.Get().PublicKeyPath)
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := publicKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	publicKeyFile.Close()

	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	rsaPub, ok := publicKeyImported.(*rsa.PublicKey)

	if !ok {
		panic(err)
	}

	return rsaPub
}

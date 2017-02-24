package models

import (
	"log"
	"time"
	"net/http"
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/trevorprater/youfie-api-2/core/etc"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            string    `json:"id" form:"-" db:"id"`
	Email         string    `json:"email" form:"email" db:"email"`
	DisplayName   string    `json:"display_name" form:"display_name" db:"display_name"`
	PasswordHash  string    `json:"hash" form:"-" db:"hash"`
	Password      string    `json:"password,omit" form:"password" db:"password"`
	Admin         bool      `json:"-" form:"-" db:"admin"`
	Disabled      bool      `json:"-" form:"-" db:"disabled"`
	DisabledUntil time.Time `json:"disabled_until" form:"disabled_until" db:"disabled_until"`
	LastLogin     time.Time `json:"last_login" form:"-" db:"last_login"`
	CreatedAt     time.Time `json:"created_at" form:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" form:"updated_at" db:"updated_at"`
}

func GetUserByEmail(email string, db sqlx.Ext) (*User, error) {
	var user User
	err := sqlx.Get(db, &user, "SELECT * FROM users where email='"+email+"'")
	return &user, err
}

func GetUserByID(id string, db sqlx.Ext) (*User, error) {
	var user User
	err := sqlx.Get(db, &user, "SELECT * FROM users where id='"+id+"'")
	return &user, err
}

func (u *User) create(db sqlx.Ext) error {
	pwHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		log.Println("could not generate a hash from the password!")
		return err
	}
	u.PasswordHash = string(pwHash)

	_, err = sqlx.NamedExec(db, `
		INSERT INTO users
		(email, display_name, hash)
		VALUES (:email, :display_name, :hash)`, u)

	return err
}

func (u *User) InsertUser(db sqlx.Ext) (int, []byte) {
	err := u.create(db)

	// unique constraint violated
	if etc.Duperr(err) {
		log.Printf("User already exists %v: %v", u.Email, err.Error())
		return http.StatusConflict, []byte("")
	}

	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, []byte("")
	}

	response, err := json.MarshalIndent(&u, "", "    ")
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusCreated, response
}

func (u *User) UpdateLastLogin(db sqlx.Ext) error {
	_, err := sqlx.NamedExec(db, `
		UPDATE users SET last_login = current_timestamp WHERE email = :email`, u)
	if err != nil {
		log.Println("could not update last login for user")
	}
	return err
}

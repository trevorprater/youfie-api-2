package models

import (
	"errors"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pborman/uuid"
	"github.com/trevorprater/youfie-api-2/core/etc"
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
	if uuid.Parse(u.ID) == nil {
		return errors.New("Invalid User ID")
	}

	_, err := sqlx.NamedExec(db, `
		INSERT INTO users
		(id, email, display_name, pw_key)
		VALUES (:id, :email, :display_name, :pw_key)`, u)

	// unique constraint violated
	if etc.Duperr(err) {
		log.Printf("Failed to create user %v: %v", u.Email, err.Error())
	}

	return err
}

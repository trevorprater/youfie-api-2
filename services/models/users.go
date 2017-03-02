package models

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pborman/uuid"
	"github.com/trevorprater/youfie-api-2/core/etc"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            string    `json:"id" form:"-" db:"id"`
	Email         string    `json:"email" form:"email" db:"email"`
	DisplayName   string    `json:"display_name" form:"display_name" db:"display_name"`
	PasswordHash  string    `json:"-" form:"-" db:"hash"`
	Password      string    `json:"password,omitempty" form:"password,omitempty" db:"password"`
	Admin         bool      `json:"-" form:"-" db:"admin"`
	Disabled      bool      `json:"-" form:"-" db:"disabled"`
	DisabledUntil time.Time `json:"-" form:"-" db:"disabled_until"`
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

func GetUserByDisplayName(displayName string, db sqlx.Ext) (*User, error) {
	var user User
	err := sqlx.Get(db, &user, "SELECT * FROM users where display_name='"+displayName+"'")
	return &user, err
}

func (u *User) Insert(db sqlx.Ext) ([]byte, int) {
	if u.isPasswordValid() {
		pwHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
		if err != nil {
			log.Println("could not generate a hash from the password!")
			return []byte("invalid password provided"), http.StatusInternalServerError
		}
		u.PasswordHash = string(pwHash)
	} else {
		return []byte("invalid password provided"), http.StatusInternalServerError
	}

	if !u.isEmailValid() {
		return []byte("invalid email provided"), http.StatusInternalServerError
	}

	if !u.isDisplayNameValid() {
		return []byte("invalid display name provided"), http.StatusInternalServerError
	}

	_, err := sqlx.NamedExec(db, `
		INSERT INTO users
		(email, display_name, hash)
		VALUES (:email, :display_name, :hash)`, u)
	if err != nil {
		log.Println(err)
		if etc.Duperr(err) {
			return []byte("user already exists: " + err.Error()), http.StatusConflict
		} else {
			return []byte("internal server error"), http.StatusInternalServerError
		}
	} else {
		createdUser, err := GetUserByEmail(u.Email, db)
		if err != nil {
			log.Println(err)
			return []byte("internal server error"), http.StatusInternalServerError
		}
		userJson, err := json.MarshalIndent(&createdUser, "", "    ")
		if err != nil {
			log.Println(err)
			return []byte("internal server error"), http.StatusInternalServerError
		}
		return userJson, http.StatusCreated
	}
}

func (u *User) isEmailValid() bool {
	if u.Email != "" && etc.ValidEmailTest.MatchString(u.Email) {
		return true
	}
	return false
}

func (u *User) isPasswordValid() bool {
	if len(u.Password) >= 6 {
		return true
	}
	return false
}

func (u *User) isDisplayNameValid() bool {
	if len(u.DisplayName) >= 4 {
		return true
	}
	return false
}

func (u *User) Update(db sqlx.Ext, updatedUser *User) ([]byte, int) {
	updatedUser.CreatedAt = u.CreatedAt.In(time.UTC)
	updatedUser.UpdatedAt = time.Now().In(time.UTC)

	q := `
		UPDATE users
		SET updated_at = :updated_at,
			email = :email,
			display_name = :display_name,
			hash = :hash
		WHERE id = :id
		`
	updatedUser.ID = u.ID

	if updatedUser.isPasswordValid() {
		pwHash, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), 10)
		if err != nil {
			log.Println("could not generate a hash from the password!")
			return []byte("invalid password"), http.StatusUnprocessableEntity
		}
		updatedUser.PasswordHash = string(pwHash)
	} else if len(updatedUser.Password) > 0 {
		return []byte("invalid password"), http.StatusUnprocessableEntity
	} else {
		updatedUser.PasswordHash = u.PasswordHash
	}
	updatedUser.Password = ""

	if !updatedUser.isEmailValid() {
		if len(updatedUser.Email) > 0 {
			return []byte("invalid email"), http.StatusUnprocessableEntity
		}
		updatedUser.Email = u.Email
	}

	if !updatedUser.isDisplayNameValid() {
		if len(updatedUser.DisplayName) > 0 {
			return []byte("invalid display name"), http.StatusUnprocessableEntity
		}
		updatedUser.DisplayName = u.DisplayName
	}

	res, err := sqlx.NamedExec(db, q, updatedUser)
	if etc.Duperr(err) {
		log.Println(err)
		return []byte("unique constraint violated: " + err.Error()), http.StatusConflict
	}

	if err != nil {
		log.Println(err)
		return []byte("unable to update user!"), http.StatusInternalServerError
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return []byte("internal server error"), http.StatusInternalServerError
	}
	if count < 1 {
		log.Println("user not found")
		return []byte("user not found"), http.StatusNotFound
	}

	updatedUser, err = GetUserByEmail(updatedUser.Email, db)
	if err != nil {
		log.Println(err)
		return []byte("user not found: " + updatedUser.Email), http.StatusNotFound
	}

	userJson, err := json.MarshalIndent(&updatedUser, "", "    ")
	if err != nil {
		log.Println(err)
		return []byte("internal server error"), http.StatusInternalServerError
	}
	return userJson, http.StatusCreated
}

func (u *User) Delete(db sqlx.Ext) ([]byte, int) {
	if uuid.Parse(u.ID) == nil {
		log.Println("user not found: " + u.ID)
		return []byte("user not found"), http.StatusNotFound
	}
	res, err := db.Exec(`
		DELETE FROM users WHERE id = $1`, u.ID,
	)
	count, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return []byte("internal server error"), http.StatusInternalServerError
	}
	if count < 1 {
		log.Println(err)
		return []byte("user not found"), http.StatusNotFound
	}
	return []byte("user deleted"), http.StatusCreated
}

func (u *User) UpdateLastLogin(db sqlx.Ext) error {
	_, err := sqlx.NamedExec(db, `
		UPDATE users SET last_login = current_timestamp WHERE email = :email`, u)
	if err != nil {
		log.Println("could not update last login for user")
	}
	return err
}

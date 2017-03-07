package models

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pborman/uuid"
)

type Match struct {
	ID               string    `json:"id" form:"-" db:"id"`
	PhotoID          string    `json:"photo_id" form:"photo_id" db:"photo_id"`
	FaceID           string    `json:"face_id" form:"face_id" db:"face_id"`
	UserID           string    `json:"user_id" form:"user_id" db:"user_id"`
	Confidence       float64   `json:"confidence" form:"confidence" db:"confidence"`
	UserAcknowledged bool      `json:"user_ackd" form:"user_ackd" db:"user_ackd"`
	Confirmed        bool      `json:"confirmed" form:"confirmed" db:"confirmed"`
	CreatedAt        time.Time `json:"created_at" form:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" form:"updated_at" db:"updated_at"`
}

func GetMatchesForUser(userID string, db sqlx.Ext) ([]*Match, error) {
	// TODO: GET OFFSET AND LIMIT
	var matches []*Match
	err := sqlx.Get(db, &matches, "SELECT * FROM matches WHERE user_id='"+userID+"'")
	return &matches, err
}

func GetMatchByID(id string, db sqlx.Ext) (*Match, error) {
	var match Match
	err := sqlx.Get(db&match, "SELECT * FROM match WHERE id='"+id+"'")
	return &match, err
}

func (m *Match) Insert(db sqlx.Ext) ([]byte, int) {
	_, err := sqlx.NamedExec(db, `
		INSERT INTO matches
		(id, photo_id, face_id, user_id, confidence, confirmed)
		VALUES (:id, :photo_id, :face_id, :user_id, :confidence, :confirmed)`, m)
	if err != nil {
		log.Println(err)
		return []byte(err.Error()), http.StatusInternalServerError
	} else {
		createdMatch, err := GetMatchByID(m.ID, db)
		if err != nil {
			log.Println(err)
			return []byte("internal server error"), http.StatusInternalServerError
		}
		matchJson, err := json.MarshalIndent(&createdMatch, "", "    ")
		if err != nil {
			log.Println(err)
			return []byte("internal server error"), http.StatusInternalServerError
		}
		return photosJson, http.StatusCreated
	}
}

func (m *Match) Update(db sqlx.Ext, updatedMatch *Match) ([]byte, int) {
	m.UpdatedAt = time.Now().In(time.UTC)

	q := `
		UPDATE matches
		SET updated_at = :updated_at,
			confirmed = :confirmed,
			user_ackd = :user_ackd,
		WHERE id = :id`

	m.Confirmed = updatedMatch.Confirmed
	m.UserSeen = updatedMatch.UserSeen

	res, err := sqlx.NamedExec(db, q, m)
	if err != nil {
		log.Println(err)
		return []byte("unable to update match!"), http.StatusInternalServerError
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return []byte("internal server error"), http.StatusInternalServerError
	}
	if count < 1 {
		log.Println("match not found")
		return []byte("match not found"), http.StatusNotFound
	}

	updatedMatch, err = GetMatchByID(updatedMatch.ID, db)
	if err != nil {
		log.Println(err)
		return []byte("match not found: " + updatedMatch.ID), http.StatusNotFound
	}
	matchJson, err := json.MarshalIndent(&updatedMatch, "", "    ")
	if err != nil {
		log.Println(err)
		return []byte("internal server error"), http.StatusInternalServerError
	}
	return matchJson, http.StatusCreated
}

func (m *Match) Delete(db sqlx.Ext) ([]byte, int) {
	if uuid.Parse(m.ID) == nil {
		log.Println("match not found: " + m.ID)
		return []byte("match not found"), http.StatusNotFound
	}
	res, err := db.Exec(`
		DELETE FROM matches WHERE id = $1`, m.ID,
	)
	if err != nil {
		log.Println("could not delete match: " + err.Error())
		return []byte("internal server error"), http.StatusInternalServerError
	}
	count, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return []byte("internal server error"), http.StatusInternalServerError
	}
	if count < 1 {
		log.Println(err)
		return []byte("match not found"), http.StatusNotFound
	}
	return []byte("match deleted"), http.StatusCreated
}

package models

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pborman/uuid"
)

type Photo struct {
	ID         string  `json:"id" form:"-" db:"id"`
	OwnerID    string  `json:"owner_id" form:"owner_id" db:"owner_id"`
	Format     string  `json:"format" form:"format" db:"format"`
	Width      int     `json:"width" form:"-" db:"width"`
	Height     int     `json:"height" form:"-" db:"height"`
	StorageURL string  `json:"url" form:"-" db:"storage_url"`
	Latitude   float64 `json:"latitude" form:"latitude" db:"latitude"`
	Longitude  float64 `json:"longitude" form:"longitude" db:"longitude"`
	Processed  bool    `json:"processed" form:"-" db:"processed"`

	CreatedAt time.Time `json:"created_at" form:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at" db:"updated_at"`
}

func GetPhotosForUser(userID string, db sqlx.Ext) ([]*Photo, error) {
	// TODO: GET OFFSET AND LIMIT
	var photos []*Photo
	err := sqlx.Get(db, &photos, "SELECT * FROM photos WHERE user_id='"+userID+"'")
	return &photos, err
}

func GetPhotoByID(photoID string, db sqlx.Ext) (*Photo, error) {
	var photo Photo
	err := sqlx.Get(db, &photo, "SELECT * FROM photos WHERE id='"+id+"'")
	return &photo, err
}

func (p *Photo) Insert(userID string, db sqlx.Ext) ([]byte, int) {
	// TODO: process photo content, width, height, validate format, generate id, add lat/lng
	_, err := sqlx.NamedExec(db, `
		INSERT INTO photos
		(id, owner_id, format, content, width, height, storage_url, latitude, longitude)
		VALUES (:id, :owner_id, :format, :content, :width, :height, :storage_url, :latitude, :longitude)`, p)
	if err != nil {
		log.Println(err)
		return []byte(err.Error()), http.StatusInternalServerError
	} else {
		createdPhoto, err := GetPhotoByID(p.ID, db)
		if err != nil {
			log.Println(err)
			return []byte("internal server error"), http.StatusInternalServerError
		}
		photoJson, err := json.MarshalIndent(&createdPhoto, "", "    ")
		if err != nil {
			log.Println(err)
			return []byte("internal server error"), http.StatusInternalServerError
		}
		return photoJson, http.StatusCreated
	}
}

func (p *Photo) Delete(db sqlx.Ext) ([]byte, int) {
	if uuid.Parse(p.ID) == nil {
		log.Println("photo not found: " + p.ID)
		return []byte("photo not found"), http.StatusNotFound
	}
	res, err := db.Exec(`
		DELETE FROM photos WHERE id = $1`, u.ID,
	)
	count, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return []byte("internal server error"), http.StatusInternalServerError
	}
	if count < 1 {
		log.Println(err)
		return []byte("photo not found"), http.StatusNotFound
	}
	return []byte("user deleted"), http.StatusCreated
}

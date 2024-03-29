package models

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pborman/uuid"
)

type Face struct {
	ID            string `json:"id" form:"-" db:"id"`
	PhotoID       string `json:"photo_id" form:"photo_id" db:"photo_id"`
	FeatureVector string `json:"feature_vector" form:"feature_vector" db:"feature_vector"`

	TopLeftX     int `json:"bb_top_left_x" form:"bb_top_left_x" db:"bb_top_left_x"`
	TopLeftY     int `json:"bb_top_left_y" form:"bb_top_left_y" db:"bb_top_left_y"`
	TopRightX    int `json:"bb_top_right_x" form:"bb_top_right_x" db:"bb_top_right_x"`
	TopRightY    int `json:"bb_top_right_y" form:"bb_top_right_y" db:"bb_top_right_y"`
	BottomLeftX  int `json:"bb_bottom_left_x" form:"bb_bottom_left_x" db:"bb_bottom_left_x"`
	BottomLeftY  int `json:"bb_bottom_left_y" form:"bb_bottom_left_y" db:"bb_bottom_left_y"`
	BottomRightX int `json:"bb_bottom_right_x" form:"bb_bottom_right_x" db:"bb_bottom_right_x"`
	BottomRightY int `json:"bb_bottom_right_y" form:"bb_bottom_right_y" db:"bb_bottom_right_y"`

	CreatedAt time.Time `json:"created_at" form:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at" db:"updated_at"`
}

func GetFacesForPhoto(photoID string, db sqlx.Ext) ([]*Face, error) {
	var faces []*Face
	rows, err := db.Queryx("SELECT * FROM faces WHERE photo_id='" + photoID + "'")
	if err != nil {
		log.Println(err)
		return faces, err
	}
	for rows.Next() {
		var f Face
		err = rows.StructScan(&f)
		if err != nil {
			log.Println(err)
		}
		faces = append(faces, &f)
	}
	return faces, err
}

func GetFaceByID(faceID string, db sqlx.Ext) (*Face, error) {
	var face Face
	err := sqlx.Get(db, &face, "SELECT * FROM faces WHERE id='"+faceID+"'")
	return &face, err
}

func (f *Face) Insert(photoID, userID string, db sqlx.Ext) ([]byte, int) {
	f.ID = uuid.New()
	_, err := sqlx.NamedExec(db, `
	INSERT INTO faces
	(id, photo_id, feature_vector, bb_top_left_x, bb_top_left_y, bb_top_right_x, bb_top_right_y, bb_bottom_left_x, bb_bottom_left_y, bb_bottom_right_x, bb_bottom_right_y)
	VALUES(:id, :photo_id, :feature_vector, :bb_top_left_x, :bb_top_left_y, :bb_top_right_x, :bb_top_right_y, :bb_bottom_left_x, :bb_bottom_left_y, :bb_bottom_right_x, :bb_bottom_right_y)`, f)
	if err != nil {
		log.Println(err)
		return []byte("internal server error"), http.StatusInternalServerError
	} else {
		createdFace, err := GetFaceByID(f.ID, db)
		if err != nil {
			log.Println(err)
			return []byte("internal server error"), http.StatusInternalServerError
		}
		faceJson, err := json.MarshalIndent(&createdFace, "", "    ")
		if err != nil {
			log.Println(err)
			return []byte("internal server error"), http.StatusInternalServerError
		}
		return faceJson, http.StatusCreated
	}
}

func (f *Face) Delete(db sqlx.Ext) ([]byte, int) {
	if uuid.Parse(f.ID) == nil {
		log.Println("face not found: " + f.ID)
		return []byte("face not found"), http.StatusNotFound
	}
	res, err := db.Exec(`
		DELETE FROM faces WHERE id = $1`, f.ID,
	)
	count, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return []byte("internal server error"), http.StatusInternalServerError
	}
	if count < 1 {
		log.Println(err)
		return []byte("internal server error"), http.StatusInternalServerError
	}
	return []byte("face deleted"), http.StatusCreated
}

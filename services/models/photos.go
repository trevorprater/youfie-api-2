package models

type Photo struct {
	ID         string  `json:"id" form:"-" db:"id"`
	OwnerID    string  `json:"owner_id" form:"owner_id" db:"owner_id"`
	Format     string  `json:"format" form:"format" db:"format"`
	Content    string  `json: "content" form:"content" db:"content"`
	Width      int     `json:"width" form:"-" db:"width"`
	Height     int     `json:"height" form:"-" db:"height"`
	StorageURL string  `json:"url" form:"-" db:"storage_url"`
	Latitude   float64 `json:"latitude" form:"latitude" db:"latitude"`
	Longitude  float64 `json:"longitude" form:"longitude" db:"longitude"`
	Processed  bool    `json:"processed" form:"-" db:"processed"`
}

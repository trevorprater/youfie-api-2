package settings

import (
	"log"
	"os"
)

type Settings struct {
	PrivateKey         []byte
	PublicKey          []byte
	JWTExpirationDelta int
}

var settings Settings = Settings{}

func LoadSettings() {
	settings = Settings{}
	settings.PrivateKey = []byte(os.Getenv("YOUFIE_PRIVATE_KEY"))
	settings.PublicKey = []byte(os.Getenv("YOUFIE_PUBLIC_KEY"))

	if len(settings.PrivateKey) == 0 || len(settings.PublicKey) == 0 {
		log.Println("ENV VARS 'YOUFIE_PRIVATE_KEY' and/or 'YOUFIE_PUBLIC_KEY' are not set!")
	}
}

func Get() Settings {
	if &settings == nil {
		LoadSettings()
	}
	return settings
}

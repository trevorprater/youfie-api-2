package routers

import (
	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router = SetUserRoutes(router)
	router = SetPhotoRoutes(router)
	router = SetFaceRoutes(router)
	router = SetMatchRoutes(router)
	router = SetConversationRoutes(router)
	return router
}

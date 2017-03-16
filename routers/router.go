package routers

import (
	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()
	router = SetUserRoutes(router)
	router = SetPhotoRoutes(router)
	router = SetFaceRoutes(router)
	router = SetMatchRoutes(router)
	router = SetConversationRoutes(router)
	return router
}

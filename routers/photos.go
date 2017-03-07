package routers

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/trevorprater/youfie-api-2/core/authentication"
	"github.com/trevorprater/youfie-backend/controllers"
)

func SetPhotoRoutes(router *mux.Router) *mux.Router {
	router.Handle("/users/{display_name}/photos",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RequireUserReadPermission),
			negroni.HandlerFunc(controllers.GetPhotos),
		)).Methods("GET")
	router.Handle("/users/{display_name}/photos",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RequireUserWritePermission),
			negroni.HandlerFunc(controllers.CreatePhoto),
		)).Methods("POST")

	router.Handle("/users/{display_name}/photos/{photo_id}",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RequireUserReadPermission),
			negroni.HandlerFunc(controllers.GetPhoto),
		)).Methods("GET")
	router.Handle("/users/{display_name}/photos/{photo_id}",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RequireUserWritePermission),
			negroni.HandlerFunc(controllers.DeletePhoto),
		)).Methods("DELETE")

	return router
}

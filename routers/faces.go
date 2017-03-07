package routers

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	"github.com/trevorprater/youfie-api-2/controllers"
	"github.com/trevorprater/youfie-api-2/core/authentication"
)

func SetFaceRoutes(router *mux.Router) *mux.Router {
	router.Handle("/users/{display_name}/photos/{photo_id}/faces",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RequireUserReadPermission),
			negroni.HandlerFunc(controllers.GetFaces),
		)).Methods("GET")
	router.Handle("/users/{display_name}/photos/{photo_id}/faces/{face_id}",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RequireUserReadPermission),
			negroni.HandlerFunc(controllers.GetFace),
		)).Methods("GET")
	router.Handle("/users/{display_name}/photos/{photo_id}/faces",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RequireUserWritePermission),
			negroni.HandlerFunc(controllers.CreateFace),
		)).Methods("POST")
	router.Handle("/users/{display_name}/photos/{photo_id}/faces/{face_id}",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RequireUserWritePermission),
			negroni.HandlerFunc(controllers.DeleteFace),
		)).Methods("DELETE")
	return router
}

package routers

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func SetFaceRoutes(router *mux.Router) *mux.Router {
	router.Handle("/users/{display_name}/photos/{photo_id}/faces",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RequireUserReadPermission),
			negroni.HandlerFunc(authentication.GetFaces),
		)).Methods("GET")
	router.Handle("/users/{display_name}/photos/{photo_id}/faces/{face_id}",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RequireUserReadPermission),
			negroni.HandlerFunc(authentication.GetFace),
		)).Methods("GET")
	router.Handle("/users/{display_name}/photos/{photo_id}/faces",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RequireUserWritePermission),
			negroni.HandlerFunc(authentication.CreateFace),
		)).Methods("POST")
	router.Handle("/users/{display_name}/photos/{photo_id}/faces/{face_id}",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RequireUserWritePermission),
			negroni.HandlerFunc(authentication.DeleteFace),
		)).Methods("DELETE")
}

package routers

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/trevorprater/youfie-api-2/controllers"
	"github.com/trevorprater/youfie-api-2/core/authentication"
)

func SetUserRoutes(router *mux.Router) *mux.Router {
	// ====================CRUD USER=================================
	router.HandleFunc("/users", controllers.CreateUser).Methods("POST")

	router.Handle("/users/{display_name}",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RequireUserWritePermission),
			negroni.HandlerFunc(controllers.UpdateUser),
		)).Methods("PUT")
	router.Handle("/users/{display_name}",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RequireUserDeletePermission),
			negroni.HandlerFunc(controllers.DeleteUser),
			negroni.HandlerFunc(controllers.Logout),
		)).Methods("DELETE")
	router.Handle("/users/{display_name}",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RequireUserReadPermission),
			negroni.HandlerFunc(controllers.GetUserByDisplayName),
		)).Methods("GET")

	// ====================LOGIN/LOGOUT/REFRESH-TOKEN========================
	router.HandleFunc("/users/{display_name}/login", controllers.Login).Methods("POST")

	router.Handle("/users/{display_name}/refresh-token",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.RefreshToken),
		)).Methods("POST")

	router.Handle("/users/{display_name}/logout",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.Logout),
		)).Methods("POST")

	return router
}

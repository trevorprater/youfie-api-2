package routers

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/trevorprater/youfie-api-2/controllers"
	"github.com/trevorprater/youfie-api-2/core/authentication"
)

func SetAuthenticationRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/auth/login", controllers.Login).Methods("POST")
	router.Handle("/auth/refresh",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.RefreshToken),
		)).Methods("POST")
	router.Handle("/auth/logout",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.Logout),
		)).Methods("POST")
	return router
}

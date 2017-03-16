package routers

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/trevorprater/youfie-api-2/controllers"
	"github.com/trevorprater/youfie-api-2/core/authentication"
)

func SetConversationRoutes(router *mux.Router) *mux.Router {
	router.Handle("/users/{display_name}/conversations"
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RequireUserWritePermission),
			negroni.HandlerFunc(controllers.CreateConversation),
		)).Methods("POST")

	router.Handle("/users/{display_name}/conversations"
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RequireUserWritePermission),
			negroni.HandlerFunc(controllers.GetConversations),
		)).Methods("GET")

	router.Handle("/users/{display_name}/conversations/{conversation_id}",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RequireUserWritePermission),
			negroni.HandlerFunc(controllers.GetConversation),
		)).Methods("GET")

	router.Handle("/users/{display_name}/conversations/{conversation_id}",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RequireUserWritePermission),
			negroni.HandlerFunc(controllers.UpdateConversation),
		)).Methods("PUT")

	router.Handle("/users/{display_name}/conversations/{conversation_id}",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RequireUserWritePermission),
			negroni.HandlerFunc(controllers.DeleteConversation),
		)).Methods("DELETE")

	return router
}

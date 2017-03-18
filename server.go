package main

import (
	"log"
	"net/http"
	"runtime"

	"github.com/trevorprater/youfie-api-2/core/postgres"
	"github.com/trevorprater/youfie-api-2/routers"
	"github.com/trevorprater/youfie-api-2/settings"

	"github.com/codegangsta/negroni"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	err := postgres.Load()
	if err != nil {
		log.Fatal(err)
	}

	settings.LoadSettings()
	router := routers.InitRoutes()
	n := negroni.Classic()
	n.UseHandler(router)
	http.ListenAndServe(":5000", n)
}

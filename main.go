package main

import (
	"DockerHttpClient/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {
	fmt.Println("Start")

	router := mux.NewRouter()

	dockerHandler := handlers.DockerMiddleware()
	stripPrefixHandler := http.StripPrefix("/docker", dockerHandler)
	mainHandler := handlers.ProtectedHandler(stripPrefixHandler)

	router.PathPrefix("/docker").Handler(mainHandler)
	router.Path("/login").HandlerFunc(handlers.ShowLoginPage).Methods(http.MethodGet)
	router.Path("/login").HandlerFunc(handlers.DoLogin).Methods(http.MethodPost)

	n := negroni.Classic()

	n.UseHandler(router)

	server := http.Server{
		Addr:    ":8080",
		Handler: n,
	}

	log.Fatal(server.ListenAndServe())

	// listen on the /docker path

	// remove the /docker prefix

	// send the reminder of the path to the docker daemon

	// create templates based on the remaining URL

	// serve the tempolates
}

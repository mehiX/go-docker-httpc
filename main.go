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

	router.PathPrefix("/docker").Handler(http.StripPrefix("/docker", handlers.DockerMiddleware()))
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

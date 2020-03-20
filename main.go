package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mehix/go-docker-httpc/handlers"
	"github.com/urfave/negroni"
)

func main() {
	fmt.Println("Start")

	router := mux.NewRouter()

	handlerStripPrefixDocker := http.StripPrefix("/docker", negroni.New(
		negroni.WrapFunc(handlers.HandleDocker),
		negroni.WrapFunc(handlers.HandleServeTemplate),
	))

	dockerRouter := mux.NewRouter().PathPrefix("/docker").Subrouter().StrictSlash(true)
	dockerRouter.PathPrefix("/images").Handler(handlerStripPrefixDocker)

	router.PathPrefix("/docker").Handler(negroni.New(
		negroni.HandlerFunc(handlers.ProtectedHandler),
		negroni.Wrap(dockerRouter),
	))
	router.Path("/login").HandlerFunc(handlers.ShowLoginPage).Methods(http.MethodGet)
	router.Path("/login").HandlerFunc(handlers.DoLogin).Methods(http.MethodPost)
	router.Path("/home").Handler(negroni.New(
		negroni.HandlerFunc(handlers.ProtectedHandler),
	)).Methods(http.MethodGet)

	n := negroni.Classic()

	n.UseHandler(router)

	server := &http.Server{
		Handler: n,
		Addr:    ":8080",
	}

	log.Fatal(server.ListenAndServe())

	// listen on the /docker path

	// remove the /docker prefix

	// send the reminder of the path to the docker daemon

	// create templates based on the remaining URL

	// serve the tempolates
}

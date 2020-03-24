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
	dockerRouter.PathPrefix("/containers").Handler(handlerStripPrefixDocker)

	router.PathPrefix("/docker").Handler(negroni.New(
		negroni.HandlerFunc(handlers.ProtectedHandler),
		negroni.Wrap(dockerRouter),
	))
	router.Path("/login").HandlerFunc(handlers.ShowLoginPage).Methods(http.MethodGet)
	router.Path("/login").HandlerFunc(handlers.DoLogin).Methods(http.MethodPost)
	router.Path("/home").Handler(negroni.New(
		negroni.HandlerFunc(handlers.ProtectedHandler),
		negroni.WrapFunc(handlers.HandleServeTemplate),
	)).Methods(http.MethodGet)

	n := negroni.Classic()

	n.UseHandler(router)

	server := &http.Server{
		Handler: n,
		Addr:    ":8443",
	}

	log.Fatal(server.ListenAndServeTLS("./certs/certificate.pem", "./certs/key.pem"))

}

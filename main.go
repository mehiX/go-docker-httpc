package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/mehix/go-docker-httpc/handlers"
	"github.com/urfave/negroni"
)

var (
	addrSecure    = ":8443"
	addrNotSecure = ":8080"
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

	serverTLS := &http.Server{
		Handler: n,
		Addr:    addrSecure,
	}

	server := &http.Server{
		Handler: http.HandlerFunc(redirectNonSecure),
		Addr:    addrNotSecure,
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		log.Println("Listen TLS ", serverTLS.Addr)
		log.Fatal(serverTLS.ListenAndServeTLS("./certs/certificate.pem", "./certs/key.pem"))
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		log.Println("Listen non TLS ", server.Addr)
		log.Fatal(server.ListenAndServe())
		wg.Done()
	}()

	wg.Wait()

}

func redirectNonSecure(w http.ResponseWriter, r *http.Request) {
	serverName := strings.SplitN(r.Host, ":", 2)[0]
	redirectURL := "https://" + serverName + addrSecure + r.RequestURI

	http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
}

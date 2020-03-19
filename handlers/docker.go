package handlers

import (
	"log"
	"net/http"

	gc "github.com/gorilla/context"
	"github.com/mehix/go-docker-httpc/client"
)

// DockerMiddleware handles incoming requests and routes them to the docker daemon
func DockerMiddleware() http.Handler {
	return dockerHandler(http.HandlerFunc(handleServeTemplate))
}

func dockerHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleDocker(w, r)
		next.ServeHTTP(w, r)
	})
}

func handleDocker(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Path
	method := r.Method

	dockerResp, err := client.DockerHttp(q, method)

	if nil != err {
		log.Println(err)
		http.Error(w, "Error talking to Docker", http.StatusInternalServerError)
		return
	}

	gc.Set(r, "docker-response", dockerResp)
}

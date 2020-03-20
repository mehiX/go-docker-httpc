package handlers

import (
	"log"
	"net/http"

	gc "github.com/gorilla/context"
	"github.com/mehix/go-docker-httpc/client"
)

// HandleDocker makes a request to the docker engine and save the response in the context
func HandleDocker(w http.ResponseWriter, r *http.Request) {

	log.Println("In HandleDocker")

	q := r.URL.Path
	method := r.Method

	dockerResp, err := client.DockerHTTP(q, method)

	if nil != err {
		log.Println(err)
		http.Error(w, "Error talking to Docker", http.StatusInternalServerError)
		return
	}

	gc.Set(r, "docker-response", dockerResp)
}

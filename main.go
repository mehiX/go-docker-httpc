package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	gc "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type Image struct {
	Containers  int               `json:"Containers"`
	Created     int               `json:"Created"`
	ID          string            `json:"Id"`
	Labels      map[string]string `json:"Labels"`
	ParentID    string            `json:"ParentId"`
	RepoDigests []string          `json:"RepoDigests"`
	RepoTags    []string          `json:"RepoTags"`
	SharedSize  int               `json:"SharedSize"`
	Size        int               `json:"Size"`
	VirtualSize int               `json:"VirtualSize"`
}

func main() {
	fmt.Println("Start")

	router := mux.NewRouter()

	dockerHandler := dockerMiddleware(http.HandlerFunc(handleServeTemplate))

	router.PathPrefix("/docker").Handler(http.StripPrefix("/docker", dockerHandler))
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

func dockerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleDocker(w, r)
		//		gc.Set(r, "docker-response", gc.Get(r, "docker-response"))
		next.ServeHTTP(w, r)
	})
}

func handleDocker(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Path
	method := r.Method

	dockerResp, err := docker(q, method)

	if nil != err {
		log.Println(err)
		http.Error(w, "Error talking to Docker", http.StatusInternalServerError)
		return
	}

	gc.Set(r, "docker-response", dockerResp)
}

func handleServeTemplate(w http.ResponseWriter, r *http.Request) {
	template, err := template.New(r.URL.Path).ParseFiles("tmpl/images.html", "tmpl/containers.html")

	if nil != err {
		log.Fatal(err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	dockerResponse := gc.Get(r, "docker-response").(string)

	log.Println("Docker response", dockerResponse)

	if strings.Index(r.URL.Path, "/images/") == 0 {
		var images []Image
		json.NewDecoder(strings.NewReader(dockerResponse)).Decode(&images)
		template.Execute(w, images)
	} else {
		http.Error(w, "Operation not supported", http.StatusBadRequest)
	}

}

func docker(q string, method string) (string, error) {
	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.Dial("unix", "/var/run/docker.sock")
			},
		},
	}

	req, err := http.NewRequest(method, "http://unix"+q, nil)
	if nil != err {
		return "", err
	}

	resp, err := client.Do(req)
	if nil != err {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return "", err
	}

	return string(body), nil
}

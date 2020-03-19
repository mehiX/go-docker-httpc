package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"

	gc "github.com/gorilla/context"
	"github.com/mehix/go-docker-httpc/data"
)

var templates map[string]*template.Template

func init() {
	templates = make(map[string]*template.Template, 0)
	templates["/images"] = template.Must(template.ParseFiles("./tmpl/images.html", "./tmpl/base.html"))
	templates["/containers"] = template.Must(template.ParseFiles("./tmpl/containers.html", "./tmpl/base.html"))
	templates["/login"] = template.Must(template.ParseFiles("tmpl/login.html", "tmpl/base.html"))

}

func handleServeTemplate(w http.ResponseWriter, r *http.Request) {
	dockerResponse := gc.Get(r, "docker-response").(string)

	dReader := strings.NewReader(dockerResponse)

	if strings.Index(r.URL.Path, "/images/") == 0 {
		var images []data.Image
		json.NewDecoder(dReader).Decode(&images)

		log.Printf("%#v\n", images)

		templates["/images"].ExecuteTemplate(w, "base", images)
	} else if 0 == strings.Index(r.URL.Path, "/containers/") {
		var containers []data.Container
		json.NewDecoder(dReader).Decode(&containers)
		templates["/containers"].ExecuteTemplate(w, "base", containers)
	} else {
		http.Error(w, "Operation not supported", http.StatusBadRequest)
	}

}

package handlers

import (
	"DockerHttpClient/data"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"

	gc "github.com/gorilla/context"
)

var templates map[string]*template.Template

func init() {
	templates = make(map[string]*template.Template, 0)
	templates["/images"] = template.Must(template.ParseFiles("./tmpl/images.html", "./tmpl/base.html"))
	templates["/containers"] = template.Must(template.ParseFiles("./tmpl/containers.html", "./tmpl/base.html"))

}

func handleServeTemplate(w http.ResponseWriter, r *http.Request) {
	dockerResponse := gc.Get(r, "docker-response").(string)

	if strings.Index(r.URL.Path, "/images/") == 0 {
		var images []data.Image
		json.NewDecoder(strings.NewReader(dockerResponse)).Decode(&images)

		log.Printf("%#v\n", images)

		templates["/images"].ExecuteTemplate(w, "base", images)
	} else {
		http.Error(w, "Operation not supported", http.StatusBadRequest)
	}

}

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
	templates["/images"] = template.Must(template.ParseFiles("./tmpl/images.html", "./tmpl/base.html", "./tmpl/menu.html"))
	templates["/containers"] = template.Must(template.ParseFiles("./tmpl/containers.html", "./tmpl/base.html", "./tmpl/menu.html"))
	templates["/login"] = template.Must(template.ParseFiles("tmpl/login.html", "tmpl/base.html"))
	templates["/home"] = template.Must(template.ParseFiles("tmpl/home.html", "tmpl/base.html"))

}

func getDockerResponse(r *http.Request) string {
	if dr := gc.Get(r, "docker-response"); nil != dr {
		return dr.(string)
	}

	return ""
}

// HandleServeTemplate fetch the docker engine response from the context
// and show the appropriate template based on the base URL
func HandleServeTemplate(w http.ResponseWriter, r *http.Request) {
	dockerResponse := getDockerResponse(r)

	if "" == dockerResponse {
		log.Println("Empty response from docker in context")
	}

	dReader := strings.NewReader(dockerResponse)

	path := r.URL.Path
	key := r.URL.Query().Get("key")
	td := data.TemplateData{Key: key}

	if strings.Index(path, "/images/") == 0 {
		var images []data.Image
		json.NewDecoder(dReader).Decode(&images)

		td.Data = images
		templates["/images"].ExecuteTemplate(w, "base", td)
	} else if 0 == strings.Index(path, "/containers/") {
		var containers []data.Container
		json.NewDecoder(dReader).Decode(&containers)
		td.Data = containers
		templates["/containers"].ExecuteTemplate(w, "base", td)
	} else if 0 == strings.Index(path, "/home") {
		templates["/home"].ExecuteTemplate(w, "base", td)
	} else {
		http.Error(w, "Operation not supported for: "+path, http.StatusBadRequest)
	}

}

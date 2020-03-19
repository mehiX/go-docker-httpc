package handlers

import (
	"DockerHttpClient/handlers/service"
	"fmt"
	"log"
	"net/http"
)

// ProtectedHandler makes sure the user is authorized for this resource
func ProtectedHandler(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Checking protected")
		if "" != r.URL.Query().Get("user") {
			next.ServeHTTP(w, r)
		} else {
			nextURL := fmt.Sprintf("/login?next=%s", r.URL.Path)
			http.Redirect(w, r, nextURL, http.StatusFound)
		}
	})
}

// ShowLoginPage display the login form
func ShowLoginPage(w http.ResponseWriter, r *http.Request) {
	templates["/login"].ExecuteTemplate(w, "base", nil)
}

// DoLogin receives a POST request and attempts to validate the credentials sent. Receives the next URL as a query parameter
// TODO: Validate the next query parameter against parameter injection
func DoLogin(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("user")
	passwd := r.FormValue("passwd")
	next := r.URL.Query().Get("next")

	user, err := service.Login(username, passwd)

	if nil != err {
		http.Redirect(w, r, "/login?next="+next, http.StatusFound)
	} else {
		http.Redirect(w, r, next+"?user="+user.Username, http.StatusFound)
	}
}

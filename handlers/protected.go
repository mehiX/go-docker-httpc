package handlers

import (
	"DockerHttpClient/handlers/service"
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
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	})
}

func ShowLoginPage(w http.ResponseWriter, r *http.Request) {
	templates["/login"].ExecuteTemplate(w, "base", nil)
}

func DoLogin(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("user")
	passwd := r.FormValue("passwd")
	next := "/docker/images/json"

	user, err := service.Login(username, passwd)

	if nil != err {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		http.Redirect(w, r, next+"?user="+user.Username, http.StatusFound)
	}
}

package handlers

import (
	"log"
	"net/http"
)

// ProtectedHandler makes sure the user is authorized for this resource
func ProtectedHandler(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Checking protected")
		if "" == r.URL.Query().Get("user") {
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	})
}

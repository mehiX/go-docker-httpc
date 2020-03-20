package handlers

import "net/http"

// HomeHandler landing page after successful login
func HomeHandler() http.Handler {
	return http.HandlerFunc(handleHomeTemplate)
}

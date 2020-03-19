package handlers

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/mehix/go-docker-httpc/handlers/service"
)

var (
	signingKey   *rsa.PrivateKey
	verifyingKey *rsa.PublicKey
)

func init() {

	var err error

	signingKey, err = getPk("keys/app.rsa")
	if nil != err {
		log.Panicf("Private key error; %v", err)
	}

	verifyingKey, err = getPubKey("keys/app.rsa.pub")
	if nil != err {
		log.Panicf("Public key error: %v", err)
	}
}

func getPk(path string) (*rsa.PrivateKey, error) {
	key, err := ioutil.ReadFile(path)

	if nil != err {
		return nil, err
	}

	return jwt.ParseRSAPrivateKeyFromPEM(key)
}

func getPubKey(path string) (*rsa.PublicKey, error) {
	key, err := ioutil.ReadFile(path)

	if nil != err {
		return nil, err
	}

	return jwt.ParseRSAPublicKeyFromPEM(key)
}

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

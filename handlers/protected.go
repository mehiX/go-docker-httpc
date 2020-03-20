package handlers

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/mehix/go-docker-httpc/data"
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
func ProtectedHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	key := r.URL.Query().Get("key")
	_, err := (&service.Token{}).FindByKey(key)

	if nil != err {
		log.Println("NOT Found token!!")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	log.Println("Found token!!")
	context.Set(r, "tokenKey", key)
	next.ServeHTTP(w, r)

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

	user, err := service.Login(username, passwd)

	if nil != err {
		// login error
		w.Header().Set("Content-Type", "text/html")
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		claims := data.UserClaims{
			Username: user.Username,
			Role:     user.Role,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: 15000,
				Issuer:    "go-docker-httpc",
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		ss, err := token.SignedString(signingKey)

		if nil != err {
			log.Printf("Signing error: %v\n", err)
		} else {
			log.Printf("Signed token: %s\n", ss)
		}

		tokenObj := service.NewStoredToken(ss).Store()

		templates["/home"].ExecuteTemplate(w, "base", tokenObj)
	}
}

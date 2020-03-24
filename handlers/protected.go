package handlers

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"github.com/mehix/go-docker-httpc/data"
	"github.com/mehix/go-docker-httpc/handlers/service"
)

var (
	signingKey   *rsa.PrivateKey
	verifyingKey *rsa.PublicKey

	// TODO change this: https://github.com/gorilla/sessions
	cookieStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	sessionName = "jwt-auth"
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

	cookieStore.Options.Secure = true
	cookieStore.Options.HttpOnly = true
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

	//key := r.URL.Query().Get("key")
	session, _ := cookieStore.Get(r, sessionName)

	var key string
	if k, ok := session.Values["key"]; ok {
		key = k.(string)
	} else {
		key = ""
	}

	_, err := (&service.Token{}).FindByKey(key)

	if nil != err {
		log.Println("NOT Found token!! Key is: ", key)
		log.Println(err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

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

		session, _ := cookieStore.Get(r, sessionName)
		session.Values["key"] = tokenObj.Key
		if err := session.Save(r, w); nil != err {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		templates["/home"].ExecuteTemplate(w, "base", nil)
	}
}

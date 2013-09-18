package hello

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"html/template"
	// "io/ioutil"
	"errors"
	"log"
	"net/http"
)

type User struct {
	Username string
	Password string
}

// templates variable
var templates = template.Must(template.ParseGlob("app/*.html"))

// session
var store = sessions.NewCookieStore([]byte("something-very-secret"))

func init() {
	http.HandleFunc("/api/products", productHandler)
	http.HandleFunc("/api/login", loginHandler)
	http.HandleFunc("/api/logout", logoutHandler)
	http.HandleFunc("/", handler)
}

func (user User) login() (string, error) {
	if user.Username == "admin" && user.Password == "admin" {
		return "really-long-string", nil
	}
	return "", errors.New("Fail login")
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-inventory")
	if session.Values["token"] != "really-long-string" {
		http.Error(w, "Unauthorized", http.StatusForbidden)
	}
	fmt.Fprint(w, "ok")
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-inventory")
	session.Values["token"] = ""
	session.Save(r, w)
	fmt.Fprint(w, "ok")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	var user User
	err := decoder.Decode(&user)
	if err != nil {
		log.Println(err.Error())
	}

	result, err := user.login()
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
	}

	log.Println(result)
	session, _ := store.Get(r, "session-inventory")
	session.Values["token"] = result
	session.Save(r, w)
	// if result == "" {
	// 	log.Println("FAIL LOGGING IN")
	// }
	// else {
	//     log.Println("OK")
	// }

	fmt.Fprint(w, "Hello, world!")
}

func handler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", nil)
	handleError(w, err)
}

func handleError(res http.ResponseWriter, err error) {
	if err != nil {
		log.Println("=====ERROR=====")
		log.Println(err.Error())
		log.Println("=====END ERROR=====")
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

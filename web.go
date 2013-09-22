package simple_inventory

import (
	"appengine"
	auth_web "auth"
	"html/template"
	"log"
	"net/http"
)

type User struct {
	Username string
	Password string
}

// templates variable
var templates = template.Must(template.ParseGlob("app/*.html"))

func init() {
	http.HandleFunc("/api/products", productHandler)
	// http.HandleFunc("/api/login", loginHandler)
	// http.HandleFunc("/api/logout", logoutHandler)
	http.HandleFunc("/", handler)
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	err := auth_web.ValidateSession(c, r)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unauthorized", http.StatusForbidden)
	}
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

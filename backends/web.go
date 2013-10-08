package simple_inventory

import (
	"backends/auth"
	"backends/products"
	mux "github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

// templates variable
var templates = template.Must(template.ParseGlob("app/*.html"))

func init() {
	r := mux.NewRouter()

	auth.InitRouter(r)
	products.InitRouter(r)

	r.HandleFunc("/", handler)

	http.Handle("/", r)
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

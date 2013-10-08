package simple_inventory

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"

	"auth"
	"fmt"
	mux "github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

type Product struct {
	Nama            string
	Kubikasi        string
	Harga_modal     string
	Harga_jual      string
	Ongkos_expedisi string
	Stock           string
	Ongkos_kirim    string
	ID              int64 `datastore:"-"`
}

// templates variable
var templates = template.Must(template.ParseGlob("app/*.html"))

func init() {
	r := mux.NewRouter()

	auth.InitRouter(r)

	r.HandleFunc("/", handler)
	r.HandleFunc("/api/products", productListHandler).Methods("GET")
	r.HandleFunc("/api/products", productCreateHandler).Methods("POST")

	http.Handle("/", r)
}

func productListHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	query := datastore.NewQuery("Products").Limit(20)
	products := make([]Product, 0, 20)

	//query all products
	keys, err := query.GetAll(c, &products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//attach the key queried to the struct
	for index := range products {
		products[index].ID = keys[index].IntID()
	}

	response, err := json.Marshal(products)
	handleError(w, err)
	fmt.Fprintf(w, "%s", response)
}

func productCreateHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	err := auth.ValidateSession(c, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// decode request
	decoder := json.NewDecoder(r.Body)
	var product Product

	err = decoder.Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// save the user
	key := datastore.NewIncompleteKey(c, "Products", nil)
	_, err = datastore.Put(c, key, &product)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

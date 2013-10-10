package products

import (
	"appengine"
	"appengine/datastore"
	"backends/auth"
	"encoding/json"
	"fmt"
	mux "github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Product struct {
	Nama            string
	Kubikasi        string
	Harga_modal     string
	Harga_jual      string
	Ongkos_expedisi string
	Stock           int `json:",string"`
	Ongkos_kirim    string
	Is_active       bool
	ID              int64 `datastore:"-"`
}

const (
	CREATE = 1
	UPDATE = 2
	DELETE = -1
)

type ProductLog struct {
	UserID    int64
	ProductID int64
	Quantity  int
	LogType   int
	Message   string
	Created   time.Time
}

func init() {

}

func InitRouter(router *mux.Router) {
	router.HandleFunc("/api/products", productListHandler).Methods("GET")
	router.HandleFunc("/api/products", productCreateHandler).Methods("POST")
	router.HandleFunc("/api/products/{product_id}", productDeleteHandler).Methods("DELETE")
}

func productListHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	err := auth.ValidateSession(c, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := datastore.NewQuery("Products").Filter("Is_active =", true).Limit(20)
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

	UserId, err := auth.GetUserFromSession(c, r)
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

	product.Is_active = true

	// save the products
	key := datastore.NewIncompleteKey(c, "Products", nil)
	key, err = datastore.Put(c, key, &product)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write the created product back as reponse
	product.ID = key.IntID()
	response, err := json.Marshal(product)

	// Write log for product creation
	log := ProductLog{
		UserID:    UserId,
		ProductID: product.ID,
		Quantity:  product.Stock,
		Message:   "Created new product",
		LogType:   CREATE,
		Created:   time.Now(),
	}
	logKey := datastore.NewIncompleteKey(c, "ProductLogs", nil)
	_, err = datastore.Put(c, logKey, &log)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", response)

}

func productDeleteHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	err := auth.ValidateSession(c, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	UserId, err := auth.GetUserFromSession(c, r)

	//get the variable
	urlVar := mux.Vars(r)
	product_id, _ := strconv.ParseInt(urlVar["product_id"], 10, 64)

	// Construct key and delete stuff
	key := datastore.NewKey(c, "Products", "", product_id, nil)

	// Get the product
	var product Product
	err = datastore.Get(c, key, &product)
	if err != nil {
		http.Error(w, "Problem with database", http.StatusInternalServerError)
	}

	// mark it as deleted, soft-delete
	product.Is_active = false
	productKey, err := datastore.Put(c, key, &product)
	if err != nil {
		http.Error(w, "Cannot delete product", http.StatusBadRequest)
		return
	}

	// Write log
	productLog := ProductLog{
		UserID:    UserId,
		ProductID: productKey.IntID(),
		Quantity:  0,
		Message:   "DELETED Product",
		LogType:   DELETE,
		Created:   time.Now(),
	}
	logKey := datastore.NewIncompleteKey(c, "ProductLogs", nil)
	_, err = datastore.Put(c, logKey, &productLog)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func handleError(res http.ResponseWriter, err error) {
	if err != nil {
		log.Println("=====ERROR=====")
		log.Println(err.Error())
		log.Println("=====END ERROR=====")
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

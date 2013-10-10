package auth

import (
	"appengine"
	"appengine/memcache"
	"encoding/json"
	"errors"
	mux "github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"strings"
	"time"
)

// session
var store = sessions.NewCookieStore([]byte("something-very-secret"))

func init() {

}

func InitRouter(router *mux.Router) {
	router.HandleFunc("/login", loginHandler)
	router.HandleFunc("/logout", logoutHandler)
	router.HandleFunc("/register", registerHandler)
}

// function to validate sessions authentication
func ValidateSession(c appengine.Context, r *http.Request) error {
	session, _ := store.Get(r, "session-inventory")
	token, _ := session.Values["token"].(string)
	ip := strings.Split(r.RemoteAddr, ":")[0]
	result, err := memcache.Get(c, "sessions:"+token)
	if err == memcache.ErrCacheMiss {
		return errors.New("validate-session: token not recognized")
	}

	if string(result.Value) != ip {
		return errors.New("validate-session: token is from different IP")
	}

	return nil
}

func GetUserFromSession(c appengine.Context, r *http.Request) (int64, error) {
	session, _ := store.Get(r, "session-inventory")
	UserID, _ := session.Values["UserID"].(int64)
	log.Println(session.Values)
	return UserID, nil
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-inventory")
	session.Values["token"] = ""
	session.Save(r, w)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// decode request
	decoder := json.NewDecoder(r.Body)
	var user User
	err := decoder.Decode(&user)
	if err != nil {
		log.Println(err.Error())
	}

	// try loggin in
	result, err := user.Authenticate(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// create session token with username, ip, secret key
	token := getAuthToken(result.Username)
	session, _ := store.Get(r, "session-inventory")
	session.Values["token"] = token
	session.Values["UserID"] = result.ID
	session.Save(r, w)

	// store in memcache
	expiration, _ := time.ParseDuration("5m")
	item := &memcache.Item{
		Key:        "sessions:" + token,
		Value:      []byte(strings.Split(r.RemoteAddr, ":")[0]),
		Expiration: expiration,
	}

	// Add the item to the memcache, if the key does not already exist
	if err := memcache.Add(c, item); err == memcache.ErrNotStored {
		c.Infof("item with key %q already exists", item.Key)
	} else if err != nil {
		c.Errorf("error adding item: %v", err)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// decode request
	decoder := json.NewDecoder(r.Body)
	var user User

	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Try to register, if cannot, should return appropriate header
	user, err = user.Register(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

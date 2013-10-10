package auth

import (
	"appengine"
	"appengine/datastore"
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"time"
)

type User struct {
	Username string
	Password string
	ID       int64 `datastore="-"`
}

const COST = 12

// Generate authentication token
func getAuthToken(username string) string {
	now := time.Now().Format("20060102150405")
	return username + now
}

// get User struct from datastore based on username
func getUserFromUsername(c appengine.Context, username string) (User, error) {
	user := User{}

	// get user based on Username
	query := datastore.NewQuery("Users").Filter("Username =", username)

	var users []User
	keys, err := query.GetAll(c, &users)
	if err != nil {
		return user, err
	}
	if len(users) > 1 {
		return user, errors.New("More than one user is returned, something bad happened")
	}

	if len(users) == 0 {
		return User{}, nil
	}
	user = users[0]
	user.ID = keys[0].IntID()
	return user, nil
}

// Authenticate function. Take a clear text password and match with the hashed password
// return nil on success, error on failure
func (user User) Authenticate(c appengine.Context) (User, error) {

	potential_user, err := getUserFromUsername(c, user.Username)
	if err != nil {
		return user, err
	}

	if potential_user != (User{}) {
		return potential_user, nil
	}

	// Check password match
	err = bcrypt.CompareHashAndPassword([]byte(potential_user.Password), []byte(user.Password))
	if err != nil {
		return user, errors.New("Username / Password Invalid")
	}
	return potential_user, nil
}

// Function to register new user. Take a clear text password and username
// Hash the password and then save into the datastore
func (user User) Register(c appengine.Context) (User, error) {

	// If existing user return error
	potential_user, err := getUserFromUsername(c, user.Username)
	if err != nil {
		return user, err
	}

	if potential_user != (User{}) {
		return user, errors.New("User with this username exists")
	}

	hashed_password, err := bcrypt.GenerateFromPassword([]byte(user.Password), COST)
	if err != nil {
		return user, err
	}
	user.Password = string(hashed_password)

	// save the user
	key := datastore.NewIncompleteKey(c, "Users", nil)
	_, err = datastore.Put(c, key, &user)

	if err != nil {
		return user, err
	}

	return user, nil
}

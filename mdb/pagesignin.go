package mdb

import (
	"encoding/json"
	"net/http"
)

func pageSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		pageSignInGet(w, r)
	} else if r.Method == http.MethodPost {
		pageSignInPost(w, r)
	}
}

func pageSignInGet(w http.ResponseWriter, r *http.Request) {
	Unauthenticate(w, r)
	templates.ExecuteTemplate(w, "signin", nil)
}

func pageSignInPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	s := OpenSchema()
	defer s.Close()
	practitionerID, hashedPassword := s.practitionerTable.FindPractitioner(username)
	if !CheckPassword(hashedPassword, password) {
		panic("Incorrect password")
	}
	a := Authentication{
		Authenticated:        true,
		PractitionerID:       practitionerID,
		PractitionerUsername: username,
	}
	Authenticate(w, r, a)
	http.Redirect(w, r, "/", 302)
}

type Authentication struct {
	Authenticated        bool
	PractitionerID       string
	PractitionerUsername string
}

func IsAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	authentication, err := store.Get(r, "authentication")
	if err != nil {
		return false
	}
	j := authentication.Values["authentication"]
	if j == nil {
		return false
	}
	var a Authentication
	err = json.Unmarshal(j.([]byte), &a)
	return a.Authenticated
}

func GetAuthentication(w http.ResponseWriter, r *http.Request) *Authentication {
	authentication, err := store.Get(r, "authentication")
	if err != nil {
		return nil
	}
	j := authentication.Values["authentication"]
	if j == nil {
		return nil
	}
	var a Authentication
	err = json.Unmarshal(j.([]byte), &a)
	return &a
}

func Authenticate(w http.ResponseWriter, r *http.Request, a Authentication) {
	authentication, err := store.Get(r, "authentication")
	if err != nil {
		panic(err)
	}
	j, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	authentication.Values["authentication"] = j
	authentication.Save(r, w)
}

func Unauthenticate(w http.ResponseWriter, r *http.Request) {
	Authenticate(w, r, Authentication{Authenticated: false})
}

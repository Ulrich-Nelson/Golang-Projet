package mdb

import (
	"net/http"
)

func pageIndex(w http.ResponseWriter, r *http.Request) {
	if IsAuthenticated(w, r) == false {
		http.Redirect(w, r, "/signin", 302)
		return
	}

	data := struct {
		Title string
		Dates []string
	}{}
	data.Title = "MDB Web Server"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := templates.ExecuteTemplate(w, "index", data)
	if err != nil {
		panic(err)
	}
}

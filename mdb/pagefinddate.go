package mdb

import (
	"net/http"
)

func pageFindDate(w http.ResponseWriter, r *http.Request) {
	dateString := r.FormValue("date") // Unix timestamp

	s := OpenSchema()
	defer s.Close()
	data := struct {
		Title     string
		Date      string
		Protocols []Protocol
	}{
		Title:     "MDB Web Server",
		Date:      dateString,
		Protocols: s.protocolTable.List(),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := templates.ExecuteTemplate(w, "finddate", data)
	if err != nil {
		panic(err)
	}
}

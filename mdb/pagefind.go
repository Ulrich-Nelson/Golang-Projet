package mdb

import (
	"net/http"
	"time"
)

func pageFind(w http.ResponseWriter, r *http.Request) {
	s := OpenSchema()
	defer s.Close()
	dates := s.experimentTable.QueryDates()
	data := struct {
		Title string
		Dates []time.Time
	}{
		Title: "MDB Web Server",
		Dates: dates,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := templates.ExecuteTemplate(w, "find", data)
	if err != nil {
		panic(err)
	}
}

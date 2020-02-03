package mdb

import (
	"net/http"
)

func pageFindExperiments(w http.ResponseWriter, r *http.Request) {
	protocolName := r.FormValue("protocolName")
	dateString := r.FormValue("date") // Unix timestamp

	if !ValidateName(protocolName) {
		panic("Invalid protocolName")
	}
	date := ParseUnixDate(dateString)

	s := OpenSchema()
	defer s.Close()
	experiments := s.protocolTable.QueryExperiments(protocolName, date)
	data := struct {
		Title        string
		ProtocolName string
		Date         string
		Experiments  []PrettyExperiment
	}{
		Title:        "MDB Web Server",
		ProtocolName: protocolName,
		Date:         dateString,
		Experiments:  experiments,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := templates.ExecuteTemplate(w, "findexperiments", data)
	if err != nil {
		panic(err)
	}
}

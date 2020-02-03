package mdb

import "net/http"

func pageDeposit(w http.ResponseWriter, r *http.Request) {
	s := OpenSchema()
	defer s.Close()
	data := struct {
		Title          string
		Authentication *Authentication
		Patients       []Patient
		Protocols      []Protocol
	}{
		Title:          "MDB Web Server",
		Authentication: GetAuthentication(w, r),
		Patients:       s.patientTable.List(),
		Protocols:      s.protocolTable.List(),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := templates.ExecuteTemplate(w, "deposit", data)
	if err != nil {
		panic(err)
	}
}

package mdb

import "net/http"

func pagePatients(w http.ResponseWriter, r *http.Request) {
	s := OpenSchema()
	defer s.Close()
	data := struct {
		Title    string
		Patients []Patient
	}{
		Title:    "MDB Web Server",
		Patients: s.patientTable.List(),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := templates.ExecuteTemplate(w, "patients", data)
	if err != nil {
		panic(err)
	}
}

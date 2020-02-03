package mdb

import (
	"net/http"
)

func pageInspectPatient(w http.ResponseWriter, r *http.Request) {
	patientID := r.FormValue("patientID")

	if !ValidateID(patientID) {
		panic("Invalid patientID")
	}

	s := OpenSchema()
	defer s.Close()
	patient := s.patientTable.FindPatient(patientID)
	data := struct {
		Title       string
		Patient     PrettyPatient
		Experiments []PrettyExperiment
	}{
		Title:       "MDB Web Server",
		Patient:     patient,
		Experiments: s.experimentTable.QueryPatientExperiments(patientID),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := templates.ExecuteTemplate(w, "inspectpatient", data)
	if err != nil {
		panic(err)
	}
}

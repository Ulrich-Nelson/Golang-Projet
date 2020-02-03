package mdb

import (
	"io/ioutil"
	"net/http"
)

func pageUpload(w http.ResponseWriter, r *http.Request) {
	practitionerID := r.FormValue("practitionerID")
	patientID := r.FormValue("patientID")
	experimentID := r.FormValue("experimentID")
	packetID := r.FormValue("packetID")
	protocolName := r.FormValue("protocolName")

	if !ValidateID(practitionerID) {
		panic("Invalid practitionerID")
	}
	if !ValidateID(patientID) {
		panic("Invalid patientID")
	}
	if experimentID != "" && !ValidateID(experimentID) {
		panic("Invalid experimentID")
	}
	if packetID != "" && !ValidateID(packetID) {
		panic("Invalid packetID")
	}
	if !ValidateName(protocolName) {
		panic("Invalid protocolName")
	}

	file, _, err := r.FormFile("csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	tempFile, err := ioutil.TempFile("csv", "csv-*")
	if err != nil {
		panic(err)
	}
	defer tempFile.Close()
	tempFile.Write(fileBytes)
	s := OpenSchema()
	defer s.Close()
	if experimentID != "" {
		e := &Experiment{
			ID:             experimentID,
			Timestamping:   GenerateTimestamping(),
			PractitionerID: practitionerID,
			PatientID:      patientID,
			PacketID:       packetID,
		}
		s.experimentTable.AddExperiment(e)
	}
	s.injection.Inject(protocolName, tempFile.Name())
	data := struct {
		Title string
	}{
		Title: "MDB Web Server",
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = templates.ExecuteTemplate(w, "uploaded", data)
	if err != nil {
		panic(err)
	}
}

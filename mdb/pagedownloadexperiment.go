package mdb

import "net/http"

func pageDownloadExperiment(w http.ResponseWriter, r *http.Request) {
	protocolName := r.FormValue("protocolName")
	experimentID := r.FormValue("experimentID")

	if !ValidateName(protocolName) {
		panic("Invalid protocolName")
	}
	if !ValidateID(experimentID) {
		panic("Invalid experimentID")
	}

	s := OpenSchema()
	defer s.Close()
	dataSet := s.protocolTable.QueryExperiment(protocolName, experimentID)
	csv := dataSet.exportCSV()

	w.Header().Set("Content-Type", "application/octet-stream; charset=utf-8")
	w.Write([]byte(csv))
}

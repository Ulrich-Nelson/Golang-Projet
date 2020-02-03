package mdb

import (
	"fmt"
	"net/http"
	"strconv"
)

func CreateWebService(port int) {
	http.HandleFunc("/id", handleID)
	fmt.Printf("Starting web server at port %d\n", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		panic(err)
	}
}

func handleID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(GenerateID()))
}

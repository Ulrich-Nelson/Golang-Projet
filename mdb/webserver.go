package mdb

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("t0p-S3cr3t"))
var templates *template.Template

func CreateWebserver(port int) {
	templates = template.Must(template.ParseGlob("mdb/templates/*.gohtml"))
	http.Handle("/js/", http.FileServer(http.Dir("mdb/static/")))
	http.Handle("/css/", http.FileServer(http.Dir("mdb/static/")))
	http.Handle("/images/", http.FileServer(http.Dir("mdb/static/")))
	http.HandleFunc("/", pageIndex)
	http.HandleFunc("/signin", pageSignIn)
	http.HandleFunc("/patients", pagePatients)
	http.HandleFunc("/inspectpatient", pageInspectPatient)
	http.HandleFunc("/find", pageFind)
	http.HandleFunc("/finddate", pageFindDate)
	http.HandleFunc("/findexperiments", pageFindExperiments)
	http.HandleFunc("/downloadexperiment", pageDownloadExperiment)
	http.HandleFunc("/deposit", pageDeposit)
	http.HandleFunc("/upload", pageUpload)
	fmt.Printf("Starting web server at port %d\n", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		panic(err)
	}
}

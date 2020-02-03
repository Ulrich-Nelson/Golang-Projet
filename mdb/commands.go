package mdb

import (
	"flag"
	"fmt"
	"os"
)

func ExecuteCommand() {
	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments")
		return
	}
	switch os.Args[1] {
	case "clear":
		doClear()
	case "id":
		doID()
	case "init":
		doInit()
	case "create-protocol":
		doCreateProtocol()
	case "create-experiment":
		doCreateExperiment()
	case "create-practitioner":
		doCreatePractitioner()
	case "create-patient":
		doCreatePatient()
	case "inject":
		doInject()
	case "experiments":
		doExperiments()
	case "webserver":
		doWebServer()
	case "webservice":
		doWebService()
	default:
		fmt.Println("Unknown command ", os.Args[1])
	}
}

func doClear() {
	s := OpenSchema()
	defer s.Close()
	s.Drop()
	fmt.Println("Schema dropped")
}

func doID() {
	fmt.Println(GenerateID())
}

func doInit() {
	s := CreateSchema()
	defer s.Close()
	fmt.Println("Schema initialized")
}

func doCreateProtocol() {
	fs := flag.NewFlagSet("create-protocol", flag.ExitOnError)
	var protocolName string
	fs.StringVar(&protocolName, "protocolName", "", "name of the protocol")
	var textFieldNames arrayFlags
	fs.Var(&textFieldNames, "text", "column name of text type")
	var realFieldNames arrayFlags
	fs.Var(&realFieldNames, "real", "column name of real type")
	fs.Parse(os.Args[2:])

	if !ValidateName(protocolName) {
		fs.PrintDefaults()
		return
	}
	var fields []Field
	for _, textFieldName := range textFieldNames {
		if !ValidateName(textFieldName) {
			fs.PrintDefaults()
			return
		}
		field := Field{FieldName: textFieldName, FieldType: TextFieldType}
		fields = append(fields, field)
	}
	for _, realFieldName := range realFieldNames {
		if !ValidateName(realFieldName) {
			fs.PrintDefaults()
			return
		}
		field := Field{FieldName: realFieldName, FieldType: RealFieldType}
		fields = append(fields, field)
	}

	s := OpenSchema()
	defer s.Close()
	p := &Protocol{Name: protocolName, Fields: fields}
	s.protocolTable.AddProtocol(p)
	fmt.Printf("Protocol %v created\n", protocolName)
}

func doCreateExperiment() {
	fs := flag.NewFlagSet("create-experiment", flag.ExitOnError)
	experimentID := fs.String("experimentID", "", "experiment id")
	practitionerID := fs.String("practitionerID", "", "practitioner id")
	patientID := fs.String("patientID", "", "patient id")
	packetID := fs.String("packetID", "", "packet id")
	fs.Parse(os.Args[2:])

	if *experimentID == "" {
		*experimentID = GenerateID()
	}
	if !ValidateID(*experimentID) ||
		!ValidateID(*practitionerID) ||
		!ValidateID(*patientID) ||
		(*packetID != "" && !ValidateID(*packetID)) {
		fs.PrintDefaults()
		return
	}

	s := OpenSchema()
	defer s.Close()
	e := &Experiment{
		ID:             *experimentID,
		Timestamping:   GenerateTimestamping(),
		PractitionerID: *practitionerID,
		PatientID:      *patientID,
		PacketID:       *packetID,
	}
	s.experimentTable.AddExperiment(e)
	fmt.Printf("Experiment %v created\n", *experimentID)
}

func doCreatePractitioner() {
	fs := flag.NewFlagSet("create-practitioner", flag.ExitOnError)
	practitionerID := fs.String("practitionerID", "", "practitioner id")
	username := fs.String("username", "", "name of the practitioner")
	password := fs.String("password", "", "password of the practitioner")
	fs.Parse(os.Args[2:])

	if *practitionerID == "" {
		*practitionerID = GenerateID()
	}
	if !ValidateID(*practitionerID) ||
		!ValidateName(*username) {
		fs.PrintDefaults()
		return
	}

	if *password == "" {
		fmt.Printf("Enter password: ")
		*password = ReadPassword()
	}
	if !ValidatePassword(*password) {
		fmt.Println("Password must be maximum 32 characters printable ascii")
		return
	}

	s := OpenSchema()
	defer s.Close()
	p := s.practitionerTable.AddPractitioner(
		*practitionerID,
		*username,
		*password,
	)
	fmt.Printf("Practitioner %v (%v) created\n", p.id, p.username)
}

func doCreatePatient() {
	fs := flag.NewFlagSet("add-patient", flag.ExitOnError)
	patientID := fs.String("patientID", "", "patient id")
	practitionerID := fs.String("practitionerID", "", "practitioner id")
	data := fs.String("data", "", "data of the patient")
	password := fs.String("password", "", "password for encryption")
	fs.Parse(os.Args[2:])

	if *patientID == "" {
		*patientID = GenerateID()
	}
	if !ValidateID(*patientID) || !ValidateID(*practitionerID) {
		fs.PrintDefaults()
		return
	}

	if *password == "" {
		fmt.Printf("Enter password for encryption: ")
		*password = ReadPassword()
	}
	if !ValidatePassword(*password) {
		fmt.Println("Password must be maximum 32 characters printable ascii")
		return
	}
	dataEncrypted := Encrypt(*password, *data)

	s := OpenSchema()
	defer s.Close()
	p := s.patientTable.AddPatient(*patientID, *practitionerID, dataEncrypted)
	fmt.Printf("Patient %v created\n", p.ID)
}

func doInject() {
	fs := flag.NewFlagSet("inject", flag.ExitOnError)
	protocolName := fs.String("protocolName", "", "protocol name")
	filename := fs.String("filename", "", "filename")
	fs.Parse(os.Args[2:])

	if !ValidateName(*protocolName) {
		fs.PrintDefaults()
		return
	}

	s := OpenSchema()
	defer s.Close()
	s.injection.Inject(*protocolName, *filename)
	fmt.Printf("File %v injected into %v\n", *filename, *protocolName)
}

func doExperiments() {
	s := OpenSchema()
	defer s.Close()
	experiments := s.experimentTable.List()
	for _, e := range experiments {
		fmt.Println(e.Timestamping.String())
	}
}

func doWebServer() {
	CreateWebserver(8080)
}

func doWebService() {
	CreateWebService(9090)
}

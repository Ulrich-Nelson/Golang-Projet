package mdb

type Schema struct {
	session           *Session
	experimentTable   *ExperimentTable
	protocolTable     *ProtocolTable
	practitionerTable *PractitionerTable
	patientTable      *PatientTable
	injection         *Injection
}

func CreateSchema() *Schema {
	s := &Schema{}
	s.session = NewSession()
	s.dbCreateSchema()
	s.experimentTable = CreateExperimentTable(s)
	s.protocolTable = CreateProtocolTable(s)
	s.practitionerTable = CreatePractitionerTable(s)
	s.patientTable = CreatePatientTable(s)
	s.injection = NewInjection(s)
	return s
}

func OpenSchema() *Schema {
	s := &Schema{}
	s.session = NewSession()
	s.experimentTable = OpenExperimentTable(s)
	s.protocolTable = OpenProtocolTable(s)
	s.practitionerTable = OpenPractitionerTable(s)
	s.patientTable = OpenPatientTable(s)
	s.injection = NewInjection(s)
	return s
}

func (s *Schema) Close() {
	defer s.session.Close()
}

func (s *Schema) Drop() {
	dbDropSchema(s.session)
}

func (s *Schema) dbCreateSchema() {
	sql := "CREATE SCHEMA mdb"
	statement, err := s.session.db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}
}

func dbDropSchema(s *Session) {
	sql := "DROP SCHEMA IF EXISTS mdb CASCADE"
	statement, err := s.db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}
}

package mdb

type PatientTable struct {
	schema *Schema
}

type Patient struct {
	ID             string
	PractitionerID string
	Profile        string
}

type PrettyPatient struct {
	Patient
	PractitionerUsername string
}

func CreatePatientTable(s *Schema) *PatientTable {
	pt := &PatientTable{schema: s}
	pt.dbCreateTable()
	return pt
}

func OpenPatientTable(s *Schema) *PatientTable {
	return &PatientTable{schema: s}
}

func (pt *PatientTable) AddPatient(
	id string,
	practitionerID string,
	profile string,
) *Patient {
	pt.dbInsert(id, practitionerID, profile)
	return &Patient{ID: id, PractitionerID: practitionerID, Profile: profile}
}

func (pt *PatientTable) List() []Patient {
	q := "SELECT id, practitioner_id, profile FROM Patient"
	statement, err := pt.schema.session.db.Prepare(q)
	if err != nil {
		panic(err)
	}
	rows, err := statement.Query()
	defer rows.Close()
	var patients []Patient
	for rows.Next() {
		var p Patient
		err = rows.Scan(&p.ID, &p.PractitionerID, &p.Profile)
		if err != nil {
			panic(err)
		}
		patients = append(patients, p)
	}
	return patients
}

func (pt *PatientTable) FindPatient(patientID string) PrettyPatient {
	q := `
	SELECT
		mdb.Patient.id,
		mdb.Patient.practitioner_id,
		mdb.Patient.profile,
		mdb.Practitioner.username
	FROM mdb.Patient
	INNER JOIN mdb.Practitioner
	ON mdb.Practitioner.id = practitioner_id
	WHERE mdb.Patient.id = $1
	`
	statement, err := pt.schema.session.db.Prepare(q)
	if err != nil {
		panic(err)
	}
	row := statement.QueryRow(patientID)
	var patient PrettyPatient
	err = row.Scan(
		&patient.ID,
		&patient.PractitionerID,
		&patient.Profile,
		&patient.PractitionerUsername,
	)
	if err != nil {
		panic(err)
	}
	return patient
}

func (pt *PatientTable) dbCreateTable() {
	q := `
	CREATE TABLE IF NOT EXISTS mdb.Patient (
		id TEXT PRIMARY KEY NOT NULL,
		practitioner_id TEXT NOT NULL,
		profile TEXT NOT NULL
	)
	`
	statement, err := pt.schema.session.db.Prepare(q)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}
}

func (pt *PatientTable) dbInsert(id string, practitionerID string, profile string) {
	q := `
	INSERT INTO mdb.Patient (
		id, practitioner_id, profile
	) VALUES (
		$1, $2, $3
	)
	`
	statement, err := pt.schema.session.db.Prepare(q)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec(id, practitionerID, profile)
	if err != nil {
		panic(err)
	}
}

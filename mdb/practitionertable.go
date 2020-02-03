package mdb

type PractitionerTable struct {
	schema *Schema
}

type Practitioner struct {
	id       string
	username string
}

func CreatePractitionerTable(s *Schema) *PractitionerTable {
	pt := &PractitionerTable{schema: s}
	pt.dbCreateTable()
	return pt
}

func OpenPractitionerTable(s *Schema) *PractitionerTable {
	return &PractitionerTable{schema: s}
}

func (pt *PractitionerTable) AddPractitioner(
	id string,
	username string,
	password string,
) *Practitioner {
	hashedPassword := HashPassword(password)
	pt.dbInsert(id, username, hashedPassword)
	return &Practitioner{id: id, username: username}
}

func (pt *PractitionerTable) FindPractitioner(username string) (string, string) {
	q := `
	SELECT id, hashed_password
	FROM mdb.Practitioner
	WHERE mdb.Practitioner.username = $1
	`
	statement, err := pt.schema.session.db.Prepare(q)
	if err != nil {
		panic(err)
	}
	row := statement.QueryRow(username)
	var practitionerID string
	var hashedPassword string
	err = row.Scan(&practitionerID, &hashedPassword)
	if err != nil {
		panic(err)
	}
	return practitionerID, hashedPassword
}

func (pt *PractitionerTable) dbCreateTable() {
	q := `
	CREATE TABLE IF NOT EXISTS mdb.Practitioner (
		id TEXT PRIMARY KEY NOT NULL,
		username TEXT NOT NULL UNIQUE,
		hashed_password TEXT NOT NULL
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

func (pt *PractitionerTable) dbInsert(
	id string,
	username string,
	hashedPassword string,
) {
	q := `
	INSERT INTO mdb.Practitioner (
		id, username, hashed_password
	) VALUES (
		$1, $2, $3
	)
	`
	statement, err := pt.schema.session.db.Prepare(q)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec(id, username, hashedPassword)
	if err != nil {
		panic(err)
	}
}

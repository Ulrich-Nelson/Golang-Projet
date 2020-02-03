package mdb

import (
	"time"
)

type ExperimentTable struct {
	schema *Schema
}

type Experiment struct {
	ID             string
	Timestamping   time.Time
	PractitionerID string
	PatientID      string
	PacketID       string
}

type PrettyExperiment struct {
	Experiment
	PractitionerUsername string
	PatientProfile       string
}

func CreateExperimentTable(s *Schema) *ExperimentTable {
	et := &ExperimentTable{schema: s}
	et.dbCreateTable()
	return et
}

func OpenExperimentTable(s *Schema) *ExperimentTable {
	return &ExperimentTable{schema: s}
}

func (et *ExperimentTable) AddExperiment(e *Experiment) {
	et.dbInsert(e)
}

func (et *ExperimentTable) List() []Experiment {
	sql := `
	SELECT id, timestamping, practitioner_id, patient_id, packet_id
	FROM mdb.Experiment
	`
	statement, err := et.schema.session.db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	rows, err := statement.Query()
	defer rows.Close()
	var experiments []Experiment
	for rows.Next() {
		var e Experiment
		err = rows.Scan(
			&e.ID,
			&e.Timestamping,
			&e.PractitionerID,
			&e.PatientID,
			&e.PacketID,
		)
		if err != nil {
			panic(err)
		}
		experiments = append(experiments, e)
	}
	return experiments
}

func (et *ExperimentTable) QueryPatientExperiments(
	patientID string,
) []PrettyExperiment {
	sql := `
	SELECT
		mdb.Experiment.id,
		mdb.Experiment.timestamping,
		mdb.Experiment.practitioner_id,
		mdb.Experiment.patient_id,
		mdb.Experiment.packet_id,
		mdb.Practitioner.username
	FROM mdb.Experiment
	INNER JOIN mdb.Practitioner
	ON mdb.Practitioner.id = mdb.Experiment.practitioner_id
	WHERE patient_id = $1
	`
	statement, err := et.schema.session.db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	rows, err := statement.Query(patientID)
	defer rows.Close()
	var experiments []PrettyExperiment
	for rows.Next() {
		var e PrettyExperiment
		err = rows.Scan(
			&e.ID,
			&e.Timestamping,
			&e.PractitionerID,
			&e.PatientID,
			&e.PacketID,
			&e.PractitionerUsername,
		)
		if err != nil {
			panic(err)
		}
		experiments = append(experiments, e)
	}
	return experiments
}

func (et *ExperimentTable) QueryDates() []time.Time {
	sql := `
	SELECT DATE_TRUNC('DAY', timestamping) AS date
	FROM mdb.Experiment
	GROUP BY date
	`
	statement, err := et.schema.session.db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	rows, err := statement.Query()
	defer rows.Close()
	var dates []time.Time
	for rows.Next() {
		var d time.Time
		err = rows.Scan(&d)
		if err != nil {
			panic(err)
		}
		dates = append(dates, d)
	}
	return dates
}

func (et *ExperimentTable) dbCreateTable() {
	sql := `
	CREATE TABLE IF NOT EXISTS mdb.Experiment (
		id TEXT PRIMARY KEY NOT NULL,
		timestamping TIMESTAMP NOT NULL,
		practitioner_id TEXT NOT NULL,
		patient_id TEXT NOT NULL,
		packet_id TEXT NOT NULL
	)
	`
	statement, err := et.schema.session.db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}
}

func (et *ExperimentTable) dbInsert(e *Experiment) {
	sql := `
	INSERT INTO mdb.Experiment (
		id,
		timestamping,
		practitioner_id,
		patient_id,
		packet_id
	) VALUES (
		$1, $2, $3, $4, $5
	)
	`
	statement, err := et.schema.session.db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec(
		e.ID,
		e.Timestamping,
		e.PractitionerID,
		e.PatientID,
		e.PacketID,
	)
	if err != nil {
		panic(err)
	}
}

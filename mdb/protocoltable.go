package mdb

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type ProtocolTable struct {
	schema *Schema
}

const (
	TextFieldType = "TEXT"
	RealFieldType = "REAL"
)

type Field struct {
	FieldName string `json:"field_name"`
	FieldType string `json:"field_type"`
}

type Protocol struct {
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}

func CreateProtocolTable(s *Schema) *ProtocolTable {
	pt := &ProtocolTable{schema: s}
	pt.dbCreateTable()
	return pt
}

func OpenProtocolTable(s *Schema) *ProtocolTable {
	return &ProtocolTable{schema: s}
}

func (pt *ProtocolTable) AddProtocol(p *Protocol) *Protocol {
	pt.dbInsert(p)
	pt.dbGenerate(p)
	return p
}

func (pt *ProtocolTable) List() []Protocol {
	q := "SELECT name FROM Protocol"
	statement, err := pt.schema.session.db.Prepare(q)
	if err != nil {
		panic(err)
	}
	rows, err := statement.Query()
	defer rows.Close()
	var protocols []Protocol
	for rows.Next() {
		var p Protocol
		err = rows.Scan(&p.Name)
		if err != nil {
			panic(err)
		}
		protocols = append(protocols, p)
	}
	return protocols
}

func (pt *ProtocolTable) QueryDates(protocolName string) []time.Time {
	q0 := `
	SELECT DATE_TRUNC('DAY', timestamping) AS date
	FROM mdb.Experiment
	INNER JOIN mdb.Protocol_%s
	ON mdb.Experiment.id = mdb.Protocol_%s.experiment_id
	GROUP BY date
	`
	q := fmt.Sprintf(q0, protocolName, protocolName)
	statement, err := pt.schema.session.db.Prepare(q)
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

func (pt *ProtocolTable) QueryExperiments(
	protocolName string,
	date time.Time,
) []PrettyExperiment {
	q0 := `
	SELECT
		mdb.Experiment.id,
		mdb.Experiment.timestamping,
		mdb.Experiment.packet_id,
		mdb.Practitioner.username,
		mdb.Patient.profile
	FROM mdb.Experiment
	INNER JOIN mdb.Protocol_%s
	ON mdb.Experiment.id = mdb.Protocol_%s.experiment_id
	INNER JOIN mdb.Practitioner
	ON mdb.Practitioner.id = mdb.Experiment.practitioner_id
	INNER JOIN mdb.Patient
	ON mdb.Patient.id = mdb.Experiment.patient_id
	WHERE DATE_TRUNC('DAY', mdb.Experiment.timestamping) = $1
	GROUP BY mdb.Experiment.id, mdb.Practitioner.username, mdb.Patient.profile
	`
	q := fmt.Sprintf(q0, protocolName, protocolName)
	statement, err := pt.schema.session.db.Prepare(q)
	if err != nil {
		panic(err)
	}
	dateTrunc := ParseTruncDate(FormatTruncDate(date))
	rows, err := statement.Query(dateTrunc)
	defer rows.Close()
	var experiments []PrettyExperiment
	for rows.Next() {
		var e PrettyExperiment
		err := rows.Scan(
			&e.ID,
			&e.Timestamping,
			&e.PacketID,
			&e.PractitionerUsername,
			&e.PatientProfile,
		)
		if err != nil {
			panic(err)
		}
		experiments = append(experiments, e)
	}
	return experiments
}

func (pt *ProtocolTable) QueryExperiment(
	protocolName string,
	experimentID string,
) *DataSet {
	q0 := "SELECT * FROM mdb.Protocol_%s WHERE experiment_id = $1"
	q := fmt.Sprintf(q0, protocolName)
	statement, err := pt.schema.session.db.Prepare(q)
	if err != nil {
		panic(err)
	}
	rows, err := statement.Query(experimentID)
	defer rows.Close()
	var columns []string
	columns, err = rows.Columns()
	if err != nil {
		panic(err)
	}
	raw := make([][]byte, len(columns))
	unknowns := make([]interface{}, len(columns))
	values := make([]string, len(columns))
	dataSet := CreateDataSet(columns)
	for i := range raw {
		unknowns[i] = &raw[i]
	}
	for rows.Next() {
		err = rows.Scan(unknowns...)
		if err != nil {
			panic(err)
		}
		for i, bytes := range raw {
			if bytes == nil {
				values[i] = "\\N"
			} else {
				values[i] = string(bytes)
			}
		}
		dataSet.Append(values)
	}
	return dataSet
}

func (pt *ProtocolTable) dbCreateTable() {
	q := `
	CREATE TABLE IF NOT EXISTS mdb.Protocol (
		name TEXT NOT NULL PRIMARY KEY,
		fields JSON NOT NULL
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

func (pt *ProtocolTable) dbInsert(p *Protocol) {
	q := `
	INSERT INTO mdb.Protocol (
		name,
		fields
	) VALUES (
		$1, $2
	)
	`
	statement, err := pt.schema.session.db.Prepare(q)
	if err != nil {
		panic(err)
	}
	var jsonFields []byte
	jsonFields, err = json.Marshal(p.Fields)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec(p.Name, string(jsonFields))
	if err != nil {
		panic(err)
	}
}

func (pt *ProtocolTable) dbGenerate(p *Protocol) {
	columns := make([]string, len(p.Fields)+1)
	columns[0] = "experiment_id TEXT NOT NULL REFERENCES mdb.Experiment(id)"
	for i := 0; i < len(p.Fields); i++ {
		columns[i+1] =
			p.Fields[i].FieldName + " " + p.Fields[i].FieldType
	}
	q := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS mdb.Protocol_%s (%s)",
		p.Name,
		strings.Join(columns, ", "),
	)
	statement, err := pt.schema.session.db.Prepare(q)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}
}

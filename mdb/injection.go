package mdb

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Injection struct {
	schema *Schema
}

func NewInjection(s *Schema) *Injection {
	inj := &Injection{schema: s}
	return inj
}

func (inj *Injection) Inject(protocolName string, csvFilename string) {
	reader := getReader(csvFilename)
	columns := getColumns(reader)
	placeholders := getPlaceholders(len(columns))
	values := make([]interface{}, len(columns))
	sql := getSQL(protocolName, columns, placeholders)
	for {
		record := getRecord(reader)
		if record == nil {
			break
		}
		for i, v := range record {
			values[i] = v
		}
		inj.insert(sql, values)
	}
}

func getReader(filename string) *csv.Reader {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	return reader
}

func getColumns(r *csv.Reader) []string {
	columns, err := r.Read()
	if err != nil {
		panic(err)
	}
	return columns
}

func getPlaceholders(n int) []string {
	placehodlers := make([]string, n)
	for i := range placehodlers {
		placehodlers[i] = "$" + strconv.Itoa(i+1)
	}
	return placehodlers
}

func getSQL(protocolName string, columns []string, placeholders []string) string {
	sql := fmt.Sprintf(
		"INSERT INTO mdb.Protocol_%s (%s) VALUES (%s)",
		protocolName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)
	return sql
}

func getRecord(r *csv.Reader) []string {
	record, err := r.Read()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		panic(err)
	}
	return record
}

func (inj *Injection) insert(sql string, values []interface{}) {
	statement, err := inj.schema.session.db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec(values...)
	if err != nil {
		panic(err)
	}
}

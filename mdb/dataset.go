package mdb

import (
	"bytes"
	"strings"
)

type DataSet struct {
	headers []string
	cells   [][]string
}

func CreateDataSet(headers []string) *DataSet {
	return &DataSet{
		headers: headers,
	}
}

func (g *DataSet) Append(values []string) {
	array := make([]string, len(values))
	copy(array, values)
	g.cells = append(g.cells, array)
}

func (g *DataSet) exportCSV() string {
	var buffer bytes.Buffer
	buffer.WriteString(strings.Join(g.headers, ";"))
	buffer.WriteString("\n")
	for i := 0; i < len(g.cells); i++ {
		line := strings.Join(g.cells[i], ";")
		buffer.WriteString(line)
		buffer.WriteString("\n")
	}
	return buffer.String()
}

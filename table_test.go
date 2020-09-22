package table

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTable_WriteTable(t *testing.T) {
	r := require.New(t)

	var buf bytes.Buffer

	// write the table
	tab := Table{
		Headers: []string{"something", "another"},
		Rows: [][]string{
			{"1", "2"},
			{"3", "4"},
			{"3", "a longer piece of text that should stretch"},
			{"but this one is longer", "shorter now"},
		},
	}
	err := tab.WriteTable(&buf)
	r.NoError(err)

	lines := strings.Split(buf.String(), "\n")
	r.Equal(6, len(lines), "the output should have %v lines", 6)

	for _, header := range tab.Headers {
		r.Contains(lines[0], header, "header column should contain %v", header)
	}

	for i, row := range tab.Rows {
		line := lines[i+1]
		for _, v := range row {
			r.Contains(line, v, "line %v should contain %v", i+1, v)
		}
	}
}

func TestTable_WriteLargeTable(t *testing.T) {
	r := require.New(t)

	var buf bytes.Buffer

	// write the table
	tab := Table{
		Headers: []string{"something"},
		Rows:    [][]string{},
	}

	for i := 0; i < 200; i++ {
		tab.Rows = append(tab.Rows, []string{"x"})
	}

	err := tab.WriteTable(&buf)
	r.NoError(err)
	r.Contains(buf.String(), "[ ID]", "id header should have a space in it")
}

func TestTable_WriteEmptyTable(t *testing.T) {
	r := require.New(t)

	var buf bytes.Buffer

	// write the table
	tab := Table{
		Headers: []string{"something"},
		Rows:    [][]string{},
	}
	err := tab.WriteTable(&buf)
	r.NoError(err)
	r.Equal(2, len(strings.Split(buf.String(), "\n")))
}

package table

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestingTable() Table {
	return Table{
		Headers: []string{"something", "another"},
		Rows: [][]string{
			{"1", "2"},
			{"3", "4"},
			{"3", "a longer piece of text that should stretch"},
			{"but this one is longer", "shorter now"},
		},
	}
}

func TestTable_WriteTable(t *testing.T) {
	r := require.New(t)

	var buf bytes.Buffer

	// write the table
	tab := TestingTable()
	err := tab.WriteTable(&buf, nil)
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
	tab := TestingTable()
	for i := 0; i < 200; i++ {
		tab.Rows = append(tab.Rows, []string{"x"})
	}

	err := tab.WriteTable(&buf, nil)
	r.NoError(err)
	r.Contains(buf.String(), "[ ID]", "id header should have a space in it")
}

func TestTable_WriteEmptyTable(t *testing.T) {
	r := require.New(t)

	var buf bytes.Buffer

	// write the table
	tab := TestingTable()
	err := tab.WriteTable(&buf, nil)
	r.NoError(err)
	r.Equal(6, len(strings.Split(buf.String(), "\n")))
}

func TestTable_WriteColorNoAlts(t *testing.T) {
	r := require.New(t)

	var buf bytes.Buffer

	// write the table
	tab := TestingTable()
	c := Config{
		ShowIndex:       true,
		Color:           true,
		AlternateColors: true,
		AltColorCodes:   nil,
		TitleColorCode:  "",
	}
	err := tab.WriteTable(&buf, &c)
	r.NoError(err)
	r.Equal(6, len(strings.Split(buf.String(), "\n")))
}

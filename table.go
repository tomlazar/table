package table

import (
	"fmt"
	"io"
	"math"

	"github.com/mattn/go-runewidth"
	"github.com/mgutz/ansi"
)

// Config is the
type Config struct {
	ShowIndex       bool     // shows the index/row number as the first column
	Color           bool     // use the color codes in the output
	AlternateColors bool     // alternate the colors when writing
	TitleColorCode  string   // the ansi code for the title row
	AltColorCodes   []string // the ansi codes to alternate between
}

// DefaultConfig returns the default config for table, if its ever left null in a method this will be the one
// used to display the table
func DefaultConfig() *Config {
	return &Config{
		ShowIndex:       true,
		Color:           true,
		AlternateColors: true,
		TitleColorCode:  ansi.ColorCode("white+buf"),
		AltColorCodes: []string{
			"",
			"\u001b[40m",
		},
	}
}

// Table is the struct used to define the structure, this can be used from a zero state, or inferred using the
// reflection based methods
type Table struct {
	Headers []string
	Rows    [][]string
}

// WriteTable writes the defined table to the writer passed in
func (t Table) WriteTable(w io.Writer, c *Config) error {
	if c == nil {
		c = DefaultConfig()
	}

	spacing := t.spacing()

	idLen := 2

	if d := digits(len(t.Rows)); d > idLen {
		idLen = d
	}

	if c.Color {
		fmt.Fprint(w, c.TitleColorCode)
	}
	if c.ShowIndex {
		fmt.Fprintf(w, " [%*v]  ", idLen, "ID")
	}
	for i, header := range t.Headers {
		fmt.Fprintf(w, "  %s", runewidth.FillRight(header, spacing[i]))
	}
	if c.Color {
		fmt.Fprint(w, ansi.Reset)
	}
	fmt.Fprintln(w)

	color := c.Color && c.AlternateColors && len(c.AltColorCodes) > 1
	for n, row := range t.Rows {
		if color {
			fmt.Fprint(w, c.AltColorCodes[n%len(c.AltColorCodes)])
		}
		if c.ShowIndex {
			fmt.Fprintf(w, " [%*v]  ", idLen, n)
		}
		for i, v := range row {
			fmt.Fprintf(w, "  %s", runewidth.FillRight(v, spacing[i]))
		}
		if color {
			fmt.Fprint(w, ansi.Reset)
		}
		fmt.Fprintln(w)
	}

	return nil
}

func (t Table) spacing() []int {
	s := make([]int, len(t.Headers))

	for i, header := range t.Headers {
		s[i] = runewidth.StringWidth(header)
	}

	for _, arr := range t.Rows {
		for i, v := range arr {
			if len(v) > s[i] {
				s[i] = runewidth.StringWidth(v)
			}
		}
	}

	return s
}

func digits(n int) int {
	if n == 0 {
		return 1
	}

	return int(math.Log10(float64(n))) + 1
}

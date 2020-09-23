package table

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
)

// Marshal the slince into a table format using reflection
func Marshal(arr interface{}, c *Config) ([]byte, error) {
	var buf bytes.Buffer
	err := MarshalTo(&buf, arr, c)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// MarshalTo writes the reflected table into the passed in io.Writer
func MarshalTo(w io.Writer, arr interface{}, c *Config) error {
	tab, err := parse(arr)
	if err != nil {
		return err
	}

	return tab.WriteTable(w, c)
}

// parse is the main method for refletion right now
func parse(arr interface{}) (*Table, error) {
	v := reflect.ValueOf(arr)
	switch v.Kind() {
	case reflect.Slice:
		return parseSlice(v)
	case reflect.Struct:
		return parseStruct(v)
	case reflect.Map:
		return parseMap(v)
	}

	return nil, errors.New("unknown interface type")
}

func parseMap(v reflect.Value) (*Table, error) {
	tab := &Table{
		Headers: []string{"Key", "Value"},
		Rows:    make([][]string, len(v.MapKeys())),
	}

	iter := v.MapRange()
	i := 0
	for iter.Next() {
		tab.Rows[i] = []string{fmt.Sprintf("%v", iter.Key().Interface()), fmt.Sprintf("%v", iter.Value().Interface())}

		i++
	}

	return tab, nil
}

func parseStruct(v reflect.Value) (*Table, error) {
	tab := &Table{
		Headers: []string{"Field", "Type", "Value"},
		Rows:    make([][]string, v.NumField()),
	}
	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)
		n, ok := f.Tag.Lookup("table")
		if !ok {
			n = f.Name
		}

		tab.Rows[i] = []string{n, v.Field(i).Type().Name(), fmt.Sprintf("%v", v.Field(i).Interface())}
	}

	return tab, nil
}

func parseSlice(v reflect.Value) (*Table, error) {
	if v.Len() < 1 {
		return nil, errors.New("arr must have at least one element")
	}

	head := v.Index(0)

	tab := Table{
		Headers: make([]string, head.NumField()),
		Rows:    make([][]string, v.Len()),
	}
	for i := 0; i < head.NumField(); i++ {
		f := head.Type().Field(i)
		n, ok := f.Tag.Lookup("table")
		if !ok {
			n = f.Name
		}

		tab.Headers[i] = n
	}

	for i := 0; i < v.Len(); i++ {
		ref := v.Index(i)

		tab.Rows[i] = make([]string, ref.NumField())
		for f := 0; f < ref.NumField(); f++ {
			tab.Rows[i][f] = fmt.Sprintf("%v", ref.Field(f).Interface())
		}
	}

	return &tab, nil
}

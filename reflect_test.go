package table

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshalToSlice(t *testing.T) {
	data := []struct {
		Name     string `table:"THE NAME"`
		Location string `table:"THE LOCATION"`
	}{
		{"name", "l"},
		{"namgfcxe", "asfdad"},
		{"namr3e", "l134151dsa"},
		{"namear", "lasd2135"},
	}

	var buf bytes.Buffer
	err := MarshalTo(&buf, data, nil)
	require.NoError(t, err)

	t.Log("\n" + buf.String())
}

type SomeStringer string

func (s SomeStringer) String() string {
	return "some string:" + string(s)
}

func TestMarshalStruct(t *testing.T) {
	data := struct {
		Name  string
		Ty    string
		Other SomeStringer
	}{"the name of th", "tybal", SomeStringer("diso")}

	buf, err := Marshal(data, nil)
	require.NoError(t, err)

	t.Log("\n" + string(buf))
}

func TestMarshalMap(t *testing.T) {
	data := map[string]interface{}{
		"Hello":    123,
		"Test Key": "some other key",
		"Nesting?": map[string]interface{}{"wow": "woo"},
	}

	buf, err := Marshal(data, nil)
	require.NoError(t, err)

	t.Log("\n" + string(buf))
}

package table

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshalTo(t *testing.T) {
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
}

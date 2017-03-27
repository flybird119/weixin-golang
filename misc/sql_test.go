package misc

import (
	"reflect"
	"testing"
)

func TestPgArray2StringSlice(t *testing.T) {
	cases := []struct {
		str      []byte
		expected []string
	}{
		{[]byte("{abc,de,f}"), []string{"abc", "de", "f"}},
		{[]byte("{abc}"), []string{"abc"}},
		{[]byte("{abc,}"), []string{"abc"}},
		{[]byte("{abc} "), []string{"abc"}},
		{[]byte("{abc"), []string{"abc"}},
		{[]byte("{}"), []string{}},
		{[]byte("}"), []string{}},
		{[]byte("{"), []string{}},
		{[]byte("x"), []string{"x"}},
		{nil, []string{}},
	}
	for i, c := range cases {
		actual := PgArray2StringSlice(c.str)
		if !reflect.DeepEqual(c.expected, actual) {
			t.Fatalf("case %d expected string array %#v, but got %#v", i, c.expected, actual)
		}
	}
}

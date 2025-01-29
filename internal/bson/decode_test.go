package bson

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/FaisonZ/bson/internal/bit"
)

func TestDecodeVersion(t *testing.T) {
	d := &Decoder{
		br:          bit.NewBitReader([]byte{0b0001_0000}),
		bsonVersion: 0,
	}

	err := d.decodeVersion()

	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	} else if d.bsonVersion != 1 {
		t.Errorf("Incorrect version %d, expected %d", d.bsonVersion, 1)
	}
}

func TestDecode(t *testing.T) {
	v, err := Decode([]byte{
		0b0001_0010,
		0b0011_0110,
		0b0011_0110,
		0b0110_0110,
		0b1111_0110,
		0b1111_0110,
		0b0011_0110,
		0b0010_0110,
		0b0001_0111,
		0b0010_0110,
		0b0011_0110,
		0b0010_0110,
		0b0001_0111,
		0b0010_0110,
		0b0011_0110,
		0b0110_0110,
		0b1111_0110,
		0b1111_0110,
		0b0100_0110,
		0b0010_0110,
		0b1001_0110,
		0b1110_0110,
		0b0111_1100,
	})

	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	}

	typ := reflect.TypeOf(v)
	fmt.Printf("%v\n", typ)
	o, ok := v.(map[string]any)
	if !ok {
		t.Errorf("Failed to decode the root object")
	} else if len(o) == 0 {
		t.Errorf("Should not be an empty object")
	}

	val, has := o["foo"]
	if !has {
		t.Errorf("Does not have key \"foo\"")
	}
	if str, has := val.(string); !has {
		t.Errorf("foo should contain a string")
	} else if str != "bar" {
		t.Errorf("foo should contain string \"bar\", but has %q", str)
	}

	val, has = o["bar"]
	if !has {
		t.Errorf("Does not have key \"bar\"")
	}
	if str, has := val.(string); !has {
		t.Errorf("bar should contain a string")
	} else if str != "foo" {
		t.Errorf("bar should contain string \"foo\", but has %q", str)
	}

	val, has = o["bing"]
	if !has {
		t.Errorf("Does not have key \"bar\"")
	}
	if val != nil {
		t.Errorf("bing should be nil")
	}
}

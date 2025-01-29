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

func TestDecode2(t *testing.T) {
	v, err := Decode([]byte{
		0b00010010, 0b00110110, 0b01000110, 0b00100110, 0b10010110, 0b11100110, 0b01111100, 0b11000110, 0b11000100, 0b11000010, 0b11110101, 0b01101100, 0b0110_1100, 0b0100_1100, 0b0010_1110, 0b0101_0100,
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

	val, has := o["bing"]
	if !has {
		t.Errorf("Does not have key \"bar\"")
	}
	if val != nil {
		t.Errorf("bing should be nil")
	}

	val, has = o["baz"]
	if !has {
		t.Errorf("Does not have key \"baz\"")
	}
	if b, ok := val.(bool); !ok {
		t.Errorf("baz should be a boolean")
	} else if !b {
		t.Errorf("baz should be true, but was %v", b)
	}

	val, has = o["bar"]
	if !has {
		t.Errorf("Does not have key \"bar\"")
	}
	if b, ok := val.(bool); !ok {
		t.Errorf("bar should be a boolean")
	} else if b {
		t.Errorf("bar should be false, but was %v", b)
	}
}

func TestDecode3(t *testing.T) {
	v, err := Decode([]byte{
		0b0001_0100,
		0b0101_0110,
		0b0001_0110,
		0b0001_1100,
		0b1100_0010,
		0b1100_0101,
		0b0110_1100,
		0b0100_1100,
		0b0010_1100,
		0b0100_0000,
	})

	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	}

	typ := reflect.TypeOf(v)
	fmt.Printf("%v\n", typ)
	a, ok := v.([]any)
	if !ok {
		t.Errorf("Failed to decode the root array")
	} else if len(a) == 0 {
		t.Errorf("Should not be an empty array")
	}

	if val, ok := a[0].(string); !ok {
		t.Errorf("first element should be a string")
	} else if val != "a" {
		t.Errorf("fist element should be \"a\", was %q", val)
	}

	if a[1] != nil {
		t.Errorf("second element should null")
	}

	if val, ok := a[2].(string); !ok {
		t.Errorf("third element should be a string")
	} else if val != "b" {
		t.Errorf("third element should be \"b\", was %q", val)
	}

	if val, ok := a[3].(bool); !ok {
		t.Errorf("fourth element should be a boolean")
	} else if val != true {
		t.Errorf("fourth element should be true, was %v", val)
	}

	if val, ok := a[4].(string); !ok {
		t.Errorf("fifth element should be a string")
	} else if val != "ab" {
		t.Errorf("fifth element should be \"ab\", was %q", val)
	}
}

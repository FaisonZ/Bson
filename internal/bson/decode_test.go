package bson

import (
	"fmt"
	"maps"
	"reflect"
	"slices"
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

func TestDecodeLongString(t *testing.T) {
	input := []byte{
		0b0111_1111,
		0b0110_0001,
		0b0110_0010,
		0b0110_0011,
		0b0110_0100,
		0b0110_0101,
		0b0110_0110,
		0b0110_0111,
		0b0110_1000,
		0b0110_1001,
		0b0110_1010,
		0b0110_1011,
		0b0110_1100,
		0b0110_1101,
		0b0110_1110,
		0b0110_1111,
		0b0111_0000,
		0b0111_0001,
		0b0111_0010,
		0b0111_0011,
		0b0111_0100,
		0b0111_0101,
		0b0111_0110,
		0b0111_0111,
		0b0111_1000,
		0b0111_1001,
		0b0111_1010,
		0b0011_0000,
		0b0011_0001,
		0b0011_0010,
		0b0011_0011,
		0b0011_0100,
		0b0000_1001,
		0b1010_1000,
	}
	expected := "abcdefghijklmnopqrstuvwxyz012345"

	d, _ := NewDecoder(input)

	// Toss out the String token for this test
	d.br.GetBits(3)
	got, err := d.decodeString()
	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	}

	if got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

func TestDecodeLargeObject(t *testing.T) {
	expected := map[string]any{
		"0": nil,
		"1": nil,
		"2": nil,
		"3": nil,
		"4": nil,
		"5": nil,
		"a": nil,
		"b": nil,
		"c": nil,
		"d": nil,
		"e": nil,
		"f": nil,
		"g": nil,
		"h": nil,
		"i": nil,
		"j": nil,
		"k": nil,
		"l": nil,
		"m": nil,
		"n": nil,
		"o": nil,
		"p": nil,
		"q": nil,
		"r": nil,
		"s": nil,
		"t": nil,
		"u": nil,
		"v": nil,
		"w": nil,
		"x": nil,
		"y": nil,
		"z": nil,
	}

	v, err := Decode([]byte{
		0b0001_0011, 0b1111_0110, 0b0001_0011, 0b0000_1100, 0b1100_0010,
		0b0110_0011, 0b1001_1000, 0b0100_1100, 0b1011_0011, 0b0000_1001,
		0b1001_1110, 0b0110_0001, 0b0011_0100, 0b1100_1100, 0b0010_0110,
		0b1011_1001, 0b1000_0101, 0b1000_0111, 0b0011_0000, 0b1011_0001,
		0b0110_0110, 0b0001_0110, 0b0011_1100, 0b1100_0010, 0b1100_1001,
		0b1001_1000, 0b0101_1001, 0b0111_0011, 0b0000_1011, 0b0011_0110,
		0b0110_0001, 0b0110_0111, 0b1100_1100, 0b0010_1101, 0b0001_1001,
		0b1000_0101, 0b1010_0111, 0b0011_0000, 0b1011_0101, 0b0110_0110,
		0b0001_0110, 0b1011_1100, 0b1100_0010, 0b1101_1001, 0b1001_1000,
		0b0101_1011, 0b0111_0011, 0b0000_1011, 0b0111_0110, 0b0110_0001,
		0b0110_1111, 0b1100_1100, 0b0010_1110, 0b0001_1001, 0b1000_0101,
		0b1100_0111, 0b0011_0000, 0b1011_1001, 0b0110_0110, 0b0001_0111,
		0b0011_1100, 0b1100_0010, 0b1110_1001, 0b1001_1000, 0b0101_1101,
		0b0111_0011, 0b0000_1011, 0b1011_0110, 0b0110_0001, 0b0111_0111,
		0b1100_1100, 0b0010_1111, 0b0001_1001, 0b1000_0101, 0b1110_0111,
		0b0000_0101, 0b1000_0101, 0b1110_1011, 0b0000_0000,
	})

	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	}

	o, ok := v.(map[string]any)
	if !ok {
		t.Errorf("Failed to decode the root object")
	} else if len(o) == 0 {
		t.Errorf("Should not be an empty object")
	}

	if !maps.Equal(o, expected) {
		t.Errorf("Expected:\n%v\nReceived:\n%v", expected, o)
	}
}

func TestDecodeLargeArray(t *testing.T) {
	expected := []any{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n",
		"o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "0", "1",
		"2", "3", "4", "5",
	}

	v, err := Decode([]byte{
		0b0001_0101, 0b1111_0110, 0b0001_0110, 0b0001_0110, 0b0001_0110,
		0b0010_0110, 0b0001_0110, 0b0011_0110, 0b0001_0110, 0b0100_0110,
		0b0001_0110, 0b0101_0110, 0b0001_0110, 0b0110_0110, 0b0001_0110,
		0b0111_0110, 0b0001_0110, 0b1000_0110, 0b0001_0110, 0b1001_0110,
		0b0001_0110, 0b1010_0110, 0b0001_0110, 0b1011_0110, 0b0001_0110,
		0b1100_0110, 0b0001_0110, 0b1101_0110, 0b0001_0110, 0b1110_0110,
		0b0001_0110, 0b1111_0110, 0b0001_0111, 0b0000_0110, 0b0001_0111,
		0b0001_0110, 0b0001_0111, 0b0010_0110, 0b0001_0111, 0b0011_0110,
		0b0001_0111, 0b0100_0110, 0b0001_0111, 0b0101_0110, 0b0001_0111,
		0b0110_0110, 0b0001_0111, 0b0111_0110, 0b0001_0111, 0b1000_0110,
		0b0001_0111, 0b1001_0110, 0b0001_0111, 0b1010_0110, 0b0001_0011,
		0b0000_0110, 0b0001_0011, 0b0001_0110, 0b0001_0011, 0b0010_0110,
		0b0001_0011, 0b0011_0110, 0b0001_0011, 0b0100_0000, 0b1011_0000,
		0b1001_1010, 0b1000_0000,
	})

	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	}

	arr, ok := v.([]any)
	if !ok {
		t.Errorf("Failed to decode the root array")
	} else if len(arr) == 0 {
		t.Errorf("Should not be an empty array")
	}

	if !slices.Equal(arr, expected) {
		t.Errorf("Expected:\n%v\nReceived:\n%v", expected, arr)
	}
}

func TestDecodeJsonInts(t *testing.T) {
	expected := []any{
		int64(10),
		int64(32021),
		int64(1503238552),
		int64(5000000000000000000),
		int64(-10),
	}

	v, err := Decode([]byte{
		0b0001_0100, 0b0101_1000, 0b0000_0101, 0b0100_0101, 0b1111_0100,
		0b0101_0110, 0b0100_1011, 0b0011_0011, 0b0011_0011, 0b0011_0011,
		0b0001_0011, 0b0100_0101, 0b0110_0011, 0b1001_0001, 0b1000_0010,
		0b0100_0100, 0b1111_0100, 0b0000_0000, 0b0000_0000, 0b1000_0111,
		0b1011_0000,
	})

	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	}

	arr, ok := v.([]any)
	if !ok {
		t.Errorf("Failed to decode the root array")
	} else if len(arr) == 0 {
		t.Errorf("Should not be an empty array")
	}

	if !slices.Equal(arr, expected) {
		t.Errorf("Expected:\n%v\nReceived:\n%v", expected, arr)
	}
}

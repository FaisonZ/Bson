package bson

import (
	"bytes"
	"testing"

	"github.com/FaisonZ/bson/internal/bit"
)

func TestEncodeJson(t *testing.T) {
	expected := []byte{
		0b00010100,
		0b01010010,
		0b00010110,
		0b00110110,
		0b01100110,
		0b11110110,
		0b11110110,
		0b00110110,
		0b00100110,
		0b00010111,
		0b00100110,
		0b01100110,
		0b01100110,
		0b11110110,
		0b11110110,
		0b00100110,
		0b00010111,
		0b00101011,
		0b10101100,
	}

	jsonBlob := []byte(`[{"foo":"bar"}, "foobar", true, false, null]`)
	bb := bit.NewBitBuilder()
	EncodeJson(jsonBlob, bb)

	if res := bytes.Compare(expected, bb.Bytes); res != 0 {
		t.Errorf("Expected: %08b\n Received: %08b", expected, bb.Bytes)
	}
}

func TestEncodeBoolean(t *testing.T) {
	expectedFalse := []byte{
		0b1010_0000,
	}
	expectedTrue := []byte{
		0b1011_0000,
	}

	bb := bit.NewBitBuilder()
	got := encodeBoolean(false, bb)

	if res := bytes.Compare(expectedFalse, bb.Bytes); res != 0 {
		t.Errorf(
			"False test expected:\n%08b\nReceived:\n%08b",
			expectedFalse,
			got,
		)
	}

	bb = bit.NewBitBuilder()
	got = encodeBoolean(true, bb)

	if res := bytes.Compare(expectedTrue, bb.Bytes); res != 0 {
		t.Errorf(
			"True test expected:\n%08b\nReceived:\n%08b",
			expectedFalse,
			got,
		)
	}
}

func TestEncodeLongString(t *testing.T) {
	expected := []byte{
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

	bb := bit.NewBitBuilder()
	encodeString("abcdefghijklmnopqrstuvwxyz012345", bb)

	if res := bytes.Compare(expected, bb.Bytes); res != 0 {
		t.Errorf("Expected:\n%08b\nReceived:\n%08b", expected, bb.Bytes)
	}
}

func TestEncodeLongString2(t *testing.T) {
	expected := []byte{
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
		0b0000_0000,
	}

	bb := bit.NewBitBuilder()
	encodeString("abcdefghijklmnopqrstuvwxyz01234", bb)

	if res := bytes.Compare(expected, bb.Bytes); res != 0 {
		t.Errorf("Expected:\n%08b\nReceived:\n%08b", expected, bb.Bytes)
	}
}

func TestEncodeLongObject(t *testing.T) {
	expected := []byte{
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
	}
	inObject := `{
    "0":null,
    "1":null,
    "2":null,
    "3":null,
    "4":null,
    "5":null,
    "a":null,
    "b":null,
    "c":null,
    "d":null,
    "e":null,
    "f":null,
    "g":null,
    "h":null,
    "i":null,
    "j":null,
    "k":null,
    "l":null,
    "m":null,
    "n":null,
    "o":null,
    "p":null,
    "q":null,
    "r":null,
    "s":null,
    "t":null,
    "u":null,
    "v":null,
    "w":null,
    "x":null,
    "y":null,
    "z":null
}`

	bb := bit.NewBitBuilder()
	EncodeJson([]byte(inObject), bb)

	if res := bytes.Compare(expected, bb.Bytes); res != 0 {
		t.Errorf("Expected:\n%08b\nReceived:\n%08b", expected, bb.Bytes)
	}
}

func TestEncodeLongObject2(t *testing.T) {
	expected := []byte{
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
		0b0000_0000,
	}
	inObject := `{
    "0":null,
    "1":null,
    "2":null,
    "3":null,
    "4":null,
    "5":null,
    "a":null,
    "b":null,
    "c":null,
    "d":null,
    "e":null,
    "f":null,
    "g":null,
    "h":null,
    "i":null,
    "j":null,
    "k":null,
    "l":null,
    "m":null,
    "n":null,
    "o":null,
    "p":null,
    "q":null,
    "r":null,
    "s":null,
    "t":null,
    "u":null,
    "v":null,
    "w":null,
    "x":null,
    "y":null
}`

	bb := bit.NewBitBuilder()
	EncodeJson([]byte(inObject), bb)

	if res := bytes.Compare(expected, bb.Bytes); res != 0 {
		t.Errorf("Expected:\n%08b\nReceived:\n%08b", expected, bb.Bytes)
	}
}

func TestEncodeJsonInts(t *testing.T) {
	expected := []byte{
		0b0001_0100,
		0b0101_1000,
		0b0000_0101,
		0b0100_0101,
		0b1111_0100,
		0b0101_0110,
		0b0100_1011,
		0b0011_0011,
		0b0011_0011,
		0b0011_0011,
		0b0001_0011,
		0b0100_0101,
		0b0110_0011,
		0b1001_0001,
		0b1000_0010,
		0b0100_0100,
		0b1111_0100,
		0b0000_0000,
		0b0000_0000,
		0b1000_0111,
		0b1011_0000,
	}

	jsonBlob := []byte(`[10,32021,1503238552,5000000000000000000,-10]`)
	bb := bit.NewBitBuilder()
	EncodeJson(jsonBlob, bb)

	if res := bytes.Compare(expected, bb.Bytes); res != 0 {
		t.Errorf("Expected:\n%08b\n Received:\n%08b", expected, bb.Bytes)
	}
}

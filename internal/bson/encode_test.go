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
	got := EncodeJson(jsonBlob, bb)

	if res := bytes.Compare(expected, bb.Bytes); res != 0 {
		t.Errorf("Expected: %08b\n Received: %08b", expected, got)
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

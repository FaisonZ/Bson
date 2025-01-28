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

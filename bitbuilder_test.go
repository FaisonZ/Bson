package main

import (
	"bytes"
	"testing"
)

func TestBitBuilder(t *testing.T) {
	expected := []byte{
		0b0001_0010,
		0b0001_0110,
		0b0011_0110,
		0b0110_0110,
		0b1111_0110,
		0b1111_0110,
		0b0011_0110,
		0b0010_0110,
		0b0001_0111,
		0b0010_0000,
	}

	b := NewBitBuilder()
	b.AddBits(0b0001, 4)
	b.AddBits(0b001, 3)
	b.AddBits(0b00001, 5)
	b.AddBits(0b011, 3)
	b.AddBits(0b00011, 5)
	b.AddBytes([]byte("foo"))
	b.AddBits(0b011, 3)
	b.AddBits(0b00011, 5)
	b.AddBytes([]byte("bar"))

	if res := bytes.Compare(expected, b.bytes); res != 0 {
		t.Errorf("Expected: %08b\n Received: %08b", expected, b.bytes)
	}
}

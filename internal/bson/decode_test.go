package bson

import (
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

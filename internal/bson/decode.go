package bson

import (
	"fmt"

	"github.com/FaisonZ/bson/internal/bit"
)

type Decoder struct {
	br          *bit.BitReader
	bsonVersion int
}

func NewDecoder(bs []byte) (*Decoder, error) {
	d := &Decoder{
		br:          bit.NewBitReader(bs),
		bsonVersion: 0,
	}

	if len(bs) < 1 {
		return d, fmt.Errorf("Empty bytes received")
	}

	err := d.decodeVersion()

	if err != nil {
		return d, err
	} else if d.bsonVersion != 1 {
		return d, fmt.Errorf("Invalid version number: %d", d.bsonVersion)
	}

	return d, nil
}

func (d *Decoder) decodeVersion() error {
	versionByte, err := d.br.GetBits(4)
	if err != nil {
		return err
	}

	d.bsonVersion = int(versionByte)

	return nil
}

func Decode(bs []byte, res *any) error {
	_, err := NewDecoder(bs)

	if err != nil {
		return err
	}

	return nil
}

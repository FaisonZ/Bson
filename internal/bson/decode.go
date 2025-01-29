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

func (d *Decoder) decode() (any, error) {
	err := d.decodeVersion()

	if err != nil {
		return nil, err
	} else if d.bsonVersion != 1 {
		return nil, fmt.Errorf("Invalid version number: %d", d.bsonVersion)
	}

	res, err := d.decodeValue()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d *Decoder) decodeValue() (any, error) {
	tokenBits, err := d.br.GetBits(3)

	fmt.Printf("Token: %03b\n", tokenBits)

	if err != nil {
		return nil, err
	}

	if tokenBits == OBJECT_TOKEN {
		return d.decodeObject()
	} else if tokenBits == STRING_TOKEN {
		return d.decodeString()
	} else if tokenBits == NULL_TOKEN {
		return nil, nil
	}

	return nil, fmt.Errorf("Invalid value token: %03b", tokenBits)
}

func (d *Decoder) decodeLength() (int, error) {
	l, err := d.br.GetBits(5)
	return int(l), err
}

func (d *Decoder) decodeObject() (any, error) {
	l, err := d.decodeLength()

	if err != nil {
		return nil, err
	}

	fmt.Printf("Object len: %d\n", l)
	o := make(map[string]any, l)

	for i := 0; i < l; i++ {
		d.br.GetBits(3)
		key, _ := d.decodeString()
		fmt.Printf("Key: %s\n", key)
		o[key], err = d.decodeValue()

		if err != nil {
			return o, err
		}
	}

	return o, nil
}

func (d *Decoder) decodeString() (string, error) {
	l, err := d.decodeLength()
	if err != nil {
		return "", err
	}

	fmt.Printf("Get %d bytes for string\n", l)
	sbytes, err := d.br.GetBytes(l)
	if err != nil {
		return "", err
	}

	return string(sbytes), nil
}

func Decode(bs []byte) (any, error) {
	d, err := NewDecoder(bs)

	if err != nil {
		return nil, err
	}

	res, err := d.decode()
	if err != nil {
		return nil, err
	}

	return res, nil
}

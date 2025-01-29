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

	if err != nil {
		return nil, err
	}

	switch tokenBits {
	case OBJECT_TOKEN:
		return d.decodeObject()
	case ARRAY_TOKEN:
		return d.decodeArray()
	case STRING_TOKEN:
		return d.decodeString()
	case NULL_TOKEN:
		return nil, nil
	case BOOLEAN_TOKEN:
		return d.decodeBoolean()
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

	o := make(map[string]any, l)

	for i := 0; i < l; i++ {
		d.br.GetBits(3)
		key, _ := d.decodeString()
		o[key], err = d.decodeValue()

		if err != nil {
			return o, err
		}
	}

	return o, nil
}

func (d *Decoder) decodeArray() (any, error) {
	l, err := d.decodeLength()

	if err != nil {
		return nil, err
	}

	o := make([]any, 0, l)

	for i := 0; i < l; i++ {
		val, err := d.decodeValue()
		if err != nil {
			return o, err
		}
		o = append(o, val)
	}

	return o, nil
}

func (d *Decoder) decodeString() (string, error) {
	l, err := d.decodeLength()
	if err != nil {
		return "", err
	}

	sbytes, err := d.br.GetBytes(l)
	if err != nil {
		return "", err
	}

	return string(sbytes), nil
}

func (d *Decoder) decodeBoolean() (bool, error) {
	l, err := d.br.GetBits(1)
	if err != nil {
		return false, err
	}

	return l == TRUE, nil
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

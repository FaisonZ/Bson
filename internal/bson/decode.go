package bson

import (
	"encoding/binary"
	"fmt"
	"maps"
	"math"

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
	case INTEGER_TOKEN:
		return d.decodeInteger()
	case FLOAT_TOKEN:
		return d.decodeFloat()
	}

	return nil, fmt.Errorf(
		"Invalid value token: %03b at %s",
		tokenBits,
		d.br.Debug(3),
	)
}

func (d *Decoder) decodeLength() (int, error) {
	l, err := d.br.GetBits(5)
	return int(l), err
}

func (d *Decoder) decodeObject() (map[string]any, error) {
	l, err := d.decodeLength()
	if err != nil {
		return nil, err
	}

	o := make(map[string]any, l)

	for i := 0; i < l; i++ {
		key, _ := d.decodeString()
		o[key], err = d.decodeValue()

		if err != nil {
			return o, err
		}
	}

	if l == MAX_CHUNK_LENGTH {
		moreO, err := d.decodeObject()
		if err != nil {
			return nil, err
		}
		maps.Copy(o, moreO)
	}

	return o, nil
}

func (d *Decoder) decodeArray() ([]any, error) {
	l, err := d.decodeLength()
	if err != nil {
		return []any{}, err
	}

	arr := make([]any, 0, l)

	for i := 0; i < l; i++ {
		val, err := d.decodeValue()
		if err != nil {
			return arr, err
		}
		arr = append(arr, val)
	}

	if l == MAX_CHUNK_LENGTH {
		moreArr, err := d.decodeArray()
		if err != nil {
			return []any{}, err
		}
		arr = append(arr, moreArr...)
	}

	return arr, nil
}

func (d *Decoder) decodeString() (string, error) {
	l, err := d.decodeLength()
	if err != nil {
		return "", err
	}
	//fmt.Printf("Len: %d\n", l)

	sbytes, err := d.br.GetBytes(l)
	if err != nil {
		return "", err
	}

	if l == MAX_CHUNK_LENGTH {
		moreBytes, err := d.decodeString()
		if err != nil {
			return "", err
		}
		sbytes = append(sbytes, []byte(moreBytes)...)
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

func (d *Decoder) decodeFloat() (float64, error) {
	fBytes, err := d.br.GetBytes(8)
	if err != nil {
		return 0, err
	}

	f := math.Float64frombits(binary.BigEndian.Uint64(fBytes))
	return f, nil
}

func (d *Decoder) decodeInteger() (i int64, err error) {
	s, err := d.br.GetBits(2)
	if err != nil {
		return 0, err
	}

	var intBytes []byte
	bytesToGet := 0

	switch s {
	case INT8_TOKEN:
		bytesToGet = 1
	case INT16_TOKEN:
		bytesToGet = 2
	case INT32_TOKEN:
		bytesToGet = 4
	case INT64_TOKEN:
		bytesToGet = 8
	}

	intBytes, err = d.br.GetBytes(bytesToGet)
	if err != nil {
		return 0, err
	}

	switch s {
	case INT8_TOKEN:
		i = int64(int8(intBytes[0]))
	case INT16_TOKEN:
		i = int64(binary.BigEndian.Uint16(intBytes))
	case INT32_TOKEN:
		i = int64(binary.BigEndian.Uint32(intBytes))
	case INT64_TOKEN:
		i = int64(binary.BigEndian.Uint64(intBytes))
	}

	return i, nil
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

package bit

import (
	"fmt"
	"io"
)

type BitBuilder struct {
	Bytes       []byte
	currBytePos int
	currByte    int
}

func NewBitBuilder() *BitBuilder {
	bb := &BitBuilder{
		Bytes:       make([]byte, 0, 10),
		currBytePos: 0,
		currByte:    0,
	}

	bb.Bytes = append(bb.Bytes, 0b0000_0000)

	return bb
}

func (bb *BitBuilder) grow() {
	bb.Bytes = append(bb.Bytes, 0b0000_0000)
	bb.currByte += 1
	bb.currBytePos = 0
}

func (bb *BitBuilder) AddBits(bits byte, len int) error {
	if bb.currBytePos == 8 {
		bb.grow()
	}

	shift := 8 - bb.currBytePos - len

	if shift < 0 {
		shiftedBits := bits >> -shift
		bb.Bytes[bb.currByte] |= shiftedBits
		bb.grow()

		shiftedBits = bits << (8 + shift)
		bb.Bytes[bb.currByte] |= shiftedBits
		bb.currBytePos += -shift
	} else {
		shiftedBits := bits << shift
		bb.Bytes[bb.currByte] |= shiftedBits
		bb.currBytePos += len
	}

	return nil
}

func (bb *BitBuilder) AddBytes(bs []byte) error {
	for _, b := range bs {
		bb.AddBits(b, 8)
	}

	return nil
}

func (bb *BitBuilder) String() string {
	return fmt.Sprintf("%08b", bb.Bytes)
}

func (bb *BitBuilder) WriteTo(w io.Writer) (n int64, err error) {
	nn, err := w.Write(bb.Bytes)
	return int64(nn), err
}

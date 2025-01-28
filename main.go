package main

import (
	"fmt"
	"io"
)

type BitBuilder struct {
	bytes       []byte
	currBytePos int
	currByte    int
}

func NewBitBuilder() *BitBuilder {
	bb := &BitBuilder{
		bytes:       make([]byte, 0, 10),
		currBytePos: 0,
		currByte:    0,
	}

	bb.bytes = append(bb.bytes, 0b0000_0000)

	return bb
}

func (bb *BitBuilder) grow() {
	bb.bytes = append(bb.bytes, 0b0000_0000)
	bb.currByte += 1
	bb.currBytePos = 0
}

func (bb *BitBuilder) AddBits(bits byte, len int) error {
	if bb.currBytePos == 8 {
		bb.grow()
	}

	shift := 8 - bb.currBytePos - len
	fmt.Printf(
		"currPos: %d ; len: %d ; shift: %d\ninput:   %08b\n",
		bb.currBytePos,
		len,
		shift,
		bits,
	)

	if shift < 0 {
		shiftedBits := bits >> -shift
		fmt.Printf("shifted: %08b\n", shiftedBits)
		bb.bytes[bb.currByte] |= shiftedBits
		bb.grow()

		shiftedBits = bits << (8 + shift)
		fmt.Printf("shifted: %08b\n", shiftedBits)
		bb.bytes[bb.currByte] |= shiftedBits
		bb.currBytePos += -shift
	} else {
		shiftedBits := bits << shift
		fmt.Printf("shifted: %08b\n", shiftedBits)
		bb.bytes[bb.currByte] |= shiftedBits
		bb.currBytePos += len
	}

	fmt.Println(bb)
	return nil
}

func (bb *BitBuilder) AddBytes(bs []byte) error {
	for _, b := range bs {
		bb.AddBits(b, 8)
	}

	return nil
}

func (bb *BitBuilder) String() string {
	return fmt.Sprintf("%08b", bb.bytes)
}

func (bb *BitBuilder) WriteTo(w io.Writer) (n int64, err error) {
	nn, err := w.Write(bb.bytes)
	return int64(nn), err
}

func main() {
	fmt.Println("Hello world")
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
}

package bit

import (
	"fmt"
	"math"
)

type BitReader struct {
	bytes   []byte
	bytePos int
	bitPos  int
}

func NewBitReader(bs []byte) *BitReader {
	return &BitReader{
		bytes:   bs,
		bytePos: 0,
		bitPos:  0,
	}
}

func (br *BitReader) GetBits(l int) (byte, error) {
	if br.bitPos == 8 {
		br.bytePos += 1
		br.bitPos = 0
	}

	shift := 8 - br.bitPos - l

	fmt.Printf(
		"shift: %d ; bitPos: %d ; bytePos: %d ; l: %d\n",
		shift,
		br.bitPos,
		br.bytePos,
		l,
	)
	if shift < 0 {
		hb := byte(math.Pow(2, float64(8-br.bitPos))) - 1
		l := br.bytes[br.bytePos] & hb
		l <<= -shift
		br.bytePos += 1
		br.bitPos = 0
		hb = byte(math.Pow(2, float64(8))) - 1
		r := br.bytes[br.bytePos] & hb
		r >>= (8 + shift)
		br.bitPos += (8 + shift)

		return l | r, nil
	} else {
		hb := byte(math.Pow(2, float64(8-br.bitPos))) - 1
		fmt.Printf("Mask: (%d) (%d) %08b\n", 8-br.bitPos, hb, hb)
		b := br.bytes[br.bytePos] & hb
		fmt.Printf("%08b\n", b)
		b >>= shift
		fmt.Printf("%08b\n", b)
		br.bitPos += l
		return b, nil
	}
}

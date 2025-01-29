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
	if l < 1 || l > 8 {
		return 0, fmt.Errorf(
			"GetBits: Invalid len (%d). len must be between 1 and 8",
			l,
		)
	}

	if br.bitPos == 8 {
		br.bytePos += 1
		br.bitPos = 0
	}

	shift := 8 - br.bitPos - l

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
		b := br.bytes[br.bytePos] & hb
		b >>= shift
		br.bitPos += l
		return b, nil
	}
}

func (br *BitReader) GetBytes(l int) ([]byte, error) {
	b := make([]byte, 0, l)

	for i := 0; i < l; i++ {
		c, _ := br.GetBits(8)
		b = append(b, c)
	}

	return b, nil
}

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

	/*
		fmt.Printf("bytePos: %d ; bitPos: %d\n", br.bytePos, br.bitPos)
		if br.bytePos > 0 {
			fmt.Printf("%08b ", br.bytes[br.bytePos-1])
		}
		fmt.Printf("%08b", br.bytes[br.bytePos])
		if br.bytePos+1 < len(br.bytes) {
			fmt.Printf(" %08b", br.bytes[br.bytePos+1])
		}
		fmt.Printf("\n")
	*/
	if br.bitPos == 8 {
		br.bytePos += 1
		br.bitPos = 0
	}

	bitsLeftInByte := 8 - br.bitPos
	shift := bitsLeftInByte - l

	// fmt.Printf("Shift: %d\n", shift)
	if shift < 0 {
		bitsFromNextByte := l - bitsLeftInByte
		lb, _ := br.GetBits(bitsLeftInByte)
		rb, _ := br.GetBits(bitsFromNextByte)

		return (lb << bitsFromNextByte) | rb, nil
	}

	hb := byte(math.Pow(2, float64(8-br.bitPos))) - 1
	b := br.bytes[br.bytePos] & hb
	b >>= shift
	br.bitPos += l
	// fmt.Printf("new bytePos: %d ; bitPos: %d\n", br.bytePos, br.bitPos)
	return b, nil
}

func (br *BitReader) GetBytes(l int) ([]byte, error) {
	b := make([]byte, 0, l)

	for i := 0; i < l; i++ {
		c, _ := br.GetBits(8)
		b = append(b, c)
	}

	return b, nil
}

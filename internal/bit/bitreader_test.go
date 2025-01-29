package bit

import (
	"bytes"
	"testing"
)

func TestGetBits1(t *testing.T) {
	bs := []byte{0b0001_0011, 0b0101_0111, 0b1110_0111}

	br := NewBitReader(bs)

	got, err := br.GetBits(4)
	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	} else if got != 0b0001 {
		t.Errorf("Got wrong value:\n%08b\nexpected\n%08b", got, 0b0001)
	}

	got, err = br.GetBits(3)
	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	} else if got != 0b001 {
		t.Errorf("Got wrong value:\n%08b\nexpected\n%08b", got, 0b001)
	}

	got, err = br.GetBits(5)
	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	} else if got != 0b10101 {
		t.Errorf("Got wrong value:\n%08b\nexpected\n%08b", got, 0b10101)
	}

	got, err = br.GetBits(4)
	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	} else if got != 0b0111 {
		t.Errorf("Got wrong value:\n%08b\nexpected\n%08b", got, 0b0111)
	}

	got, err = br.GetBits(8)
	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	} else if got != 0b1110_0111 {
		t.Errorf("Got wrong value:\n%08b\nexpected\n%08b", got, 0b1110_0111)
	}
}

func TestGetBytes(t *testing.T) {
	bs := []byte{
		0b0001_0110,
		0b0011_1000,
		0b0001_0100,
		0b0010_0010,
		0b0100_0000,
	}

	expectedBytes := []byte{
		0b1000_0001,
		0b0100_0010,
		0b0010_0100,
	}

	br := NewBitReader(bs)

	got, err := br.GetBits(4)
	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	} else if got != 0b0001 {
		t.Errorf("Got wrong value:\n%08b\nexpected\n%08b", got, 0b0001)
	}

	got, err = br.GetBits(3)
	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	} else if got != 0b011 {
		t.Errorf("Got wrong value:\n%08b\nexpected\n%08b", got, 0b011)
	}

	got, err = br.GetBits(5)
	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	} else if got != 0b00011 {
		t.Errorf("Got wrong value:\n%08b\nexpected\n%08b", got, 0b00011)
	}

	bytesGot, err := br.GetBytes(3)
	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	} else if res := bytes.Compare(expectedBytes, bytesGot); res != 0 {
		t.Errorf("Got wrong value:\n%08b\nexpected\n%08b", bytesGot, expectedBytes)
	}
}

func TestGetBitsEdgeCases(t *testing.T) {
	// 0111_1100 1000_1100
	//       ^
	//           ^ < should be here, but end up at 6
	// get 6 0b01_1111
	// get 3 0b001
	// get 5 0b0_0011

	bs := []byte{
		0b0111_1100,
		0b1000_1100,
	}

	br := NewBitReader(bs)

	got, err := br.GetBits(6)
	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	} else if got != 0b01_1111 {
		t.Errorf("Got wrong value:\n%08b\nexpected\n%08b", got, 0b01_1111)
	}

	got, err = br.GetBits(3)
	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	} else if got != 0b001 {
		t.Errorf("Got wrong value:\n%08b\nexpected\n%08b", got, 0b001)
	}

	got, err = br.GetBits(5)
	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	} else if got != 0b0_0011 {
		t.Errorf("Got wrong value:\n%08b\nexpected\n%08b", got, 0b0_0011)
	}
}

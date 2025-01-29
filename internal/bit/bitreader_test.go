package bit

import "testing"

func TestGetBits(t *testing.T) {
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

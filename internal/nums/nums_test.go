package nums

import "testing"

func TestIsInt(t *testing.T) {
	isIntTests := []struct {
		name     string
		input    float64
		expected bool
	}{
		{
			name:     "Returns true for positive int",
			input:    1232123,
			expected: true,
		},
		{
			name:     "Returns true for negative int",
			input:    -1232123,
			expected: true,
		},
		{
			name:     "Returns false for float",
			input:    -1232.123,
			expected: false,
		},
		{
			name:     "Returns true for int with zero after dot",
			input:    -1232.0,
			expected: true,
		},
	}

	for _, tt := range isIntTests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsInt(tt.input)
			if got != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestMinIntSize(t *testing.T) {
	minSizeTests := []struct {
		name    string
		inInt   float64
		outSize int
	}{
		{
			name:    "Returns 8 for 10",
			inInt:   10,
			outSize: 8,
		},
		{
			name:    "Returns 8 for -10",
			inInt:   -10,
			outSize: 8,
		},
		{
			name:    "Returns 16 for 30,000",
			inInt:   30_000,
			outSize: 16,
		},
		{
			name:    "Returns 32 for 32 bit numbers",
			inInt:   -2_147_483_648,
			outSize: 32,
		},
		{
			name:    "Returns 64",
			inInt:   5_000_000_000_000_000_000,
			outSize: 64,
		},
	}

	for _, tt := range minSizeTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MinIntSize(tt.inInt)
			if err != nil {
				t.Errorf("Unexpected error: %q", err)
			} else if got != tt.outSize {
				t.Errorf("Expected %v, got %v", tt.outSize, got)
			}
		})
	}
}

func TestIntFitsInSize(t *testing.T) {
	intFitsTests := []struct {
		name     string
		inInt    float64
		inSize   int
		expected bool
	}{
		{
			name:     "Returns true for 8 bit int",
			inInt:    10,
			inSize:   8,
			expected: true,
		},
		{
			name:     "Returns true for negative 8 bit int",
			inInt:    -10,
			inSize:   8,
			expected: true,
		},
		{
			name:     "Returns true for 16 bit int",
			inInt:    30_000,
			inSize:   16,
			expected: true,
		},
		{
			name:     "Returns false for 16 bit int with size 8",
			inInt:    30_000,
			inSize:   8,
			expected: false,
		},
		{
			name:     "Returns true for 32 bit int",
			inInt:    -2_147_483_648,
			inSize:   32,
			expected: true,
		},
		{
			name:     "Returns true for 64 bit int",
			inInt:    5_000_000_000_000_000_000,
			inSize:   64,
			expected: true,
		},
	}

	for _, tt := range intFitsTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IntFitsInSize(tt.inInt, tt.inSize)
			if err != nil {
				t.Errorf("Unexpected error: %q", err)
			} else if got != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestIntFitsInSizeErrors(t *testing.T) {
	_, err := IntFitsInSize(1.0, 0)
	if err == nil {
		t.Error("Should return error for incorrect int size")
	}

	_, err = IntFitsInSize(0.1, 8)
	if err == nil {
		t.Error("Should return error for passing in a non int float")
	}
}

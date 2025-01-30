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

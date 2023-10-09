package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldReturnArithmeticMean(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected float64
	}{
		{
			name:     "empty input",
			input:    nil,
			expected: 0.0,
		},
		{
			name:     "one element",
			input:    []int{3},
			expected: 3.0,
		},
		{
			name:     "multiple elements",
			input:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			expected: 5.0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// when
			res := ArithmeticMean(test.input...)

			// then
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestShouldReturnStandardDeviation(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected float64
	}{
		{
			name:     "empty input",
			input:    nil,
			expected: 0.0,
		},
		{
			name:     "one element",
			input:    []int{3},
			expected: 0.0,
		},
		{
			name:     "multiple elements",
			input:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			expected: 2.581988897471611,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// when
			res := StandardDeviation(test.input...)

			// then
			assert.Equal(t, test.expected, res)
		})
	}
}

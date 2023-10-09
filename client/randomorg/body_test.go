package randomorg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldParseBytesToIntegers(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []int
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty string",
			input:    []byte(""),
			expected: nil,
		},
		{
			name:     "one line",
			input:    []byte("17"),
			expected: []int{17},
		},
		{
			name:     "multiple lines",
			input:    []byte("2\n53\n-31"),
			expected: []int{2, 53, -31},
		},
		{
			name:     "multiple lines with whitespaces",
			input:    []byte(" 41\n -7 \n	0\n\n"),
			expected: []int{41, -7, 0},
		},
	}
	for _, test := range tests {
		// given
		contentType := "text/plain; charset=utf-8"
		sut := NewBodyParser()

		// when
		integers, err := sut.ParseIntegers(test.input, contentType)

		// then
		assert.NoError(t, err)
		assert.Equal(t, test.expected, integers)
	}
}

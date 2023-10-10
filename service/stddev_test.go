package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldReturnNothingWhenEmptyInputIsGiven(t *testing.T) {
	// given
	sut := NewStdDevService()
	pipe := make(chan []int)
	close(pipe)

	// when
	resultPipe := sut.Calculate(pipe)

	// then
	results := read(resultPipe)
	assert.Empty(t, results)
}

func TestShouldReturnStandardDeviationResult(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]int
		expected []StdDevResult
	}{
		{
			name: "one set",
			input: [][]int{
				{3},
			},
			expected: []StdDevResult{
				{
					StdDev: 0,
					Data:   []int{3},
				},
				{
					StdDev: 0,
					Data:   []int{3},
				},
			},
		},
		{
			name: "multiple sets",
			input: [][]int{
				{1, 2, 3, 4, 5},
				{6, 7, 8, 9},
			},
			expected: []StdDevResult{
				{
					StdDev: 1.4142135623730951,
					Data:   []int{1, 2, 3, 4, 5},
				},
				{
					StdDev: 1.118033988749895,
					Data:   []int{6, 7, 8, 9},
				},
				{
					StdDev: 2.581988897471611,
					Data:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			sut := NewStdDevService()
			pipe := make(chan []int)
			go func() {
				defer close(pipe)
				for _, set := range test.input {
					pipe <- set
				}
			}()

			// when
			resultPipe := sut.Calculate(pipe)

			// then
			results := read(resultPipe)
			assert.ElementsMatch(t, test.expected, results)
		})
	}
}

func read(pipe <-chan StdDevResult) []StdDevResult {
	var res []StdDevResult
	for r := range pipe {
		res = append(res, r)
	}
	return res
}

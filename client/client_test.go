package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldReturnErrorWhenResponseStatusCodeIsNotOK(t *testing.T) {
	// given
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	req, _ := http.NewRequest(http.MethodGet, fakeServer.URL, nil)
	sut := New()

	// when
	payload, contentType, err := sut.Send(req)

	// then
	assert.ErrorIs(t, err, ErrResponse)
	assert.Zero(t, payload)
	assert.Zero(t, contentType)
}

func TestShouldReturnPayloadBytesAndContentTypeWhenNoError(t *testing.T) {
	// given
	expectedContentType := "application/json; charset=utf-8"
	content := "A B C"
	expectedBytes := []byte(fmt.Sprintf("\"%s\"\n", content))
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", expectedContentType)
		json.NewEncoder(w).Encode(content)
	}))
	req, _ := http.NewRequest(http.MethodGet, fakeServer.URL, nil)
	sut := New()

	// when
	payload, contentType, err := sut.Send(req)

	// then
	assert.NoError(t, err)
	assert.Equal(t, expectedBytes, payload)
	assert.Equal(t, expectedContentType, contentType)
}

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

		// when
		integers, err := ParseIntegers(test.input, contentType)

		// then
		assert.NoError(t, err)
		assert.Equal(t, test.expected, integers)
	}
}

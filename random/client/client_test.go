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

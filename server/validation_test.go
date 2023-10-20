package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/koenno/standard-deviation-service/server/mocks"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnBadRequestWhenRequiredQueryParamsAreNotValid(t *testing.T) {
	tests := []struct {
		name            string
		requests        string
		length          string
		expectedPayload string
	}{
		{
			name:            "missing requests",
			requests:        "",
			length:          "1",
			expectedPayload: "requests parameter must be an integer",
		},
		{
			name:            "negative requests",
			requests:        "-1",
			length:          "1",
			expectedPayload: "requests parameter must be a positive integer",
		},
		{
			name:            "zero requests",
			requests:        "0",
			length:          "1",
			expectedPayload: "requests parameter must be a positive integer",
		},
		{
			name:            "missing length",
			requests:        "1",
			length:          "",
			expectedPayload: "length parameter must be an integer",
		},
		{
			name:            "negative length",
			requests:        "1",
			length:          "-1",
			expectedPayload: "length parameter must be a positive integer",
		},
		{
			name:            "zero length",
			requests:        "1",
			length:          "0",
			expectedPayload: "length parameter must be a positive integer",
		},
		{
			name:            "missing both",
			requests:        "",
			length:          "",
			expectedPayload: "requests parameter must be an integer",
		},
		{
			name:            "no number in requests",
			requests:        "A",
			length:          "1",
			expectedPayload: "requests parameter must be an integer",
		},
		{
			name:            "no number in length",
			requests:        "1",
			length:          "A",
			expectedPayload: "length parameter must be an integer",
		},
		{
			name:            "no number in both",
			requests:        "A",
			length:          "B",
			expectedPayload: "requests parameter must be an integer",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			URL := fmt.Sprintf("/random/mean?requests=%s&length=%s", test.requests, test.length)
			req := httptest.NewRequest(http.MethodGet, URL, nil)
			w := httptest.NewRecorder()
			httpHandlerMock := mocks.NewHandler(t)
			sut := validationMiddleware(httpHandlerMock)

			// when
			sut.ServeHTTP(w, req)

			// then
			res := w.Result()
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedPayload, string(data))
			assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		})
	}
}

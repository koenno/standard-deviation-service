package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/koenno/standard-deviation-service/server/mocks"
	"github.com/koenno/standard-deviation-service/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShouldReturnBadRequestWhenRequiredQueryParamsAreNotValid(t *testing.T) {
	tests := []struct {
		name     string
		requests string
		length   string
	}{
		{
			name:     "missing requests",
			requests: "",
			length:   "1",
		},
		{
			name:     "missing length",
			requests: "1",
			length:   "",
		},
		{
			name:     "missing both",
			requests: "",
			length:   "",
		},
		{
			name:     "no number in requests",
			requests: "A",
			length:   "1",
		},
		{
			name:     "no number in length",
			requests: "1",
			length:   "A",
		},
		{
			name:     "no number in both",
			requests: "A",
			length:   "B",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			URL := fmt.Sprintf("/random/mean?requests=%s&length=%s", test.requests, test.length)
			req := httptest.NewRequest(http.MethodGet, URL, nil)
			w := httptest.NewRecorder()
			generatorMock := mocks.NewRandomIntegerGenerator(t)
			calculatorMock := mocks.NewStdDevCalculator(t)
			sut := NewRandomServer(generatorMock, calculatorMock)

			// when
			sut.Mean(w, req)

			// then
			res := w.Result()
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			assert.Empty(t, data)
			assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		})
	}
}

func TestShouldReturnInternalServiceErrorWhenGeneratorFails(t *testing.T) {
	// given
	requests, length := "1", "2"
	URL := fmt.Sprintf("/random/mean?requests=%s&length=%s", requests, length)
	req := httptest.NewRequest(http.MethodGet, URL, nil)
	w := httptest.NewRecorder()
	generatorMock := mocks.NewRandomIntegerGenerator(t)
	calculatorMock := mocks.NewStdDevCalculator(t)
	sut := NewRandomServer(generatorMock, calculatorMock)

	generatorMock.EXPECT().Integers(mock.Anything, mock.Anything).Return(nil, errors.New("failure")).Once()

	calcPipe := make(chan service.StdDevResult)
	close(calcPipe)
	calculatorMock.EXPECT().Calculate(mock.Anything).Return(calcPipe).Once()

	// when
	sut.Mean(w, req)

	// then
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Empty(t, data)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestShouldReturnStandardDeviationCalculations(t *testing.T) {
	// given
	requests, length := "2", "5"
	URL := fmt.Sprintf("/random/mean?requests=%s&length=%s", requests, length)
	req := httptest.NewRequest(http.MethodGet, URL, nil)
	w := httptest.NewRecorder()
	generatorMock := mocks.NewRandomIntegerGenerator(t)
	calculatorMock := mocks.NewStdDevCalculator(t)
	sut := NewRandomServer(generatorMock, calculatorMock)

	genRes := []int{0, 1, 2, 3, 4}
	generatorMock.EXPECT().Integers(mock.Anything, mock.Anything).Return(genRes, nil).Twice()

	expectedResult := []service.StdDevResult{
		{
			StdDev: 1,
			Data:   genRes,
		},
		{
			StdDev: 1,
			Data:   genRes,
		},
		{
			StdDev: 2,
			Data:   append(genRes, genRes...),
		},
	}
	calcPipe := make(chan service.StdDevResult)
	go func() {
		defer close(calcPipe)
		calcPipe <- expectedResult[0]
		calcPipe <- expectedResult[1]
		calcPipe <- expectedResult[2]
	}()
	calculatorMock.EXPECT().Calculate(mock.Anything).Run(func(input <-chan []int) {
		go func() {
			for _ = range input {
			}
		}()
	}).Return(calcPipe).Once()

	// when
	sut.Mean(w, req)

	// then
	res := w.Result()
	defer res.Body.Close()
	var stddevResult []service.StdDevResult
	err := json.NewDecoder(res.Body).Decode(&stddevResult)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedResult, stddevResult)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

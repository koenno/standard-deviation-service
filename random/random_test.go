package random

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/koenno/standard-deviation-service/random/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShouldReturnErrorWhenRequestFails(t *testing.T) {
	// given
	senderMock := mocks.NewRequestSender(t)
	parserMock := mocks.NewResponseParser(t)
	reqFactoryMock := mocks.NewRequestFactory(t)
	sut := NewRandom(senderMock, parserMock, reqFactoryMock)
	quantity := 4

	req, err := http.NewRequest(http.MethodGet, "some.domain.com", nil)
	reqFactoryMock.EXPECT().NewRequest(mock.Anything, mock.Anything).Return(req, err).Once()
	senderMock.EXPECT().Send(req).Return(nil, "", errors.New("failure")).Once()

	// when
	ints, err := sut.Integers(context.Background(), quantity)

	// then
	assert.ErrorIs(t, err, ErrGenerator)
	assert.Zero(t, ints)
	parserMock.AssertNotCalled(t, "ParseIntegers")
}

func TestShouldReturnErrorWhenParsingResponseFails(t *testing.T) {
	// given
	senderMock := mocks.NewRequestSender(t)
	parserMock := mocks.NewResponseParser(t)
	reqFactoryMock := mocks.NewRequestFactory(t)
	sut := NewRandom(senderMock, parserMock, reqFactoryMock)
	quantity := 4
	contentType := "text/plain"
	response := []byte("")

	req, err := http.NewRequest(http.MethodGet, "some.domain.com", nil)
	reqFactoryMock.EXPECT().NewRequest(mock.Anything, mock.Anything).Return(req, err).Once()
	senderMock.EXPECT().Send(mock.AnythingOfType("*http.Request")).Return(response, contentType, nil).Once()
	parserMock.EXPECT().ParseIntegers(mock.Anything, mock.Anything).Return(nil, errors.New("failure")).Once()

	// when
	ints, err := sut.Integers(context.Background(), quantity)

	// then
	assert.ErrorIs(t, err, ErrItems)
	assert.Zero(t, ints)
}

func TestShouldReturnGeneratedIntegers(t *testing.T) {
	// given
	senderMock := mocks.NewRequestSender(t)
	parserMock := mocks.NewResponseParser(t)
	reqFactoryMock := mocks.NewRequestFactory(t)
	sut := NewRandom(senderMock, parserMock, reqFactoryMock)
	quantity := 4
	contentType := "text/plain"
	response := []byte("")
	expectedInts := []int{1, 7, 4}

	req, err := http.NewRequest(http.MethodGet, "some.domain.com", nil)
	reqFactoryMock.EXPECT().NewRequest(mock.Anything, mock.Anything).Return(req, err).Once()
	senderMock.EXPECT().Send(mock.AnythingOfType("*http.Request")).Return(response, contentType, nil).Once()
	parserMock.EXPECT().ParseIntegers(response, contentType).Return(expectedInts, nil).Once()

	// when
	ints, err := sut.Integers(context.Background(), quantity)

	// then
	assert.NoError(t, err)
	assert.Equal(t, expectedInts, ints)
}

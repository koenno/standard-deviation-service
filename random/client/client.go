package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	ErrSendRequest = errors.New("failed to send request")
	ErrResponse    = errors.New("response failure")

	httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
)

// // go:generate mockery --name=Converter --case underscore --with-expecter

type Client struct {
}

func New() Client {
	return Client{}
}

func (c Client) Send(req *http.Request) ([]byte, string, error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %v", ErrSendRequest, err)
	}

	defer resp.Body.Close()
	payloadBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("%w: unable to read body: %v", ErrResponse, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("%w: status code %d; body %s", ErrResponse, resp.StatusCode, string(payloadBytes))
	}

	return payloadBytes, resp.Header.Get("content-type"), nil
}

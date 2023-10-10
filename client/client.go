package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/exp/slog"
)

var (
	ErrSendRequest = errors.New("failed to send request")
	ErrResponse    = errors.New("response failure")

	httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
)

//go:generate mockery --name=RateLimiter --case underscore --with-expecter
type RateLimiter interface {
	Wait(ctx context.Context) (err error)
}

type Client struct {
	rateLimiter RateLimiter
}

func New(rateLimiter RateLimiter) Client {
	return Client{
		rateLimiter: rateLimiter,
	}
}

func (c Client) Send(req *http.Request) ([]byte, string, error) {
	if c.rateLimiter != nil {
		if err := c.rateLimiter.Wait(req.Context()); err != nil {
			return nil, "", fmt.Errorf("failed to limit a rate: %v", err)
		}
	}

	slog.Info("client sends a request", "method", req.Method, "url", req.URL.String())

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

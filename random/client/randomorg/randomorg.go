package randomorg

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type options struct {
	min      int
	max      int
	quantity int
}

func newOptions() *options {
	return &options{
		min:      1,
		max:      10,
		quantity: 5,
	}
}

type Option func(*options)

func WithMin(min int) Option {
	return func(o *options) {
		o.min = min
	}
}

func WithMax(max int) Option {
	return func(o *options) {
		o.max = max
	}
}

func WithQuantity(quantity int) Option {
	return func(o *options) {
		o.quantity = quantity
	}
}

func NewRequest(ctx context.Context, opts ...Option) (*http.Request, error) {
	cfg := newOptions()
	for _, o := range opts {
		o(cfg)
	}

	rawURL := "https://www.random.org/integers/?num=5&min=1&max=100&col=1&base=10&format=plain&rnd=new"
	URL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %v", err)
	}

	query := URL.Query()
	query.Set("min", strconv.Itoa(cfg.min))
	query.Set("max", strconv.Itoa(cfg.max))
	query.Set("num", strconv.Itoa(cfg.quantity))
	URL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	return req, nil
}

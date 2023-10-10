package randomorg

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/koenno/standard-deviation-service/client"
)

type RequestFactory struct {
}

func NewRequestFactory() RequestFactory {
	return RequestFactory{}
}

func (f RequestFactory) NewRequest(ctx context.Context, opts ...client.Option) (*http.Request, error) {
	cfg := client.NewOptions(opts...)

	rawURL := "https://www.random.org/integers/?num=5&min=1&max=100&col=1&base=10&format=plain&rnd=new"
	URL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %v", err)
	}

	query := URL.Query()
	query.Set("min", strconv.Itoa(cfg.Min))
	query.Set("max", strconv.Itoa(cfg.Max))
	query.Set("num", strconv.Itoa(cfg.Quantity))
	URL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	return req, nil
}

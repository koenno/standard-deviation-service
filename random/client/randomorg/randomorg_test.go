package randomorg

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldReturnProperURLWithDefaultValues(t *testing.T) {
	// when
	req, err := NewRequest(context.Background())

	// then
	assert.NoError(t, err)
	assert.Equal(t, http.MethodGet, req.Method)
	assert.Equal(t, "www.random.org", req.URL.Host)
	assert.Equal(t, "https", req.URL.Scheme)
	assert.Equal(t, "/integers/", req.URL.Path)
	query, err := url.ParseQuery(req.URL.RawQuery)
	assert.NoError(t, err)
	assert.Equal(t, "5", query.Get("num"))
	assert.Equal(t, "1", query.Get("min"))
	assert.Equal(t, "10", query.Get("max"))
	assert.Equal(t, "1", query.Get("col"))
	assert.Equal(t, "10", query.Get("base"))
	assert.Equal(t, "plain", query.Get("format"))
	assert.Equal(t, "new", query.Get("rnd"))
}

func TestShouldReturnProperURLWithCustomValues(t *testing.T) {
	// given
	min := -11
	max := 435
	quantity := 23

	// when
	req, err := NewRequest(context.Background(), WithQuantity(quantity), WithMin(min), WithMax(max))

	// then
	assert.NoError(t, err)
	assert.Equal(t, http.MethodGet, req.Method)
	assert.Equal(t, "www.random.org", req.URL.Host)
	assert.Equal(t, "https", req.URL.Scheme)
	assert.Equal(t, "/integers/", req.URL.Path)
	query, err := url.ParseQuery(req.URL.RawQuery)
	assert.NoError(t, err)
	assert.Equal(t, "23", query.Get("num"))
	assert.Equal(t, "-11", query.Get("min"))
	assert.Equal(t, "435", query.Get("max"))
	assert.Equal(t, "1", query.Get("col"))
	assert.Equal(t, "10", query.Get("base"))
	assert.Equal(t, "plain", query.Get("format"))
	assert.Equal(t, "new", query.Get("rnd"))
}

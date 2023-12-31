package random

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/koenno/standard-deviation-service/client"
)

var (
	ErrInit      = errors.New("failed to initialize random generator")
	ErrGenerator = errors.New("random generator failure")
	ErrItems     = errors.New("failed to obtain random items")
)

//go:generate mockery --name=RequestFactory --case underscore --with-expecter
type RequestFactory interface {
	NewRequest(ctx context.Context, opts ...client.Option) (*http.Request, error)
}

//go:generate mockery --name=RequestSender --case underscore --with-expecter
type RequestSender interface {
	Send(req *http.Request) ([]byte, string, error)
}

//go:generate mockery --name=ResponseParser --case underscore --with-expecter
type ResponseParser interface {
	ParseIntegers(bb []byte, contentType string) ([]int, error)
}

type Random struct {
	reqFactory RequestFactory
	reqSender  RequestSender
	respParser ResponseParser
}

func NewRandom(reqSender RequestSender, respParser ResponseParser, reqFactory RequestFactory) Random {
	return Random{
		reqFactory: reqFactory,
		reqSender:  reqSender,
		respParser: respParser,
	}
}

func (r Random) Integers(ctx context.Context, quantity int) ([]int, error) {
	req, err := r.reqFactory.NewRequest(ctx, client.WithQuantity(quantity))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInit, err)
	}

	bb, contentType, err := r.reqSender.Send(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGenerator, err)
	}

	ints, err := r.respParser.ParseIntegers(bb, contentType)
	if err != nil {
		return nil, fmt.Errorf("%w (integers): %v", ErrItems, err)
	}

	return ints, nil
}

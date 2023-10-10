package client

type Options struct {
	Min      int
	Max      int
	Quantity int
}

func NewOptions(opts ...Option) *Options {
	cfg := defaultOptions()
	for _, o := range opts {
		o(cfg)
	}
	return cfg
}

func defaultOptions() *Options {
	return &Options{
		Min:      1,
		Max:      10,
		Quantity: 5,
	}
}

type Option func(*Options)

func WithMin(min int) Option {
	return func(o *Options) {
		o.Min = min
	}
}

func WithMax(max int) Option {
	return func(o *Options) {
		o.Max = max
	}
}

func WithQuantity(quantity int) Option {
	return func(o *Options) {
		o.Quantity = quantity
	}
}

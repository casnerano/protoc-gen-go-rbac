package interceptor

import "context"

type options struct {
	debug bool
}

type Option func(*options)

func WithDebug() Option {
	return func(options *options) {
		options.debug = true
	}
}

type Authenticator func(ctx context.Context) (bool, error)

func WithAuthentification() Option {
	return func(o *options) {

	}
}

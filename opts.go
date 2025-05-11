package zerr

import (
	"net/http"
)

type options struct {
	httpStatusCode int
	redacted       bool
}

type OptFunc = func(*options)

func resolveOptions(optFuncs []OptFunc) *options {
	opts := &options{
		httpStatusCode: http.StatusInternalServerError,
		redacted:       true,
	}
	for _, fn := range optFuncs {
		fn(opts)
	}
	return opts
}

func WithHTTPStatusCode(code int) OptFunc {
	return func(opts *options) {
		opts.httpStatusCode = code
	}
}

func WithRedacted(r bool) OptFunc {
	return func(opts *options) {
		opts.redacted = r
	}
}

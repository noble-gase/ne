package httpzip

import (
	"net/http"

	"github.com/go-resty/resty/v2"
)

type Option func(r *Reader)

func WithClient(c *http.Client) Option {
	return func(r *Reader) {
		r.client = resty.NewWithClient(c)
	}
}

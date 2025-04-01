package cors

import (
	"net/http"
	"slices"
	"strconv"
	"strings"
)

const (
	wildcard   = "*"
	HeaderVary = "Vary"
	HeaderACAO = "Access-Control-Allow-Origin"
	HeaderACAM = "Access-Control-Allow-Methods"
	HeaderACAH = "Access-Control-Allow-Headers"
	HeaderACEH = "Access-Control-Expose-Headers"
	HeaderACAC = "Access-Control-Allow-Credentials"
	HeaderACMA = "Access-Control-Max-Age"
)

type Cors struct {
	allowOrigins     []string
	allowMethods     []string
	allowHeaders     []string
	exposeHeaders    []string
	allowCredentials bool
	maxAge           int
}

func (c *Cors) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if len(origin) != 0 {
			w.Header().Set(HeaderVary, "Origin")
		}

		// Access-Control-Allow-Origin
		if slices.Contains(c.allowOrigins, wildcard) {
			w.Header().Set(HeaderACAO, wildcard)
		} else {
			if slices.Contains(c.allowOrigins, origin) {
				w.Header().Set(HeaderACAO, origin)
			}
			// Access-Control-Allow-Credentials
			if c.allowCredentials {
				w.Header().Set(HeaderACAC, "true")
			}
		}

		// Access-Control-Allow-Methods
		w.Header().Set(HeaderACAM, strings.Join(c.allowMethods, ", "))

		// Access-Control-Allow-Headers
		if slices.Contains(c.allowHeaders, wildcard) {
			w.Header().Set(HeaderACAH, wildcard)
		} else {
			w.Header().Set(HeaderACAH, strings.Join(c.allowHeaders, ", "))
		}

		// Access-Control-Expose-Headers
		if len(c.exposeHeaders) != 0 {
			w.Header().Set(HeaderACEH, strings.Join(c.exposeHeaders, ", "))
		}

		// Access-Control-Max-Age
		if c.maxAge > 0 {
			w.Header().Set(HeaderACMA, strconv.Itoa(c.maxAge))
		}

		// Preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})
}

// New 创建一个 CORS 中间件，默认允许所有跨域请求
func New(opts ...Option) *Cors {
	c := &Cors{
		allowOrigins: []string{wildcard},
		allowMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		allowHeaders: []string{wildcard},
	}
	for _, f := range opts {
		f(c)
	}
	return c
}

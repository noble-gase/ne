package protokit

import (
	"net/http"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
)

// ClientOption option for resty.Client
type ClientOption func(c *resty.Client)

// WithDebug enables the debug mode on the Resty client. The client logs details
// of every request and response.
//
//	WithDebug(true)
//
// Also, it can be enabled at the request level for a particular request; see [WithReqDebug].
//   - For [Request], it logs information such as HTTP verb, Relative URL path,
//     Host, Headers, and Body if it has one.
//   - For [Response], it logs information such as Status, Response Time, Headers,
//     and Body if it has one.
func WithDebug(b bool) ClientOption {
	return func(c *resty.Client) {
		c.SetDebug(b)
	}
}

// WithLogger sets given writer for logging Resty request and response details.
//
// Compliant to interface [resty.Logger]
func WithLogger(l resty.Logger) ClientOption {
	return func(c *resty.Client) {
		c.SetLogger(l)
	}
}

// WithRetryCount enables retry on Resty client and allows you
// to set no. of retry count. Resty uses a Backoff mechanism.
func WithRetryCount(n int) ClientOption {
	return func(c *resty.Client) {
		c.SetRetryCount(n)
	}
}

// WithRetryWaitTime sets the default wait time for sleep before retrying
// request.
//
// Default is 100 milliseconds.
func WithRetryWaitTime(d time.Duration) ClientOption {
	return func(c *resty.Client) {
		c.SetRetryWaitTime(d)
	}
}

// WithRetryMaxWaitTime sets the max wait time for sleep before retrying
// request.
//
// Default is 2 seconds.
func WithRetryMaxWaitTime(d time.Duration) ClientOption {
	return func(c *resty.Client) {
		c.SetRetryMaxWaitTime(d)
	}
}

// WithRetryAfter sets a callback to calculate the wait time between retries.
// Default (nil) implies exponential backoff with jitter
func WithRetryAfter(fn resty.RetryAfterFunc) ClientOption {
	return func(c *resty.Client) {
		c.SetRetryAfter(fn)
	}
}

// RequestOption option for resty.Request
type RequestOption func(req *resty.Request)

// WithReqHeader sets multiple header fields and their values as a list of strings in the current request.
//
// For Example: To set `Accept` as `text/html, application/xhtml+xml, application/xml;q=0.9, image/webp, */*;q=0.8`
//
//	WithReqHeader(http.Header{
//		"Accept": []string{"text/html", "application/xhtml+xml", "application/xml;q=0.9", "image/webp", "*/*;q=0.8"},
//	})
//
// It overrides the header value set at the client instance level.
func WithReqHeader(h http.Header) RequestOption {
	return func(req *resty.Request) {
		req.SetHeaderMultiValues(h)
	}
}

// WithReqQuery appends multiple parameters with multi-value
// ([url.Values]) at one go in the current request. It will be formed as
// query string for the request.
//
// For Example: `status=pending&status=approved&status=open` in the URL after the `?` mark.
//
//	WithReqQuery(url.Values{
//		"status": []string{"pending", "approved", "open"},
//	})
//
// It overrides the query parameter value set at the client instance level.
func WithReqQuery(v url.Values) RequestOption {
	return func(req *resty.Request) {
		req.SetQueryParamsFromValues(v)
	}
}

// WithReqDebug enables the debug mode on the current request. It logs
// the details current request and response.
//
//	WithReqDebug(true)
//
// Also, it can be enabled at the request level for a particular request; see [Request.SetDebug].
//   - For [Request], it logs information such as HTTP verb, Relative URL path,
//     Host, Headers, and Body if it has one.
//   - For [Response], it logs information such as Status, Response Time, Headers,
//     and Body if it has one.
func WithReqDebug(b bool) RequestOption {
	return func(req *resty.Request) {
		req.SetDebug(b)
	}
}

// WithReqLogger sets given writer for logging Resty request and response details.
// By default, requests and responses inherit their logger from the client.
//
// Compliant to interface [resty.Logger].
//
// It overrides the logger value set at the client instance level.
func WithReqLogger(l resty.Logger) RequestOption {
	return func(req *resty.Request) {
		req.SetLogger(l)
	}
}

// WithReqRetryCondition adds a retry condition function to the request's
// array of functions is checked to determine if the request can be retried.
// The request will retry if any functions return true and the error is nil.
//
// NOTE: The request level retry conditions are checked before all retry
// conditions from the client instance.
func WithReqRetryCondition(fn resty.RetryConditionFunc) RequestOption {
	return func(req *resty.Request) {
		req.AddRetryCondition(fn)
	}
}

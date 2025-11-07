package helper

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const XTraceId = "x-trace-id"

const (
	HeaderAccept        = "Accept"
	HeaderAuthorization = "Authorization"
	HeaderContentType   = "Content-Type"
)

const (
	ContentText          = "text/plain; charset=utf-8"
	ContentJSON          = "application/json"
	ContentXML           = "application/xml"
	ContentForm          = "application/x-www-form-urlencoded"
	ContentStream        = "application/octet-stream"
	ContentMultipartForm = "multipart/form-data"
)

func ContentType(h http.Header) string {
	content := h.Get(HeaderContentType)
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

// RestyClient default client for http request
var RestyClient = resty.NewWithClient(NewHttpClient())

// NewHttpClient returns a http client
func NewHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			MaxIdleConns:          0,
			MaxIdleConnsPerHost:   1000,
			MaxConnsPerHost:       1000,
			IdleConnTimeout:       60 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: time.Second,
		},
	}
}

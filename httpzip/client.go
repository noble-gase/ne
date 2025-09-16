package httpzip

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	once   sync.Once
	client *resty.Client
)

func Client() *resty.Client {
	if client == nil {
		once.Do(func() {
			client = resty.NewWithClient(&http.Client{
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
			})
		})
	}
	return client
}

func SetClient(c *http.Client) {
	client = resty.NewWithClient(c)
}

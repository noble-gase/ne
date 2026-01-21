package helper

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strings"
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

type ApiResult[T any] struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func (r *ApiResult[T]) Error(ok int) error {
	if r == nil || r.Code == ok {
		return nil
	}

	msgs := make([]string, 0, 2)
	if len(r.Msg) != 0 {
		msgs = append(msgs, r.Msg)
	}
	if len(r.Message) != 0 {
		msgs = append(msgs, r.Message)
	}
	return fmt.Errorf("[%d] %s", r.Code, strings.Join(msgs, "; "))
}

func HttpGet(ctx context.Context, url string, query url.Values, header ...http.Header) (resp *resty.Response, err error) {
	start := time.Now()
	defer func() {
		cost := time.Since(start)

		level := slog.LevelInfo

		attrs := []slog.Attr{
			slog.String("url", url),
			slog.String("method", http.MethodGet),
			slog.Any("query", query),
			slog.String("duration", cost.String()),
		}

		if err != nil {
			level = slog.LevelError
			attrs = append(attrs, slog.Any("err", err))
		}
		if resp != nil {
			attrs = append(attrs, slog.String("resp", resp.String()), slog.Int("status", resp.StatusCode()))
		}

		slog.LogAttrs(ctx, level, "http request", attrs...)
	}()

	req := RestyClient.R().
		SetContext(ctx).
		SetQueryParamsFromValues(query)
	if len(header) != 0 {
		req.SetHeaderMultiValues(header[0])
	}

	resp, err = req.Get(url)
	return
}

func HttpGetX[T any](ctx context.Context, url string, query url.Values, header ...http.Header) (ret *ApiResult[T], err error) {
	var resp *resty.Response

	start := time.Now()
	defer func() {
		cost := time.Since(start)

		level := slog.LevelInfo

		attrs := []slog.Attr{
			slog.String("url", url),
			slog.String("method", http.MethodGet),
			slog.Any("query", query),
			slog.String("duration", cost.String()),
		}

		if err != nil {
			level = slog.LevelError
			attrs = append(attrs, slog.Any("err", err))
		}
		if resp != nil {
			attrs = append(attrs, slog.String("resp", resp.String()), slog.Int("status", resp.StatusCode()))
		}

		slog.LogAttrs(ctx, level, "http request", attrs...)
	}()

	ret = new(ApiResult[T])

	req := RestyClient.R().
		SetContext(ctx).
		SetQueryParamsFromValues(query).
		SetResult(ret)
	if len(header) != 0 {
		req.SetHeaderMultiValues(header[0])
	}

	resp, err = req.Get(url)
	if err == nil && resp.StatusCode() != http.StatusOK {
		err = errors.New(resp.Status())
	}
	return
}

func HttpPost(ctx context.Context, url string, body any, header ...http.Header) (resp *resty.Response, err error) {
	start := time.Now()
	defer func() {
		cost := time.Since(start)

		level := slog.LevelInfo

		attrs := []slog.Attr{
			slog.String("url", url),
			slog.String("method", http.MethodPost),
			slog.Any("body", body),
			slog.String("duration", cost.String()),
		}

		if err != nil {
			level = slog.LevelError
			attrs = append(attrs, slog.Any("err", err))
		}
		if resp != nil {
			attrs = append(attrs, slog.String("resp", resp.String()), slog.Int("status", resp.StatusCode()))
		}

		slog.LogAttrs(ctx, level, "http request", attrs...)
	}()

	req := RestyClient.R().
		SetContext(ctx).
		SetBody(body)
	if len(header) != 0 {
		req.SetHeaderMultiValues(header[0])
	}

	resp, err = req.Post(url)
	return
}

func HttpPostX[T any](ctx context.Context, url string, body any, header ...http.Header) (ret *ApiResult[T], err error) {
	var resp *resty.Response

	start := time.Now()
	defer func() {
		cost := time.Since(start)

		level := slog.LevelInfo

		attrs := []slog.Attr{
			slog.String("url", url),
			slog.String("method", http.MethodPost),
			slog.Any("body", body),
			slog.String("duration", cost.String()),
		}

		if err != nil {
			level = slog.LevelError
			attrs = append(attrs, slog.Any("err", err))
		}
		if resp != nil {
			attrs = append(attrs, slog.String("resp", resp.String()), slog.Int("status", resp.StatusCode()))
		}

		slog.LogAttrs(ctx, level, "http request", attrs...)
	}()

	ret = new(ApiResult[T])

	req := RestyClient.R().
		SetContext(ctx).
		SetBody(body).
		SetResult(ret)
	if len(header) != 0 {
		req.SetHeaderMultiValues(header[0])
	}

	resp, err = req.Post(url)
	if err == nil && resp.StatusCode() != http.StatusOK {
		err = errors.New(resp.Status())
	}
	return
}

func HttpForm(ctx context.Context, url string, form url.Values, header ...http.Header) (resp *resty.Response, err error) {
	start := time.Now()
	defer func() {
		cost := time.Since(start)

		level := slog.LevelInfo

		attrs := []slog.Attr{
			slog.String("url", url),
			slog.String("method", http.MethodPost),
			slog.Any("form", form),
			slog.String("duration", cost.String()),
		}

		if err != nil {
			level = slog.LevelError
			attrs = append(attrs, slog.Any("err", err))
		}
		if resp != nil {
			attrs = append(attrs, slog.String("resp", resp.String()), slog.Int("status", resp.StatusCode()))
		}

		slog.LogAttrs(ctx, level, "http request", attrs...)
	}()

	req := RestyClient.R().
		SetContext(ctx).
		SetFormDataFromValues(form)
	if len(header) != 0 {
		req.SetHeaderMultiValues(header[0])
	}

	resp, err = req.Post(url)
	return
}

func HttpFormX[T any](ctx context.Context, url string, form url.Values, header ...http.Header) (ret *ApiResult[T], err error) {
	var resp *resty.Response

	start := time.Now()
	defer func() {
		cost := time.Since(start)

		level := slog.LevelInfo

		attrs := []slog.Attr{
			slog.String("url", url),
			slog.String("method", http.MethodPost),
			slog.Any("form", form),
			slog.String("duration", cost.String()),
		}

		if err != nil {
			level = slog.LevelError
			attrs = append(attrs, slog.Any("err", err))
		}
		if resp != nil {
			attrs = append(attrs, slog.String("resp", resp.String()), slog.Int("status", resp.StatusCode()))
		}

		slog.LogAttrs(ctx, level, "http request", attrs...)
	}()

	ret = new(ApiResult[T])

	req := RestyClient.R().
		SetContext(ctx).
		SetFormDataFromValues(form).
		SetResult(ret)
	if len(header) != 0 {
		req.SetHeaderMultiValues(header[0])
	}

	resp, err = req.Get(url)
	if err == nil && resp.StatusCode() != http.StatusOK {
		err = errors.New(resp.Status())
	}
	return
}

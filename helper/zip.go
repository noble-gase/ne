package helper

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// HTTPReaderAt 用 HTTP Range 请求实现 io.ReaderAt
type HTTPReaderAt struct {
	ctx  context.Context
	url  string
	size int64
}

func NewHTTPReaderAt(ctx context.Context, url string) (*HTTPReaderAt, error) {
	resp, err := RestyClient.R().SetContext(ctx).Head(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status())
	}

	cl := resp.Header().Get("Content-Length")
	if len(cl) == 0 {
		return nil, fmt.Errorf("missing Content-Length header")
	}

	size, err := strconv.ParseInt(cl, 10, 64)
	if err != nil {
		return nil, err
	}

	return &HTTPReaderAt{
		ctx:  ctx,
		url:  url,
		size: size,
	}, nil
}

// 实现 io.ReaderAt
func (r *HTTPReaderAt) ReadAt(p []byte, off int64) (n int, err error) {
	end := off + int64(len(p)) - 1

	resp, err := RestyClient.R().
		SetContext(r.ctx).
		SetHeader("Range", fmt.Sprintf("bytes=%d-%d", off, end)).
		SetDoNotParseResponse(true).
		Get(r.url)
	if err != nil {
		return 0, err
	}
	defer resp.RawResponse.Body.Close()

	if resp.StatusCode() != http.StatusPartialContent && resp.StatusCode() != http.StatusOK {
		return 0, fmt.Errorf("bad status: %s", resp.Status())
	}
	return io.ReadFull(resp.RawResponse.Body, p)
}

// Size 返回文件大小
func (r *HTTPReaderAt) Size() int64 {
	return r.size
}

func ReadFileFromZipURL(ctx context.Context, url, filename string) ([]byte, error) {
	hr, err := NewHTTPReaderAt(ctx, url)
	if err != nil {
		return nil, err
	}

	// 交给标准库 zip.NewReader 处理 ZIP64
	zr, err := zip.NewReader(hr, hr.Size())
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	// 读取指定文件
	for _, f := range zr.File {
		if f.Name == filename {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			if _, err = io.Copy(buf, rc); err != nil {
				return nil, err
			}
			break
		}
	}
	return buf.Bytes(), nil
}

func ReadFileFromZipURLToWriter(ctx context.Context, url, filename string, w io.Writer) error {
	hr, err := NewHTTPReaderAt(ctx, url)
	if err != nil {
		return err
	}

	// 交给标准库 zip.NewReader 处理 ZIP64
	zr, err := zip.NewReader(hr, hr.Size())
	if err != nil {
		return err
	}

	// 读取指定文件
	for _, f := range zr.File {
		if f.Name == filename {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			if _, err = io.Copy(w, rc); err != nil {
				return err
			}
			break
		}
	}
	return nil
}

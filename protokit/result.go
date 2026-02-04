package protokit

import (
	"encoding/base64"
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"
)

// ApiResult is the API response with proto.Message
type ApiResult[T proto.Message] struct {
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

// Bytes converts the given string representation of a byte sequence into a slice of bytes
// A bytes sequence is encoded in URL-safe base64 without padding
func Bytes(s string) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		b, err = base64.URLEncoding.DecodeString(s)
		if err != nil {
			return nil, err
		}
	}
	return b, nil
}

package helper

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/noble-gase/ne/protokit"
	"github.com/noble-gase/ne/validkit"
	"google.golang.org/protobuf/proto"
)

const MaxFormMemory = 32 << 20

// BindJSON 解析JSON请求体并校验
func BindJSON(r *http.Request, obj any) error {
	if r.Body != nil && r.Body != http.NoBody {
		defer io.Copy(io.Discard, r.Body)
		if err := json.NewDecoder(r.Body).Decode(obj); err != nil {
			return err
		}
	}
	return validkit.ValidateStruct(obj)
}

// BindProto 解析Proto请求体并校验
func BindProto(r *http.Request, msg proto.Message) error {
	// GET请求
	if r.Method == http.MethodGet {
		if err := protokit.ValuesToMessage(msg, r.URL.Query()); err != nil {
			return err
		}
		return protokit.Validate(msg)
	}

	// 根据Content-Type解析请求体
	switch ContentType(r.Header) {
	case ContentForm:
		if err := r.ParseForm(); err != nil {
			return err
		}
		if err := protokit.ValuesToMessage(msg, r.PostForm); err != nil {
			return err
		}
	case ContentMultipartForm:
		if err := r.ParseMultipartForm(MaxFormMemory); err != nil {
			if err != http.ErrNotMultipart {
				return err
			}
		}
		if err := protokit.ValuesToMessage(msg, r.PostForm); err != nil {
			return err
		}
	case ContentJSON:
		if r.Body != nil && r.Body != http.NoBody {
			defer io.Copy(io.Discard, r.Body)
			if err := json.NewDecoder(r.Body).Decode(msg); err != nil {
				return err
			}
		}
	default:
		return errors.New("unsupported Content-Type")
	}
	return protokit.Validate(msg)
}

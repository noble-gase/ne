package ne

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/noble-gase/ne/protos"
	"github.com/noble-gase/ne/validator"
	"google.golang.org/protobuf/proto"
)

// BindJSON 解析JSON请求体并校验
func BindJSON(r *http.Request, obj any) error {
	if r.Body != nil && r.Body != http.NoBody {
		defer io.Copy(io.Discard, r.Body)
		if err := json.NewDecoder(r.Body).Decode(obj); err != nil {
			return err
		}
	}
	return validator.ValidateStruct(obj)
}

// BindForm 解析Form表单并校验
func BindForm(r *http.Request, obj any) error {
	switch ContentType(r.Header) {
	case ContentForm:
		if err := r.ParseForm(); err != nil {
			return err
		}
	case ContentMultipartForm:
		if err := r.ParseMultipartForm(MaxFormMemory); err != nil {
			if err != http.ErrNotMultipart {
				return err
			}
		}
	}
	if err := MapForm(obj, r.Form); err != nil {
		return err
	}
	return validator.ValidateStruct(obj)
}

// BindProto 解析Proto请求体并校验
func BindProto(r *http.Request, msg proto.Message) error {
	// GET请求
	if r.Method == http.MethodGet {
		err := protos.QueryToMessage(msg, r.URL.Query())
		if err != nil {
			return err
		}
		return protos.Validate(msg)
	}
	// 取请求Body
	if r.Body != nil && r.Body != http.NoBody {
		defer io.Copy(io.Discard, r.Body)
		if err := json.NewDecoder(r.Body).Decode(msg); err != nil {
			return err
		}
	}
	return protos.Validate(msg)
}

package result

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/noble-gase/ne/codes"
)

// Result the result definition for API
type Result interface {
	// JSON outputs json result
	JSON(w http.ResponseWriter, r *http.Request)
}

type result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

func (ret *result) JSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(ret); err != nil {
		w.Write([]byte(err.Error()))
	}
}

func New(code codes.Code, data ...any) Result {
	ret := &result{
		Code: code.Val(),
		Msg:  code.Msg(),
	}
	if len(data) != 0 && data[0] != nil {
		ret.Data = data[0]
	}
	return ret
}

func OK(data ...any) Result {
	return New(codes.OK, data...)
}

func Err(err error, data ...any) Result {
	var code codes.Code
	if errors.As(err, &code) {
		return New(code, data...)
	}
	return New(codes.Unknown, data...)
}

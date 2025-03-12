package result

import (
	"encoding/json"
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

func New(code codes.Code, data any) Result {
	return &result{
		Code: code.V(),
		Msg:  code.M(),
		Data: data,
	}
}

func OK(data ...any) Result {
	if len(data) == 0 {
		return New(codes.OK, nil)
	}
	return New(codes.OK, data[0])
}

func Err(err error, msg ...string) Result {
	var (
		code codes.Code
		text string
	)
	switch v := err.(type) {
	case codes.Code:
		code = v
	default:
		code = codes.Unknown
		text = err.Error()
	}
	if len(msg) != 0 {
		text = msg[0]
	}
	if len(text) != 0 {
		code = code.New(text)
	}
	return New(code, nil)
}

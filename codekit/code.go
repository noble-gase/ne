package codekit

import (
	"errors"
	"fmt"
)

// Code the code definition for API
type Code interface {
	error
	// Value returns the code value
	Value() int
	// Message returns the code message
	Message() string
	// WithMsg returns a new Code with the same value but a different message
	WithMsg(msg string) Code
	// WithMsgF returns a new Code with the same value and a formatted message
	WithMsgF(format string, args ...any) Code
}

type code struct {
	val int
	msg string
}

func (c code) Error() string {
	return fmt.Sprintf("[%d] %s", c.val, c.msg)
}

func (c code) Value() int {
	return c.val
}

func (c code) Message() string {
	return c.msg
}

func (c code) WithMsg(msg string) Code {
	return code{val: c.val, msg: msg}
}

func (c code) WithMsgF(format string, args ...any) Code {
	return code{val: c.val, msg: fmt.Sprintf(format, args...)}
}

func New(val int, msg string) Code {
	return code{val: val, msg: msg}
}

var (
	OK  = New(0, "OK")
	Err = New(-1, "System Exception")
)

// Is reports whether the err is the target code
func Is(err error, target Code) bool {
	if err == nil || target == nil {
		return err == target
	}

	var c code
	if errors.As(err, &c) {
		return c.Value() == target.Value()
	}
	return false
}

// FromError returns a Code representation of err.
func FromError(err error) Code {
	if err == nil {
		return OK
	}

	var c code
	if errors.As(err, &c) {
		return c
	}
	return Err.WithMsg(err.Error())
}

package codes

import "fmt"

// Code the code definition for API
type Code interface {
	error
	// V returns the code value
	Val() int
	// M returns the code message
	Msg() string
	// New returns a newly allocated code with the same value.
	New(msg string, args ...any) Code
}

type code struct {
	v int
	m string
}

func (c code) Error() string {
	return fmt.Sprintf("%d | %s", c.v, c.m)
}

func (c code) Val() int {
	return c.v
}

func (c code) Msg() string {
	return c.m
}

func (c code) New(format string, args ...any) Code {
	if len(args) == 0 {
		return code{v: c.v, m: format}
	}
	return code{v: c.v, m: fmt.Sprintf(format, args...)}
}

func New(val int, msg string) Code {
	return code{v: val, m: msg}
}

var (
	OK      = New(0, "OK")
	Unknown = New(-1, "System Exception")
)

// Is reports whether the err is the target code
func Is(err error, target Code) bool {
	if err == nil || target == nil {
		return err == target
	}
	c, ok := err.(Code)
	if !ok {
		return false
	}
	return c.Val() == target.Val()
}

package codes

import (
	"fmt"
)

// Code the code definition for API
type Code interface {
	error
	// V returns the code value
	V() int
	// M returns the code message
	M() string
	// New returns a newly allocated code with the same value.
	New(msg string, args ...any) Code
}

type code struct {
	v int
	m string
}

func (c code) V() int {
	return c.v
}

func (c code) M() string {
	return c.m
}

func (c code) Error() string {
	return fmt.Sprintf("%d | %s", c.v, c.m)
}

func (c code) New(msg string, args ...any) Code {
	if len(args) == 0 {
		return code{v: c.v, m: msg}
	}
	return code{v: c.v, m: fmt.Sprintf(msg, args...)}
}

func New(v int, m string) Code {
	return code{v: v, m: m}
}

var (
	OK      = New(0, "OK")
	Unknown = New(-1, "unknown")
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
	return c.V() == target.V()
}

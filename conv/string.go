package conv

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

// BoolToStr returns the string representation of a boolean value.
func BoolToStr(b bool) string {
	return strconv.FormatBool(b)
}

// IntToStr returns the string representation of a signed integer.
func IntToStr[T constraints.Signed](v T) string {
	return strconv.FormatInt(int64(v), 10)
}

// UintToStr returns the string representation of an unsigned integer.
func UintToStr[T constraints.Unsigned](v T) string {
	return strconv.FormatUint(uint64(v), 10)
}

// FloatToStr returns the string representation of a float value.
func FloatToStr[T constraints.Float](v T) string {
	return strconv.FormatFloat(float64(v), 'f', -1, 64)
}

// StrToBool parses a boolean value from a string.
func StrToBool(s string) bool {
	v, _ := strconv.ParseBool(s)
	return v
}

// StrToInt parses a signed integer from a string.
func StrToInt[T constraints.Signed](s string) T {
	v, _ := strconv.ParseInt(s, 10, 64)
	return T(v)
}

// StrToUint parses an unsigned integer from a string.
func StrToUint[T constraints.Unsigned](s string) T {
	v, _ := strconv.ParseUint(s, 10, 64)
	return T(v)
}

// StrToFloat parses a float value from a string.
func StrToFloat[T constraints.Float](s string) T {
	v, _ := strconv.ParseFloat(s, 64)
	return T(v)
}

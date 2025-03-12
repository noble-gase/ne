package conv

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

func BoolToStr(b bool) string {
	return strconv.FormatBool(b)
}

func IntToStr[T constraints.Signed](v T) string {
	return strconv.FormatInt(int64(v), 10)
}

func UintToStr[T constraints.Unsigned](v T) string {
	return strconv.FormatUint(uint64(v), 10)
}

func FloatToStr[T constraints.Float](v T) string {
	return strconv.FormatFloat(float64(v), 'f', -1, 64)
}

func StrToBool(s string) bool {
	v, _ := strconv.ParseBool(s)
	return v
}

func StrToInt[T constraints.Signed](s string) T {
	v, _ := strconv.ParseInt(s, 10, 64)
	return T(v)
}

func StrToUint[T constraints.Unsigned](s string) T {
	v, _ := strconv.ParseUint(s, 10, 64)
	return T(v)
}

func StrToFloat[T constraints.Float](s string) T {
	v, _ := strconv.ParseFloat(s, 64)
	return T(v)
}

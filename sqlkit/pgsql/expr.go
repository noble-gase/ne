package pgsql

import . "github.com/go-jet/jet/v2/postgres"

func Exprs[T any](values []T, fn func(T) Expression) []Expression {
	exprs := make([]Expression, 0, len(values))
	for _, v := range values {
		exprs = append(exprs, fn(v))
	}
	return exprs
}

func Bools(values []bool) []Expression {
	return Exprs(values, func(v bool) Expression { return Bool(v) })
}

func Ints(values []int) []Expression {
	return Exprs(values, func(v int) Expression { return Int(int64(v)) })
}

func Int8s(values []int8) []Expression {
	return Exprs(values, func(v int8) Expression { return Int8(v) })
}

func Int16s(values []int16) []Expression {
	return Exprs(values, func(v int16) Expression { return Int16(v) })
}

func Int32s(values []int32) []Expression {
	return Exprs(values, func(v int32) Expression { return Int32(v) })
}

func Int64s(values []int64) []Expression {
	return Exprs(values, func(v int64) Expression { return Int64(v) })
}

func Uints(values []uint) []Expression {
	return Exprs(values, func(v uint) Expression { return Uint64(uint64(v)) })
}

func Uint8s(values []uint8) []Expression {
	return Exprs(values, func(v uint8) Expression { return Uint8(v) })
}

func Uint16s(values []uint16) []Expression {
	return Exprs(values, func(v uint16) Expression { return Uint16(v) })
}

func Uint32s(values []uint32) []Expression {
	return Exprs(values, func(v uint32) Expression { return Uint32(v) })
}

func Uint64s(values []uint64) []Expression {
	return Exprs(values, func(v uint64) Expression { return Uint64(v) })
}

func Floats(values []float64) []Expression {
	return Exprs(values, func(v float64) Expression { return Float(v) })
}

func Decimals(values []string) []Expression {
	return Exprs(values, func(v string) Expression { return Decimal(v) })
}

func Strings(values []string) []Expression {
	return Exprs(values, func(v string) Expression { return String(v) })
}

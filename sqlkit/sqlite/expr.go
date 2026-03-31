package sqlite

import . "github.com/go-jet/jet/v2/sqlite"

func ExprList[T any](values []T, fn func(T) Expression) []Expression {
	exprs := make([]Expression, 0, len(values))
	for _, v := range values {
		exprs = append(exprs, fn(v))
	}
	return exprs
}

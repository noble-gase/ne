package array

import (
	"math/rand/v2"
	"slices"
)

type Constraint[E comparable] interface {
	Element() E
}

// In 返回指定元素是否都在集合中
func In[T comparable](list []T, elems ...T) bool {
	listLen := len(list)
	elemLen := len(elems)
	if elemLen == 0 || listLen < elemLen {
		return false
	}
	// 单元素
	if elemLen == 1 {
		return slices.Contains(list, elems[0])
	}
	// 多元素
	m := make(map[T]struct{}, listLen)
	for _, v := range list {
		m[v] = struct{}{}
	}
	for _, v := range elems {
		if _, ok := m[v]; !ok {
			return false
		}
	}
	return true
}

// InT 返回指定元素是否都在集合中
func InT[T Constraint[E], E comparable](list []T, elems ...E) bool {
	listLen := len(list)
	elemLen := len(elems)
	if elemLen == 0 || listLen < elemLen {
		return false
	}
	// 单元素
	if elemLen == 1 {
		e := elems[0]
		for _, v := range list {
			if v.Element() == e {
				return true
			}
		}
		return false
	}
	// 多元素
	m := make(map[E]struct{}, len(list))
	for _, v := range list {
		m[v.Element()] = struct{}{}
	}
	for _, v := range elems {
		if _, ok := m[v]; !ok {
			return false
		}
	}
	return true
}

// Unique 集合去重
func Unique[T comparable](list []T) []T {
	var ret []T
	if len(list) == 0 {
		return ret
	}

	ret = make([]T, 0, len(list))
	m := make(map[T]struct{}, len(list))
	for _, v := range list {
		if _, ok := m[v]; !ok {
			ret = append(ret, v)
			m[v] = struct{}{}
		}
	}
	return ret
}

// UniqueT 集合去重
func UniqueT[T Constraint[E], E comparable](list []T) []T {
	var ret []T
	if len(list) == 0 {
		return ret
	}

	ret = make([]T, 0, len(list))
	m := make(map[E]struct{}, len(list))
	for _, v := range list {
		if _, ok := m[v.Element()]; !ok {
			ret = append(ret, v)
			m[v.Element()] = struct{}{}
		}
	}
	return ret
}

// Diff 返回两个集合之间的差异
func Diff[T comparable](list1 []T, list2 []T) (ret1 []T, ret2 []T) {
	m1 := map[T]struct{}{}
	m2 := map[T]struct{}{}
	for _, v := range list1 {
		m1[v] = struct{}{}
	}
	for _, v := range list2 {
		m2[v] = struct{}{}
	}

	ret1 = make([]T, 0)
	ret2 = make([]T, 0)
	for _, v := range list1 {
		if _, ok := m2[v]; !ok {
			ret1 = append(ret1, v)
		}
	}
	for _, v := range list2 {
		if _, ok := m1[v]; !ok {
			ret2 = append(ret2, v)
		}
	}
	return ret1, ret2
}

// DiffT 返回两个集合之间的差异
func DiffT[T Constraint[E], E comparable](list1 []T, list2 []T) (ret1 []T, ret2 []T) {
	m1 := map[E]struct{}{}
	m2 := map[E]struct{}{}
	for _, v := range list1 {
		m1[v.Element()] = struct{}{}
	}
	for _, v := range list2 {
		m2[v.Element()] = struct{}{}
	}

	ret1 = make([]T, 0)
	ret2 = make([]T, 0)
	for _, v := range list1 {
		if _, ok := m2[v.Element()]; !ok {
			ret1 = append(ret1, v)
		}
	}
	for _, v := range list2 {
		if _, ok := m1[v.Element()]; !ok {
			ret2 = append(ret2, v)
		}
	}
	return ret1, ret2
}

// Without 返回不包括所有给定值的切片
func Without[T comparable](list []T, exclude ...T) []T {
	var ret []T
	if len(list) == 0 {
		return ret
	}

	m := make(map[T]struct{}, len(exclude))
	for _, v := range exclude {
		m[v] = struct{}{}
	}

	ret = make([]T, 0, len(list))
	for _, v := range list {
		if _, ok := m[v]; !ok {
			ret = append(ret, v)
		}
	}
	return ret
}

// WithoutT 返回不包括所有给定值的切片
func WithoutT[T Constraint[E], E comparable](list []T, exclude ...E) []T {
	var ret []T
	if len(list) == 0 {
		return ret
	}

	m := make(map[E]struct{}, len(exclude))
	for _, v := range exclude {
		m[v] = struct{}{}
	}

	ret = make([]T, 0, len(list))
	for _, v := range list {
		if _, ok := m[v.Element()]; !ok {
			ret = append(ret, v)
		}
	}
	return ret
}

// Intersect 返回两个集合的交集
func Intersect[T comparable](list1 []T, list2 []T) []T {
	m := make(map[T]struct{})
	for _, v := range list1 {
		m[v] = struct{}{}
	}

	ret := make([]T, 0)
	for _, v := range list2 {
		if _, ok := m[v]; ok {
			ret = append(ret, v)
		}
	}
	return ret
}

// IntersectT 返回两个集合的交集
func IntersectT[T Constraint[E], E comparable](list1 []T, list2 []T) []T {
	m := make(map[E]struct{})
	for _, v := range list1 {
		m[v.Element()] = struct{}{}
	}

	ret := make([]T, 0)
	for _, v := range list2 {
		if _, ok := m[v.Element()]; ok {
			ret = append(ret, v)
		}
	}
	return ret
}

// Union 返回两个集合的并集
func Union[T comparable](lists ...[]T) []T {
	ret := make([]T, 0)
	m := make(map[T]struct{})
	for _, list := range lists {
		for _, v := range list {
			if _, ok := m[v]; !ok {
				ret = append(ret, v)
				m[v] = struct{}{}
			}
		}
	}
	return ret
}

// UnionT 返回两个集合的并集
func UnionT[T Constraint[E], E comparable](lists ...[]T) []T {
	ret := make([]T, 0)
	m := make(map[E]struct{})
	for _, list := range lists {
		for _, v := range list {
			if _, ok := m[v.Element()]; !ok {
				ret = append(ret, v)
				m[v.Element()] = struct{}{}
			}
		}
	}
	return ret
}

// Rand 返回一个指定随机挑选个数的切片
// 若 n == -1 or n >= len(list)，则返回打乱的切片
func Rand[T any](list []T, n int) []T {
	var ret []T
	if n == 0 || n < -1 {
		return ret
	}

	len := len(list)
	ret = make([]T, len)
	copy(ret, list)
	rand.Shuffle(len, func(i, j int) {
		ret[i], ret[j] = ret[j], ret[i]
	})
	if n == -1 || n >= len {
		return ret
	}
	return ret[:n]
}

// PinTop 置顶集合中的一个元素
func PinTop[T any](list []T, index int) {
	if index <= 0 || index >= len(list) {
		return
	}
	for i := index; i > 0; i-- {
		list[i], list[i-1] = list[i-1], list[i]
	}
}

// PinTopF 置顶集合中满足条件的一个元素
func PinTopFunc[T any](list []T, fn func(v T) bool) {
	index := 0
	for i, v := range list {
		if fn(v) {
			index = i
			break
		}
	}
	for i := index; i > 0; i-- {
		list[i], list[i-1] = list[i-1], list[i]
	}
}

// Chunk 集合分片
func Chunk[T any](list []T, size int) [][]T {
	var ret [][]T
	if size <= 0 {
		return ret
	}

	len := len(list)
	mod := len % size
	cap := len / size
	if mod != 0 {
		cap += 1
	}
	end := len - mod
	ret = make([][]T, 0, cap)
	for i := 0; i < end; i += size {
		ret = append(ret, list[i:i+size])
	}
	if mod != 0 {
		ret = append(ret, list[end:len])
	}
	return ret
}

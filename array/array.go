package array

import (
	"math/rand/v2"
	"slices"
)

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
	m := make(map[T]struct{})
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

// InFunc 返回指定元素是否都在集合中
func InFunc[T any, E comparable](fn func(v T) E, list []T, elems ...T) bool {
	listLen := len(list)
	elemLen := len(elems)
	if elemLen == 0 || listLen < elemLen {
		return false
	}

	// 单元素
	if elemLen == 1 {
		e := fn(elems[0])
		for _, v := range list {
			if fn(v) == e {
				return true
			}
		}
		return false
	}

	// 多元素
	m := make(map[E]struct{})
	for _, v := range list {
		m[fn(v)] = struct{}{}
	}
	for _, v := range elems {
		if _, ok := m[fn(v)]; !ok {
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

	m := make(map[T]struct{})
	for _, v := range list {
		if _, ok := m[v]; !ok {
			ret = append(ret, v)
			m[v] = struct{}{}
		}
	}
	return ret
}

// UniqueFunc 集合去重
func UniqueFunc[T any, E comparable](fn func(v T) E, list []T) []T {
	var ret []T
	if len(list) == 0 {
		return ret
	}

	m := make(map[E]struct{})
	for _, v := range list {
		e := fn(v)
		if _, ok := m[e]; !ok {
			ret = append(ret, v)
			m[e] = struct{}{}
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

// DiffFunc 返回两个集合之间的差异
func DiffFunc[T any, E comparable](fn func(v T) E, list1 []T, list2 []T) (ret1 []T, ret2 []T) {
	m1 := map[E]struct{}{}
	m2 := map[E]struct{}{}
	for _, v := range list1 {
		m1[fn(v)] = struct{}{}
	}
	for _, v := range list2 {
		m2[fn(v)] = struct{}{}
	}

	for _, v := range list1 {
		if _, ok := m2[fn(v)]; !ok {
			ret1 = append(ret1, v)
		}
	}
	for _, v := range list2 {
		if _, ok := m1[fn(v)]; !ok {
			ret2 = append(ret2, v)
		}
	}
	return ret1, ret2
}

// Exclude 返回不包括所有给定值的切片
func Exclude[T comparable](list []T, excludes ...T) []T {
	if len(list) == 0 || len(excludes) == 0 {
		return list
	}

	var ret []T

	// 单元素
	if len(excludes) == 1 {
		e := excludes[0]
		for _, v := range list {
			if v != e {
				ret = append(ret, v)
			}
		}
		return ret
	}

	// 多元素
	m := make(map[T]struct{})
	for _, v := range excludes {
		m[v] = struct{}{}
	}
	for _, v := range list {
		if _, ok := m[v]; !ok {
			ret = append(ret, v)
		}
	}
	return ret
}

// ExcludeFunc 返回不包括所有给定值的切片
func ExcludeFunc[T any, E comparable](fn func(v T) E, list []T, excludes ...T) []T {
	if len(list) == 0 || len(excludes) == 0 {
		return list
	}

	var ret []T

	// 单元素
	if len(excludes) == 1 {
		e := fn(excludes[0])
		for _, v := range list {
			if fn(v) != e {
				ret = append(ret, v)
			}
		}
		return ret
	}

	// 多元素
	m := make(map[E]struct{})
	for _, v := range excludes {
		m[fn(v)] = struct{}{}
	}
	for _, v := range list {
		if _, ok := m[fn(v)]; !ok {
			ret = append(ret, v)
		}
	}
	return ret
}

// Intersect 返回两个集合的交集
func Intersect[T comparable](list1 []T, list2 []T) []T {
	var ret []T

	m := make(map[T]struct{})
	for _, v := range list1 {
		m[v] = struct{}{}
	}

	for _, v := range list2 {
		if _, ok := m[v]; ok {
			ret = append(ret, v)
		}
	}
	return ret
}

// IntersectFunc 返回两个集合的交集
func IntersectFunc[T any, E comparable](fn func(v T) E, list1 []T, list2 []T) []T {
	var ret []T

	m := make(map[E]struct{})
	for _, v := range list1 {
		m[fn(v)] = struct{}{}
	}

	for _, v := range list2 {
		if _, ok := m[fn(v)]; ok {
			ret = append(ret, v)
		}
	}
	return ret
}

// Union 返回集合的并集
func Union[T comparable](lists ...[]T) []T {
	var ret []T
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

// UnionFunc 返回集合的并集
func UnionFunc[T any, E comparable](fn func(v T) E, lists ...[]T) []T {
	var ret []T
	m := make(map[E]struct{})
	for _, list := range lists {
		for _, v := range list {
			e := fn(v)
			if _, ok := m[e]; !ok {
				ret = append(ret, v)
				m[e] = struct{}{}
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

// PinTopFunc 置顶集合中的一个元素
func PinTopFunc[T any](fn func(v T) bool, list []T) {
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

// Filter 过滤集合
func Filter[T any](fn func(v T) bool, list []T) []T {
	var ret []T
	for _, v := range list {
		if fn(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// Map 返回处理后的新集合
func Map[T any, E any](fn func(v T) E, list []T) []E {
	var ret []E
	if len(list) == 0 {
		return ret
	}

	ret = make([]E, 0, len(list))
	for _, v := range list {
		ret = append(ret, fn(v))
	}
	return ret
}

// UniqMap 返回处理后的新集合(去重)
func UniqMap[T any, E comparable](fn func(v T) E, list []T) []E {
	var ret []E
	if len(list) == 0 {
		return ret
	}

	m := make(map[E]struct{})
	for _, v := range list {
		e := fn(v)
		if _, ok := m[e]; !ok {
			ret = append(ret, e)
			m[e] = struct{}{}
		}
	}
	return ret
}

// FilterMap 返回过滤并处理后的新集合
func FilterMap[T any, E any](fn func(v T) (E, bool), list []T) []E {
	var ret []E
	for _, v := range list {
		if e, ok := fn(v); ok {
			ret = append(ret, e)
		}
	}
	return ret
}

// Associate 序列化一个集合为Map
func Associate[T any, K comparable, V any](fn func(v T) (K, V), list []T) map[K]V {
	m := make(map[K]V, len(list))
	for _, v := range list {
		k, e := fn(v)
		m[k] = e
	}
	return m
}

// FilterAssociate 过滤并序列化一个集合为Map
func FilterAssociate[T any, K comparable, V any](fn func(v T) (K, V, bool), list []T) map[K]V {
	m := make(map[K]V, len(list))
	for _, v := range list {
		if k, e, ok := fn(v); ok {
			m[k] = e
		}
	}
	return m
}

// Group 序列化一个集合为分组合并后的Map
func Group[T any, K comparable, V any](fn func(v T) (K, V), list []T) map[K][]V {
	m := make(map[K][]V)
	for _, v := range list {
		k, e := fn(v)
		m[k] = append(m[k], e)
	}
	return m
}

// FilterGroupBy 过滤并序列化一个集合为分组合并后的Map
func FilterGroup[T any, K comparable, V any](fn func(v T) (K, V, bool), list []T) map[K][]V {
	m := make(map[K][]V)
	for _, v := range list {
		if k, e, ok := fn(v); ok {
			m[k] = append(m[k], e)
		}
	}
	return m
}

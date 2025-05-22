package array

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Foo[T comparable] struct {
	ID   T
	Name string
}

func TestIn(t *testing.T) {
	assert.True(t, In([]int{1, 2, 3, 4, 5}, 4))
	assert.True(t, In([]int{1, 2, 3, 4, 5}, 2, 4))
	assert.True(t, In([]int64{1, 2, 3, 4, 5}, 2, 4))
	assert.True(t, In([]float64{1.01, 2.02, 3.03, 4.04, 5.05}, 2.02, 4.04))
	assert.True(t, In([]string{"h", "e", "l", "l", "o"}, "e", "o"))
}

func TestInFunc(t *testing.T) {
	fooArr := []Foo[int]{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}
	assert.True(t, InFunc(func(v Foo[int]) int { return v.ID }, fooArr, Foo[int]{ID: 2}))
	assert.True(t, InFunc(func(v Foo[int]) int { return v.ID }, fooArr, Foo[int]{ID: 2}, Foo[int]{ID: 4}))

	barArr := []Foo[string]{{ID: "1"}, {ID: "2"}, {ID: "3"}, {ID: "4"}, {ID: "5"}}
	assert.True(t, InFunc(func(v Foo[string]) string { return v.ID }, barArr, Foo[string]{ID: "2"}))
	assert.True(t, InFunc(func(v Foo[string]) string { return v.ID }, barArr, Foo[string]{ID: "2"}, Foo[string]{ID: "4"}))
}

func TestUnique(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3, 4}, Unique([]int{1, 2, 1, 3, 4, 3}))
	assert.Equal(t, []int64{1, 2, 3, 4}, Unique([]int64{1, 2, 1, 3, 4, 3}))
	assert.Equal(t, []float64{1.01, 2.02, 3.03, 4.04}, Unique([]float64{1.01, 2.02, 1.01, 3.03, 4.04, 3.03}))
	assert.Equal(t, []string{"h", "e", "l", "o"}, Unique([]string{"h", "e", "l", "l", "o"}))
}

func TestUniqueFunc(t *testing.T) {
	fooArr := []Foo[int]{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 2}, {ID: 3}}
	assert.Equal(t, UniqueFunc(func(v Foo[int]) int { return v.ID }, fooArr), []Foo[int]{{ID: 1}, {ID: 2}, {ID: 3}})

	barArr := []Foo[string]{{ID: "1"}, {ID: "2"}, {ID: "3"}, {ID: "2"}, {ID: "3"}}
	assert.Equal(t, UniqueFunc(func(v Foo[string]) string { return v.ID }, barArr), []Foo[string]{{ID: "1"}, {ID: "2"}, {ID: "3"}})
}

func TestDiff(t *testing.T) {
	left1, right1 := Diff([]int{0, 1, 2, 3, 4, 5}, []int{0, 2, 6})
	assert.Equal(t, []int{1, 3, 4, 5}, left1)
	assert.Equal(t, []int{6}, right1)

	left2, right2 := Diff([]int{1, 2, 3, 4, 5}, []int{0, 6})
	assert.Equal(t, []int{1, 2, 3, 4, 5}, left2)
	assert.Equal(t, []int{0, 6}, right2)

	left3, right3 := Diff([]int{0, 1, 2, 3, 4, 5}, []int{0, 1, 2, 3, 4, 5})
	assert.Nil(t, left3)
	assert.Nil(t, right3)
}

func TestDiffFunc(t *testing.T) {
	left1, right1 := DiffFunc(func(v Foo[int]) int { return v.ID }, []Foo[int]{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}, []Foo[int]{{ID: 0}, {ID: 2}, {ID: 6}})
	assert.Equal(t, []Foo[int]{{ID: 1}, {ID: 3}, {ID: 4}, {ID: 5}}, left1)
	assert.Equal(t, []Foo[int]{{ID: 6}}, right1)

	left2, right2 := DiffFunc(func(v Foo[int]) int { return v.ID }, []Foo[int]{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}, []Foo[int]{{ID: 0}, {ID: 6}})
	assert.Equal(t, []Foo[int]{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}, left2)
	assert.Equal(t, []Foo[int]{{ID: 0}, {ID: 6}}, right2)

	left3, right3 := DiffFunc(func(v Foo[int]) int { return v.ID }, []Foo[int]{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}, []Foo[int]{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}})
	assert.Nil(t, left3)
	assert.Nil(t, right3)
}

func TestExclude(t *testing.T) {
	result1 := Exclude([]int{0, 2, 10}, 0, 1, 2, 3, 4, 5)
	assert.Equal(t, []int{10}, result1)

	result2 := Exclude([]int{0, 7}, 0, 1, 2, 3, 4, 5)
	assert.Equal(t, []int{7}, result2)

	result3 := Exclude([]int{}, 0, 1, 2, 3, 4, 5)
	assert.Equal(t, 0, len(result3))

	result4 := Exclude([]int{0, 1, 2}, 0, 1, 2)
	assert.Equal(t, 0, len(result4))

	result5 := Exclude([]int{})
	assert.Equal(t, 0, len(result5))
}

func TestExcludeFunc(t *testing.T) {
	result1 := ExcludeFunc(func(v Foo[int]) int { return v.ID }, []Foo[int]{{ID: 0}, {ID: 2}, {ID: 10}}, Foo[int]{ID: 0}, Foo[int]{ID: 1}, Foo[int]{ID: 2}, Foo[int]{ID: 3}, Foo[int]{ID: 4}, Foo[int]{ID: 5})
	assert.Equal(t, []Foo[int]{{ID: 10}}, result1)

	result2 := ExcludeFunc(func(v Foo[int]) int { return v.ID }, []Foo[int]{{ID: 0}, {ID: 7}}, Foo[int]{ID: 0}, Foo[int]{ID: 1}, Foo[int]{ID: 2}, Foo[int]{ID: 3}, Foo[int]{ID: 4}, Foo[int]{ID: 5})
	assert.Equal(t, []Foo[int]{{ID: 7}}, result2)

	result3 := ExcludeFunc(func(v Foo[int]) int { return v.ID }, []Foo[int]{}, Foo[int]{ID: 0}, Foo[int]{ID: 1}, Foo[int]{ID: 2}, Foo[int]{ID: 3}, Foo[int]{ID: 4}, Foo[int]{ID: 5})
	assert.Equal(t, 0, len(result3))

	result4 := ExcludeFunc(func(v Foo[int]) int { return v.ID }, []Foo[int]{{ID: 0}, {ID: 1}, {ID: 2}}, Foo[int]{ID: 0}, Foo[int]{ID: 1}, Foo[int]{ID: 2})
	assert.Equal(t, 0, len(result4))

	result5 := ExcludeFunc(func(v Foo[int]) int { return v.ID }, []Foo[int]{})
	assert.Equal(t, 0, len(result5))
}

func TestIntersect(t *testing.T) {
	result1 := Intersect([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
	assert.Equal(t, []int{0, 2}, result1)

	result2 := Intersect([]int{0, 1, 2, 3, 4, 5}, []int{0, 6})
	assert.Equal(t, []int{0}, result2)

	result3 := Intersect([]int{0, 1, 2, 3, 4, 5}, []int{-1, 6})
	assert.Equal(t, 0, len(result3))

	result4 := Intersect([]int{0, 6}, []int{0, 1, 2, 3, 4, 5})
	assert.Equal(t, []int{0}, result4)

	result5 := Intersect([]int{0, 6, 0}, []int{0, 1, 2, 3, 4, 5})
	assert.Equal(t, []int{0}, result5)
}

func TestIntersectFunc(t *testing.T) {
	result1 := IntersectFunc(func(v Foo[int]) int { return v.ID }, []Foo[int]{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}, []Foo[int]{{ID: 0}, {ID: 2}})
	assert.Equal(t, []Foo[int]{{ID: 0}, {ID: 2}}, result1)

	result2 := IntersectFunc(func(v Foo[int]) int { return v.ID }, []Foo[int]{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}, []Foo[int]{{ID: 0}, {ID: 6}})
	assert.Equal(t, []Foo[int]{{ID: 0}}, result2)

	result3 := IntersectFunc(func(v Foo[int]) int { return v.ID }, []Foo[int]{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}, []Foo[int]{{ID: -1}, {ID: 6}})
	assert.Equal(t, 0, len(result3))

	result4 := IntersectFunc(func(v Foo[int]) int { return v.ID }, []Foo[int]{{ID: 0}, {ID: 6}}, []Foo[int]{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}})
	assert.Equal(t, []Foo[int]{{ID: 0}}, result4)

	result5 := IntersectFunc(func(v Foo[int]) int { return v.ID }, []Foo[int]{{ID: 0}, {ID: 6}, {ID: 0}}, []Foo[int]{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}})
	assert.Equal(t, []Foo[int]{{ID: 0}}, result5)
}

func TestUnion(t *testing.T) {
	result1 := Union([]int{0, 1, 2, 3, 4, 5}, []int{0, 2, 10}, []int{0, 1, 11})
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 10, 11}, result1)

	result2 := Union([]int{0, 1, 2, 3, 4, 5}, []int{6, 7}, []int{8, 9})
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, result2)

	result3 := Union([]int{0, 1, 2, 3, 4, 5}, []int{}, []int{})
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5}, result3)

	result4 := Union([]int{0, 1, 2}, []int{0, 1, 2}, []int{0, 1, 2})
	assert.Equal(t, []int{0, 1, 2}, result4)

	result5 := Union([]int{}, []int{}, []int{})
	assert.Equal(t, 0, len(result5))
}

func TestUnionFunc(t *testing.T) {
	result1 := UnionFunc(func(v Foo[int]) int { return v.ID }, []Foo[int]{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}, []Foo[int]{{ID: 0}, {ID: 2}, {ID: 10}}, []Foo[int]{{ID: 0}, {ID: 1}, {ID: 11}})
	assert.Equal(t, []Foo[int]{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 10}, {ID: 11}}, result1)

	result2 := UnionFunc(func(v Foo[int]) int { return v.ID }, []Foo[int]{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}, []Foo[int]{{ID: 6}, {ID: 7}}, []Foo[int]{{ID: 8}, {ID: 9}})
	assert.Equal(t, []Foo[int]{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7}, {ID: 8}, {ID: 9}}, result2)

	result3 := UnionFunc(func(v Foo[int]) int { return v.ID }, []Foo[int]{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}, []Foo[int]{}, []Foo[int]{})
	assert.Equal(t, []Foo[int]{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}, result3)

	result4 := UnionFunc(func(v Foo[int]) int { return v.ID }, []Foo[int]{{ID: 0}, {ID: 1}, {ID: 2}}, []Foo[int]{{ID: 0}, {ID: 1}, {ID: 2}}, []Foo[int]{{ID: 0}, {ID: 1}, {ID: 2}})
	assert.Equal(t, []Foo[int]{{ID: 0}, {ID: 1}, {ID: 2}}, result4)

	result5 := UnionFunc(func(v Foo[int]) int { return v.ID }, []Foo[int]{}, []Foo[int]{}, []Foo[int]{})
	assert.Equal(t, 0, len(result5))
}

func TestRand(t *testing.T) {
	arr1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ret1 := Rand(arr1, 2)
	assert.Equal(t, 2, len(ret1))
	assert.NotEqual(t, arr1[:2], ret1)

	arr2 := []float64{1.01, 2.02, 3.03, 4.04, 5.05, 6.06, 7.07, 8.08, 9.09, 10.10}
	ret2 := Rand(arr2, 2)
	assert.Equal(t, 2, len(ret2))
	assert.NotEqual(t, arr2[:2], ret2)

	arr3 := []string{"h", "e", "l", "l", "o", "w", "o", "r", "l", "d"}
	ret3 := Rand(arr3, 2)
	assert.Equal(t, 2, len(ret3))
	assert.NotEqual(t, arr3[:2], ret3)

	arr4 := []Foo[int]{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7}, {ID: 8}, {ID: 9}, {ID: 10}}

	ret4 := Rand(arr4, 2)
	assert.Equal(t, 2, len(ret4))
	assert.NotEqual(t, arr4[:2], ret4)

	ret5 := Rand(arr4, -1)
	assert.Equal(t, len(arr4), len(ret5))
	assert.NotEqual(t, arr4, ret5)
}

func TestPinTop(t *testing.T) {
	arr1 := []int{1, 2, 3, 4, 5}
	PinTop(arr1, 3)
	assert.Equal(t, []int{4, 1, 2, 3, 5}, arr1)

	arr2 := []int{1, 2, 3, 4, 5}
	PinTop(arr2, 0)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, arr2)

	arr3 := []int{1, 2, 3, 4, 5}
	PinTop(arr3, -1)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, arr3)

	arr4 := []int{1, 2, 3, 4, 5}
	PinTop(arr4, 5)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, arr4)
}

func TestPinTopF(t *testing.T) {
	arr := []Foo[int]{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}
	PinTopFunc(func(v Foo[int]) bool {
		return v.ID == 4
	}, arr)
	assert.Equal(t, []Foo[int]{{ID: 4}, {ID: 1}, {ID: 2}, {ID: 3}, {ID: 5}}, arr)
}

func TestChunk(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	ret1 := Chunk(arr, 2)
	assert.Equal(t, [][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}, {9, 10}}, ret1)

	ret2 := Chunk(arr, 3)
	assert.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10}}, ret2)

	ret3 := Chunk(arr, 4)
	assert.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10}}, ret3)
}

func TestFilter(t *testing.T) {
	arr1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ret1 := Filter(func(v int) bool {
		return v%2 == 0
	}, arr1)
	assert.Equal(t, []int{2, 4, 6, 8, 10}, ret1)

	arr2 := []Foo[int]{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7}, {ID: 8}, {ID: 9}, {ID: 10}}
	ret2 := Filter(func(v Foo[int]) bool {
		return v.ID%2 == 0
	}, arr2)
	assert.Equal(t, []Foo[int]{{ID: 2}, {ID: 4}, {ID: 6}, {ID: 8}, {ID: 10}}, ret2)
}

func TestMap(t *testing.T) {
	arr1 := []Foo[int]{{ID: 1}, {ID: 2}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 4}, {ID: 5}}
	ret1 := Map(func(v Foo[int]) int {
		return v.ID
	}, arr1)
	assert.Equal(t, []int{1, 2, 2, 3, 4, 4, 5}, ret1)

	arr2 := []Foo[int]{{ID: 1}, {ID: 2}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 4}, {ID: 5}}
	ret2 := Map(func(v Foo[int]) int {
		return v.ID * 2
	}, arr2)
	assert.Equal(t, []int{2, 4, 4, 6, 8, 8, 10}, ret2)
}

func TestUniqMap(t *testing.T) {
	arr1 := []Foo[int]{{ID: 1}, {ID: 2}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 4}, {ID: 5}}
	ret1 := UniqMap(func(v Foo[int]) int {
		return v.ID
	}, arr1)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, ret1)

	arr2 := []Foo[int]{{ID: 1}, {ID: 2}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 4}, {ID: 5}}
	ret2 := UniqMap(func(v Foo[int]) int {
		return v.ID * 2
	}, arr2)
	assert.Equal(t, []int{2, 4, 6, 8, 10}, ret2)
}

func TestFilterMap(t *testing.T) {
	arr1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ret1 := FilterMap(func(v int) (int, bool) {
		return v * 2, v%2 == 0
	}, arr1)
	assert.Equal(t, []int{4, 8, 12, 16, 20}, ret1)

	arr2 := []Foo[int]{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7}, {ID: 8}, {ID: 9}, {ID: 10}}
	ret2 := FilterMap(func(v Foo[int]) (int, bool) {
		return v.ID * 2, v.ID%2 == 0
	}, arr2)
	assert.Equal(t, []int{4, 8, 12, 16, 20}, ret2)
}

func TestAssociate(t *testing.T) {
	arr := []Foo[int]{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}
	ret := Associate(func(v Foo[int]) (int, int) {
		return v.ID, v.ID * 2
	}, arr)
	assert.Equal(t, map[int]int{1: 2, 2: 4, 3: 6, 4: 8, 5: 10}, ret)
}

func TestFilterAssociate(t *testing.T) {
	arr := []Foo[int]{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}}
	ret := FilterAssociate(func(v Foo[int]) (int, int, bool) {
		return v.ID, v.ID * 2, v.ID%2 == 0
	}, arr)
	assert.Equal(t, map[int]int{2: 4, 4: 8}, ret)
}

func TestGroup(t *testing.T) {
	arr := []Foo[int]{{ID: 1, Name: "a"}, {ID: 1, Name: "b"}, {ID: 1, Name: "c"}, {ID: 2, Name: "d"}, {ID: 2, Name: "e"}}
	ret := Group(func(v Foo[int]) (int, string) {
		return v.ID, v.Name
	}, arr)
	assert.Equal(t, map[int][]string{1: {"a", "b", "c"}, 2: {"d", "e"}}, ret)
}

func TestFilterGroup(t *testing.T) {
	arr := []Foo[int]{{ID: 1, Name: "a"}, {ID: 1, Name: "b"}, {ID: 1, Name: "c"}, {ID: 2, Name: "d"}, {ID: 2, Name: "e"}}
	ret := FilterGroup(func(v Foo[int]) (int, string, bool) {
		return v.ID, v.Name, v.ID%2 == 0
	}, arr)
	assert.Equal(t, map[int][]string{2: {"d", "e"}}, ret)
}

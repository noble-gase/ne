package array

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Foo struct {
	ID   int
	Name string
}

func (f *Foo) Element() int {
	return f.ID
}

type Bar struct {
	ID   string
	Name string
}

func (b *Bar) Element() string {
	return b.ID
}

func TestIn(t *testing.T) {
	assert.True(t, In([]int{1, 2, 3, 4, 5}, 4))
	assert.True(t, In([]int{1, 2, 3, 4, 5}, 2, 4))
	assert.True(t, In([]int64{1, 2, 3, 4, 5}, 2, 4))
	assert.True(t, In([]float64{1.01, 2.02, 3.03, 4.04, 5.05}, 2.02, 4.04))
	assert.True(t, In([]string{"h", "e", "l", "l", "o"}, "e", "o"))
}

func TestInT(t *testing.T) {
	fooArr := []*Foo{
		{
			ID:   1,
			Name: "foo-1",
		},
		{
			ID:   2,
			Name: "foo-2",
		},
		{
			ID:   3,
			Name: "foo-3",
		},
		{
			ID:   4,
			Name: "foo-4",
		},
		{
			ID:   5,
			Name: "foo-5",
		},
	}
	assert.True(t, InT(fooArr, 2))
	assert.True(t, InT(fooArr, 2, 4))

	barArr := []*Bar{
		{
			ID:   "1",
			Name: "foo-1",
		},
		{
			ID:   "2",
			Name: "foo-2",
		},
		{
			ID:   "3",
			Name: "foo-3",
		},
		{
			ID:   "4",
			Name: "foo-4",
		},
		{
			ID:   "5",
			Name: "foo-5",
		},
	}
	assert.True(t, InT(barArr, "2"))
	assert.True(t, InT(barArr, "2", "4"))
}

func TestUniq(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3, 4}, Unique([]int{1, 2, 1, 3, 4, 3}))
	assert.Equal(t, []int64{1, 2, 3, 4}, Unique([]int64{1, 2, 1, 3, 4, 3}))
	assert.Equal(t, []float64{1.01, 2.02, 3.03, 4.04}, Unique([]float64{1.01, 2.02, 1.01, 3.03, 4.04, 3.03}))
	assert.Equal(t, []string{"h", "e", "l", "o"}, Unique([]string{"h", "e", "l", "l", "o"}))
}

func TestUniqT(t *testing.T) {
	fooArr := []*Foo{
		{
			ID:   1,
			Name: "foo-1",
		},
		{
			ID:   2,
			Name: "foo-2",
		},
		{
			ID:   2,
			Name: "foo-2",
		},
		{
			ID:   3,
			Name: "foo-3",
		},
		{
			ID:   3,
			Name: "foo-3",
		},
	}
	assert.Equal(t, UniqueT(fooArr), []*Foo{
		{
			ID:   1,
			Name: "foo-1",
		},
		{
			ID:   2,
			Name: "foo-2",
		},
		{
			ID:   3,
			Name: "foo-3",
		},
	})

	barArr := []*Bar{
		{
			ID:   "1",
			Name: "foo-1",
		},
		{
			ID:   "2",
			Name: "foo-2",
		},
		{
			ID:   "2",
			Name: "foo-2",
		},
		{
			ID:   "3",
			Name: "foo-3",
		},
		{
			ID:   "3",
			Name: "foo-3",
		},
	}
	assert.Equal(t, UniqueT(barArr), []*Bar{
		{
			ID:   "1",
			Name: "foo-1",
		},
		{
			ID:   "2",
			Name: "foo-2",
		},
		{
			ID:   "3",
			Name: "foo-3",
		},
	})
}

func TestDiff(t *testing.T) {
	left1, right1 := Diff([]int{0, 1, 2, 3, 4, 5}, []int{0, 2, 6})
	assert.Equal(t, []int{1, 3, 4, 5}, left1)
	assert.Equal(t, []int{6}, right1)

	left2, right2 := Diff([]int{1, 2, 3, 4, 5}, []int{0, 6})
	assert.Equal(t, []int{1, 2, 3, 4, 5}, left2)
	assert.Equal(t, []int{0, 6}, right2)

	left3, right3 := Diff([]int{0, 1, 2, 3, 4, 5}, []int{0, 1, 2, 3, 4, 5})
	assert.Equal(t, []int{}, left3)
	assert.Equal(t, []int{}, right3)
}

func TestDiffT(t *testing.T) {
	fooArr1 := []*Foo{
		{
			ID:   1,
			Name: "foo-1",
		},
		{
			ID:   2,
			Name: "foo-2",
		},
		{
			ID:   3,
			Name: "foo-3",
		},
		{
			ID:   4,
			Name: "foo-4",
		},
		{
			ID:   5,
			Name: "foo-5",
		},
	}
	fooArr2 := []*Foo{
		{
			ID:   2,
			Name: "foo-2",
		},
		{
			ID:   4,
			Name: "foo-4",
		},
		{
			ID:   6,
			Name: "foo-6",
		},
	}
	left1, right1 := DiffT(fooArr1, fooArr2)
	assert.Equal(t, left1, []*Foo{
		{
			ID:   1,
			Name: "foo-1",
		},
		{
			ID:   3,
			Name: "foo-3",
		},
		{
			ID:   5,
			Name: "foo-5",
		},
	})
	assert.Equal(t, right1, []*Foo{
		{
			ID:   6,
			Name: "foo-6",
		},
	})

	barArr1 := []*Bar{
		{
			ID:   "1",
			Name: "foo-1",
		},
		{
			ID:   "2",
			Name: "foo-2",
		},
		{
			ID:   "3",
			Name: "foo-3",
		},
		{
			ID:   "4",
			Name: "foo-4",
		},
		{
			ID:   "5",
			Name: "foo-5",
		},
	}
	barArr2 := []*Bar{
		{
			ID:   "2",
			Name: "foo-2",
		},
		{
			ID:   "4",
			Name: "foo-4",
		},
		{
			ID:   "6",
			Name: "foo-6",
		},
	}
	left2, right2 := DiffT(barArr1, barArr2)
	assert.Equal(t, left2, []*Bar{
		{
			ID:   "1",
			Name: "foo-1",
		},
		{
			ID:   "3",
			Name: "foo-3",
		},
		{
			ID:   "5",
			Name: "foo-5",
		},
	})
	assert.Equal(t, right2, []*Bar{
		{
			ID:   "6",
			Name: "foo-6",
		},
	})
}

func TestWithout(t *testing.T) {
	result1 := Without([]int{0, 2, 10}, 0, 1, 2, 3, 4, 5)
	assert.Equal(t, []int{10}, result1)

	result2 := Without([]int{0, 7}, 0, 1, 2, 3, 4, 5)
	assert.Equal(t, []int{7}, result2)

	result3 := Without([]int{}, 0, 1, 2, 3, 4, 5)
	assert.Nil(t, result3)

	result4 := Without([]int{0, 1, 2}, 0, 1, 2)
	assert.Equal(t, []int{}, result4)

	result5 := Without([]int{})
	assert.Nil(t, result5)
}

func TestWithoutT(t *testing.T) {
	fooArr := []*Foo{
		{
			ID:   1,
			Name: "foo-1",
		},
		{
			ID:   2,
			Name: "foo-2",
		},
		{
			ID:   3,
			Name: "foo-3",
		},
	}
	result1 := WithoutT(fooArr, 1, 3)
	assert.Equal(t, result1, []*Foo{
		{
			ID:   2,
			Name: "foo-2",
		},
	})

	barArr := []*Bar{
		{
			ID:   "1",
			Name: "foo-1",
		},
		{
			ID:   "2",
			Name: "foo-2",
		},
		{
			ID:   "3",
			Name: "foo-3",
		},
	}
	result2 := WithoutT(barArr, "1", "3")
	assert.Equal(t, result2, []*Bar{
		{
			ID:   "2",
			Name: "foo-2",
		},
	})
}

func TestIntersect(t *testing.T) {
	result1 := Intersect([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
	assert.Equal(t, []int{0, 2}, result1)

	result2 := Intersect([]int{0, 1, 2, 3, 4, 5}, []int{0, 6})
	assert.Equal(t, []int{0}, result2)

	result3 := Intersect([]int{0, 1, 2, 3, 4, 5}, []int{-1, 6})
	assert.Equal(t, []int{}, result3)

	result4 := Intersect([]int{0, 6}, []int{0, 1, 2, 3, 4, 5})
	assert.Equal(t, []int{0}, result4)

	result5 := Intersect([]int{0, 6, 0}, []int{0, 1, 2, 3, 4, 5})
	assert.Equal(t, []int{0}, result5)
}

func TestIntersectT(t *testing.T) {
	fooArr1 := []*Foo{
		{
			ID:   1,
			Name: "foo-1",
		},
		{
			ID:   2,
			Name: "foo-2",
		},
		{
			ID:   3,
			Name: "foo-3",
		},
		{
			ID:   4,
			Name: "foo-4",
		},
		{
			ID:   5,
			Name: "foo-5",
		},
	}
	fooArr2 := []*Foo{
		{
			ID:   2,
			Name: "foo-2",
		},
		{
			ID:   4,
			Name: "foo-4",
		},
		{
			ID:   6,
			Name: "foo-6",
		},
	}
	result1 := IntersectT(fooArr1, fooArr2)
	assert.Equal(t, result1, []*Foo{
		{
			ID:   2,
			Name: "foo-2",
		},
		{
			ID:   4,
			Name: "foo-4",
		},
	})

	barArr1 := []*Bar{
		{
			ID:   "1",
			Name: "foo-1",
		},
		{
			ID:   "2",
			Name: "foo-2",
		},
		{
			ID:   "3",
			Name: "foo-3",
		},
		{
			ID:   "4",
			Name: "foo-4",
		},
		{
			ID:   "5",
			Name: "foo-5",
		},
	}
	barArr2 := []*Bar{
		{
			ID:   "2",
			Name: "foo-2",
		},
		{
			ID:   "4",
			Name: "foo-4",
		},
		{
			ID:   "6",
			Name: "foo-6",
		},
	}
	result2 := IntersectT(barArr1, barArr2)
	assert.Equal(t, result2, []*Bar{
		{
			ID:   "2",
			Name: "foo-2",
		},
		{
			ID:   "4",
			Name: "foo-4",
		},
	})
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
	assert.Equal(t, []int{}, result5)
}

func TestUnionT(t *testing.T) {
	fooArr1 := []*Foo{
		{
			ID:   1,
			Name: "foo-1",
		},
	}
	fooArr2 := []*Foo{
		{
			ID:   2,
			Name: "foo-2",
		},
		{
			ID:   4,
			Name: "foo-4",
		},
	}
	fooArr3 := []*Foo{
		{
			ID:   3,
			Name: "foo-3",
		},
		{
			ID:   5,
			Name: "foo-5",
		},
	}
	result1 := UnionT(fooArr1, fooArr2, fooArr3)
	assert.Equal(t, result1, []*Foo{
		{
			ID:   1,
			Name: "foo-1",
		},
		{
			ID:   2,
			Name: "foo-2",
		},
		{
			ID:   4,
			Name: "foo-4",
		},
		{
			ID:   3,
			Name: "foo-3",
		},
		{
			ID:   5,
			Name: "foo-5",
		},
	})

	barArr1 := []*Bar{
		{
			ID:   "1",
			Name: "foo-1",
		},
	}
	barArr2 := []*Bar{
		{
			ID:   "2",
			Name: "foo-2",
		},
		{
			ID:   "4",
			Name: "foo-4",
		},
	}
	barArr3 := []*Bar{
		{
			ID:   "3",
			Name: "foo-3",
		},
		{
			ID:   "5",
			Name: "foo-5",
		},
	}
	result2 := UnionT(barArr1, barArr2, barArr3)
	assert.Equal(t, result2, []*Bar{
		{
			ID:   "1",
			Name: "foo-1",
		},
		{
			ID:   "2",
			Name: "foo-2",
		},
		{
			ID:   "4",
			Name: "foo-4",
		},
		{
			ID:   "3",
			Name: "foo-3",
		},
		{
			ID:   "5",
			Name: "foo-5",
		},
	})
}

func TestRand(t *testing.T) {
	a1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ret1 := Rand(a1, 2)
	assert.Equal(t, 2, len(ret1))
	assert.NotEqual(t, a1[:2], ret1)

	a2 := []float64{1.01, 2.02, 3.03, 4.04, 5.05, 6.06, 7.07, 8.08, 9.09, 10.10}
	ret2 := Rand(a2, 2)
	assert.Equal(t, 2, len(ret2))
	assert.NotEqual(t, a2[:2], ret2)

	a3 := []string{"h", "e", "l", "l", "o", "w", "o", "r", "l", "d"}
	ret3 := Rand(a3, 2)
	assert.Equal(t, 2, len(ret3))
	assert.NotEqual(t, a3[:2], ret3)

	type User struct {
		ID   int64
		Name string
	}

	a4 := []User{
		{
			ID:   1,
			Name: "h",
		},
		{
			ID:   2,
			Name: "e",
		},
		{
			ID:   3,
			Name: "l",
		},
		{
			ID:   4,
			Name: "l",
		},
		{
			ID:   5,
			Name: "o",
		},
		{
			ID:   6,
			Name: "w",
		},
		{
			ID:   7,
			Name: "o",
		},
		{
			ID:   8,
			Name: "r",
		},
		{
			ID:   9,
			Name: "l",
		},
		{
			ID:   10,
			Name: "d",
		},
	}

	ret4 := Rand(a4, 2)
	assert.Equal(t, 2, len(ret4))
	assert.NotEqual(t, a4[:2], ret4)

	ret5 := Rand(a4, -1)
	assert.Equal(t, len(a4), len(ret5))
	assert.NotEqual(t, a4, ret5)
}

func TestPinTop(t *testing.T) {
	a1 := []int{1, 2, 3, 4, 5}
	PinTop(a1, 3)
	assert.Equal(t, []int{4, 1, 2, 3, 5}, a1)

	a2 := []int{1, 2, 3, 4, 5}
	PinTop(a1, 0)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, a2)

	a3 := []int{1, 2, 3, 4, 5}
	PinTop(a1, -1)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, a3)

	a4 := []int{1, 2, 3, 4, 5}
	PinTop(a1, 5)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, a4)
}

func TestPinTopF(t *testing.T) {
	type Demo struct {
		ID   int
		Name string
	}
	arr := []Demo{
		{
			ID:   1,
			Name: "h",
		},
		{
			ID:   2,
			Name: "e",
		},
		{
			ID:   3,
			Name: "l",
		},
		{
			ID:   4,
			Name: "o",
		},
		{
			ID:   5,
			Name: "w",
		},
	}
	PinTopFunc(arr, func(v Demo) bool {
		return v.Name == "o"
	})
	assert.Equal(t, []Demo{
		{
			ID:   4,
			Name: "o",
		},
		{
			ID:   1,
			Name: "h",
		},
		{
			ID:   2,
			Name: "e",
		},
		{
			ID:   3,
			Name: "l",
		},
		{
			ID:   5,
			Name: "w",
		},
	}, arr)
}

func TestChunk(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	ret1 := Chunk(a, 2)
	assert.Equal(t, [][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}, {9, 10}}, ret1)

	ret2 := Chunk(a, 3)
	assert.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10}}, ret2)

	ret3 := Chunk(a, 4)
	assert.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10}}, ret3)
}

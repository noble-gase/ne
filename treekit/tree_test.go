package treekit

import (
	"encoding/json"
	"testing"
)

type Foo struct {
	Uid  int    `json:"uid"`
	Pid  int    `json:"pid"`
	Name string `json:"name"`
}

func (f *Foo) ID() int {
	return f.Uid
}

func (f *Foo) BelongTo() int {
	return f.Pid
}

func TestTreeInt(t *testing.T) {
	data := []*Foo{
		{
			Uid:  1,
			Pid:  0,
			Name: "foo-1",
		},
		{
			Uid:  2,
			Pid:  0,
			Name: "foo-2",
		},
		{
			Uid:  3,
			Pid:  1,
			Name: "foo-3",
		},
		{
			Uid:  4,
			Pid:  1,
			Name: "foo-4",
		},
		{
			Uid:  5,
			Pid:  2,
			Name: "foo-5",
		},
		{
			Uid:  6,
			Pid:  2,
			Name: "foo-6",
		},
		{
			Uid:  7,
			Pid:  3,
			Name: "foo-7",
		},
		{
			Uid:  8,
			Pid:  3,
			Name: "foo-8",
		},
		{
			Uid:  9,
			Pid:  4,
			Name: "foo-9",
		},
		{
			Uid:  10,
			Pid:  4,
			Name: "foo-10",
		},
	}

	tree := NewTree(data, 0)
	b, _ := json.Marshal(tree)
	t.Log(string(b))
}

type Bar struct {
	Uid  string `json:"uid"`
	Pid  string `json:"pid"`
	Name string `json:"name"`
}

func (b *Bar) ID() string {
	return b.Uid
}

func (b *Bar) BelongTo() string {
	return b.Pid
}

func TestTreeStr(t *testing.T) {
	data := []*Bar{
		{
			Uid:  "1",
			Pid:  "",
			Name: "bar-1",
		},
		{
			Uid:  "2",
			Pid:  "",
			Name: "bar-2",
		},
		{
			Uid:  "3",
			Pid:  "1",
			Name: "bar-3",
		},
		{
			Uid:  "4",
			Pid:  "1",
			Name: "bar-4",
		},
		{
			Uid:  "5",
			Pid:  "2",
			Name: "bar-5",
		},
		{
			Uid:  "6",
			Pid:  "2",
			Name: "bar-6",
		},
		{
			Uid:  "7",
			Pid:  "3",
			Name: "bar-7",
		},
		{
			Uid:  "8",
			Pid:  "3",
			Name: "bar-8",
		},
		{
			Uid:  "9",
			Pid:  "4",
			Name: "bar-9",
		},
		{
			Uid:  "10",
			Pid:  "4",
			Name: "bar-10",
		},
	}

	tree := NewTree(data, "")
	b, _ := json.Marshal(tree)
	t.Log(string(b))
}

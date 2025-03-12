package leveltree

import (
	"encoding/json"
	"testing"
)

type Foo struct {
	ID   int    `json:"id"`
	Pid  int    `json:"pid"`
	Name string `json:"name"`
}

func (f *Foo) GetID() int {
	return f.ID
}

func (f *Foo) GetPid() int {
	return f.Pid
}

func TestTreeInt(t *testing.T) {
	data := map[int][]*Foo{
		0: {
			{
				ID:   1,
				Pid:  0,
				Name: "foo-1",
			},
			{
				ID:   2,
				Pid:  0,
				Name: "foo-2",
			},
		},
		1: {
			{
				ID:   3,
				Pid:  1,
				Name: "foo-3",
			},
			{
				ID:   4,
				Pid:  1,
				Name: "foo-4",
			},
		},
		2: {
			{
				ID:   5,
				Pid:  2,
				Name: "foo-5",
			},
			{
				ID:   6,
				Pid:  2,
				Name: "foo-6",
			},
		},
		3: {
			{
				ID:   7,
				Pid:  3,
				Name: "foo-7",
			},
			{
				ID:   8,
				Pid:  3,
				Name: "foo-8",
			},
		},
		4: {
			{
				ID:   9,
				Pid:  4,
				Name: "foo-9",
			},
			{
				ID:   10,
				Pid:  4,
				Name: "foo-10",
			},
		},
	}

	tree := New(data, 0)
	b, _ := json.Marshal(tree)
	t.Log(string(b))
}

type Bar struct {
	ID   string `json:"id"`
	Pid  string `json:"pid"`
	Name string `json:"name"`
}

func (b *Bar) GetID() string {
	return b.ID
}

func (b *Bar) GetPid() string {
	return b.Pid
}

func TestTreeStr(t *testing.T) {
	data := map[string][]*Bar{
		"0": {
			{
				ID:   "1",
				Pid:  "0",
				Name: "bar-1",
			},
			{
				ID:   "2",
				Pid:  "0",
				Name: "bar-2",
			},
		},
		"1": {
			{
				ID:   "3",
				Pid:  "1",
				Name: "bar-3",
			},
			{
				ID:   "4",
				Pid:  "1",
				Name: "bar-4",
			},
		},
		"2": {
			{
				ID:   "5",
				Pid:  "2",
				Name: "bar-5",
			},
			{
				ID:   "6",
				Pid:  "2",
				Name: "bar-6",
			},
		},
		"3": {
			{
				ID:   "7",
				Pid:  "3",
				Name: "bar-7",
			},
			{
				ID:   "8",
				Pid:  "3",
				Name: "bar-8",
			},
		},
		"4": {
			{
				ID:   "9",
				Pid:  "4",
				Name: "bar-9",
			},
			{
				ID:   "10",
				Pid:  "4",
				Name: "bar-10",
			},
		},
	}

	tree := New(data, "0")
	b, _ := json.Marshal(tree)
	t.Log(string(b))
}

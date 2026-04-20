package stepkit

import "iter"

type Step struct {
	Head int
	Tail int
}

// New returns steps for a slice.
//
// Example:
//
//	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
//	for step := range stepkit.Split(len(arr), 6) {
//		cur := arr[step.Head:step.Tail]
//		// todo: do something
//	}
func Split(length, step int) iter.Seq[Step] {
	return func(yield func(Step) bool) {
		if length <= 0 || step <= 0 {
			return
		}
		for i := 0; i < length; i += step {
			tail := min(i+step, length)
			if !yield(Step{Head: i, Tail: tail}) {
				return
			}
		}
	}
}

package stepkit

type Step struct {
	Head int
	Tail int
}

// New returns steps for a slice.
//
// Example:
//
//	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
//	for _, step := range stepkit.Split(len(arr), 6) {
//		cur := arr[step.Head:step.Tail]
//		// todo: do something
//	}
func Split(len, step int) []Step {
	mod := len % step
	cap := len / step
	if mod != 0 {
		cap += 1
	}
	end := len - mod
	steps := make([]Step, 0, cap)
	for i := 0; i < end; i += step {
		steps = append(steps, Step{
			Head: i,
			Tail: i + step,
		})
	}
	if mod != 0 {
		steps = append(steps, Step{
			Head: end,
			Tail: len,
		})
	}
	return steps
}

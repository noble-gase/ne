package stepkit

import (
	"log"
	"testing"
)

func TestStep(t *testing.T) {
	log.Println("------------------ arr0 ------------------")
	arr0 := []int{}
	for step := range Split(len(arr0), 6) {
		ids := arr0[step.Head:step.Tail]
		log.Printf("step[%d, %d] = %+v\n", step.Head, step.Tail, ids)
	}

	log.Println("------------------ arr1 ------------------")
	arr1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18}
	for step := range Split(len(arr1), 6) {
		ids := arr1[step.Head:step.Tail]
		log.Printf("step[%d, %d] = %+v\n", step.Head, step.Tail, ids)
	}

	log.Println("------------------ arr2 ------------------")
	arr2 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	for step := range Split(len(arr2), 6) {
		ids := arr2[step.Head:step.Tail]
		log.Printf("step[%d, %d] = %+v\n", step.Head, step.Tail, ids)
	}
}

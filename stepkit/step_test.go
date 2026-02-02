package stepkit

import (
	"log"
	"testing"
)

func TestStep(t *testing.T) {
	log.Println("------------------ arr1 ------------------")
	arr1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18}
	for index, step := range Split(len(arr1), 6) {
		ids := arr1[step.Head:step.Tail]
		log.Printf("step[%d], slice: %d\n", index, ids)
	}

	log.Println("------------------ arr2 ------------------")
	arr2 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	for index, step := range Split(len(arr2), 6) {
		ids := arr2[step.Head:step.Tail]
		log.Printf("step[%d], slice: %d\n", index, ids)
	}
}

package closes

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

func TestClose(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(3)
	go func() {
		defer wg.Done()
		Add("test-1", P1, func() error {
			fmt.Println("test-1 bye-bye")
			return nil
		})
	}()
	go func() {
		defer wg.Done()
		Add("test-3", P3, func() error {
			fmt.Println("test-3 bye-bye")
			return nil
		})
	}()
	go func() {
		defer wg.Done()
		Add("test-2", P2, func() error {
			fmt.Println("test-2 bye-bye")
			return errors.New("test-2 oh no")
		})
	}()
	wg.Wait()

	Close()
}

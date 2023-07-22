package go_common

import (
	"fmt"
	"sync"
	"testing"
)

func TestOne(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	q := NewArrayBlockQueue[int](10)

	go func() {
		for i := 0; i < 100; i++ {
			fmt.Println("put", i)
			q.Put(i)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 100; i++ {
			fmt.Println("take", q.Take())
		}
		wg.Done()
	}()

	wg.Wait()

}

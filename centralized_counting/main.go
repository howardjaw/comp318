package main

import (
	"fmt"
	"sync"

	"github.com/howardjaw/comp318/centralized_counting/de_counting"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1000)

	ctr := de_counting.New()
	for itr := 0; itr < 1000; itr++ {
		go func() {
			defer wg.Done()
			ctr.Increment(1)
		}()
	}
	wg.Wait()

	fmt.Println(ctr.Read())
}

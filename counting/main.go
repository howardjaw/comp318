package main

import (
	"fmt"
	"sync"

	"github.com/howardjaw/comp318/counting/counter"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1000)
	ctr := counter.New()

	for itr := 0; itr < 1000; itr++ {
		go func() {
			defer wg.Done()
			ctr.Increment()
		}()
	}

	wg.Wait()
	count := ctr.Read()
	fmt.Printf("count: %d\n", count)
}

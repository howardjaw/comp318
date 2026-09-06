// /*
// Concurrent program for count that is atomic.
// */

package main

import (
	"fmt"
	// "sync"
)

// import (
// 	"fmt"
// 	"sync"
// )

// var count int
// var mtx sync.Mutex

// func increment(wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	mtx.Lock()
// 	defer mtx.Unlock()
// 	count++
// }

// func main() {
// 	var wg sync.WaitGroup
// 	wg.Add(1000)
// 	for iter := 0; iter < 1000; iter++ {
// 		go increment(&wg)
// 	}
// 	wg.Wait()
// 	mtx.Lock()
// 	defer mtx.Unlock()
// 	fmt.Printf("count: %d\n", count)
// }

// /*
// Encapsulation
// */

// type Counter struct {
// 	mtx   sync.Mutex
// 	count int
// }

// func (c *Counter) increment() {
// 	c.mtx.Lock()
// 	defer c.mtx.Unlock()
// 	c.count++
// }

// func main() {
// 	var c Counter
// 	for iter := 0; iter < 100; iter++ {
// 		c.increment()
// 	}
// 	c.mtx.Lock()
// 	defer c.mtx.Unlock()
// 	fmt.Printf("count: %d\n", c.count)
// }

/*
Go channels
*/

func main() {
	ch := make(chan int)

	go func() {
		ch <- 1
	}()

	result := <-ch
	fmt.Println(result)
}

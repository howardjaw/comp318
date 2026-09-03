package de_counting

type Counter struct {
	increment chan int
	reset     chan bool
	read      chan int
}

func (c *Counter) Increment(val int) {
	c.increment <- val
}

func (c *Counter) Reset() {
	c.reset <- true
}

func (c *Counter) Read() int {
	return <-c.read
}

func New() Counter {
	increment := make(chan int)
	reset := make(chan bool)
	read := make(chan int)

	go func() {
		count := 0
		for {
			select {
			case val := <-increment:
				count += val
			case <-reset:
				count = 0
			case read <- count:
			}
		}
	}()

	return Counter{increment, reset, read}
}

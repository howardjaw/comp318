package counter

type Counter struct {
	count chan int
}

func (c *Counter) Increment() {
	curr := <-c.count
	curr++
	go func() {
		c.count <- curr
	}()
}

func (c *Counter) Read() int {
	val := <-c.count
	go func() {
		c.count <- val
	}()
	return val
}

func New() Counter {
	count := make(chan int)
	go func() {
		count <- 0
	}()
	return Counter{count: count}
}

package sync1

import "sync"

type Counter struct {
	mu sync.Mutex //A Mutex is a mutual 
	// exclusion lock. The zero value for a Mutex is an unlocked mutex.
	value int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) Value() int {
	return c.value
}


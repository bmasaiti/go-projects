package sync1

import (
	"sync"
	"testing"
)

func TestCounter(t *testing.T) {

	t.Run("Increamenting the counter 3 times leaves it at 3", func(t *testing.T) {
		counter := NewCounter()
		
		counter.Inc()
		counter.Inc()
		counter.Inc()
		assertCounter(t, counter, 3)
		// if counter.Value() != 3 {
		// 	t.Errorf("got %d, want %d ", counter.Value(), 3)
		// }
	})	

	t. Run("it runs safely concurrently", func(t *testing.T){
		wantedCount := 1000
		counter := NewCounter()

		var wg sync.WaitGroup
		wg.Add(wantedCount) //how many routines to wait for

		for i := 0; i<wantedCount; i++{
			go func(){ //kickoff a good routine for each count
				counter.Inc()
				wg.Done()
			}()

		}
		wg.Wait() //wait for all go routines to return
		assertCounter(t, counter, wantedCount)  //compare results
	})
}

//func assertCounter(t testing.TB, got Counter , want int)

func NewCounter() *Counter {
	return &Counter{}
}
func assertCounter(t testing.TB, got *Counter, want int){ //A Mutex must not be copied after first use.
	t.Helper()	
	if got.Value()!= want {
		t.Errorf("got %d, want %d", got.Value(), want)
	}
}
// the interesting thing here is all the go routines are trying to 
//change the value at the same memory locatiion, certainly some will
//overwrite each other.,
package selecta

import (
	"fmt"
	"net/http"
	"time"
)

// httptest : A convenient way of creating test servers so you can have reliable and controllable tests.
// select :   you wait on multiple channels.
// time.After() : Sometimes you'll want to include time.After in one of your cases to prevent your system blocking forever.
//go routine (go func())
// http.Get()

func Racer(a, b string) (winner string) {
	//We don't really care about the exact response times of the requests, we just want to know which one comes back first.
	aDuration := measureResponseTime(a)
	bDuration := measureResponseTime(b)
	// startA :=time.Now()
	// http.Get(a)
	// aDuration:=time.Since(startA)

	// startB := time.Now()
	// http.Get(b)
	// bDuration := time.Since(startB)

	if aDuration < bDuration {
		return a
	}
	return b

}

func measureResponseTime(url string) time.Duration {
	start := time.Now()
	http.Get(url)
	return time.Since(start)
}

var tenSecondTimeout = 10 * time.Second

func Racer3(a, b string) (winner string, error error) {
	return ConfigurableRacer(a, b, tenSecondTimeout)
}

func ConfigurableRacer(a, b string, timeout time.Duration) (winner string, error error) {
	select { //arguments to select are evaluated at the same time , which means ping)a) and ping(b) are called at the same time
	//which menas their response times are impacted only by the actual urls called.
	// Select waits for these responses to check which one was first.
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timed out waiting for %s and %s ", a, b)
	}
}

//select allows you to wait on multiple channels. The first one to send a value "wins" and the code underneath the case is executed.

func ping(url string) chan struct{} {
	ch := make(chan struct{})

	go func() { //remember go func() creates a goroutine the unit of concurrency .
		http.Get(url)
		close(ch)
	}()
	return ch
}

/*
	Why struct{} and not another type like a bool? Well, a chan struct{} is the smallest data type available from
	a memory perspective so we get no allocation versus a bool. Since we are closing and not sending anything on the
	chan, why allocate anything?

	Always make channels

	Notice how we have to use make when creating a channel; rather than say var ch chan struct{}. When you use var the
	 variable will be initialised with the "zero" value of the type. So for string it is "", int it is 0, etc.

	For channels the zero value is nil and if you try and send to it with <- it will block forever because
	you cannot send to nil channels
*/

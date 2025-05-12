package selecta


// httptest
// select
//go routine (go func())
// http.Get()

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// func TestRacer(t *testing.T){

// 	slowUel := "http://www.facebook.com"
// 	fastUrl:= "http://www.quii.dev"

// 	want := fastUrl
// 	got:= Racer(slowUel,fastUrl)

// 	if got != want {
// 		t.Errorf("got %q , want %q ", got ,want)
// 	}
// }

func TestRacer(t *testing.T) {
	// slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter , r *http.Request){
	// 	time.Sleep(20* time.Millisecond)
	// 	w.WriteHeader(http.StatusOK)
	// }))

	// fastServer := httptest.NewServer(http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
	// 	w.WriteHeader(http.StatusOK)
	// }))

	slowServer := makeDelayedServer(20 * time.Millisecond)
	fastServer := makeDelayedServer(0 * time.Millisecond)

	slowUrl := slowServer.URL
	fastUrl := fastServer.URL

	want := fastUrl
	got := Racer(slowUrl, fastUrl)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	//By prefixing a function call with defer it will now call that function at the end of the containing function.
	defer slowServer.Close()
	defer fastServer.Close()
}

func makeDelayedServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
	}))
}

func TestRacer2(t *testing.T) {
	t.Run("compares speeds of servers ,returning the url of the fastest one", func(t *testing.T) {
		slowServer := makeDelayedServer(20 * time.Millisecond)
		fastServer := makeDelayedServer(0 * time.Millisecond)

		defer slowServer.Close()
		defer fastServer.Close()

		slowUrl := slowServer.URL
		fastUrl := fastServer.URL

		want := fastUrl
		got, err := Racer3(slowUrl, fastUrl)

		if err != nil {
			t.Fatalf("did not expect an error but got one %v", err)
		}

		if got != want {
			t.Errorf("got %q , want %q", got, want)
		}
	})

	t.Run("returns an error if a server doesn't respond within 10s", func(t *testing.T) {
		server := makeDelayedServer(25 * time.Millisecond)
		

		defer server.Close()
		
		_, err := ConfigurableRacer(server.URL, server.URL, 20*time.Millisecond)


		if err == nil {
			t.Error("expected an error but didn't get one")
		}
	})
}

package main

import "testing"

//1. basic test
// func TestHello(t *testing.T) {
// 	got := Hello("Chris")
// 	want := "Hello , Chris"

// 	if got != want {
// 		t.Errorf("got %q want %q", got, want)
// 	}
// }

//2. subtests
// func Test2Hello(t *testing.T) {
// 	t.Run("Saying hello to people ", func(t *testing.T){
// 		got := Hello("Chris")
// 		want := "Hello , Chris"

// 		if got != want {
// 			t.Errorf("got %q want %q", got, want)
// 		}
// 	})

// 	t.Run("say 'Hello, World' when an empty string is supplied", func(t *testing.T) {
// 		got := Hello("")
// 		want := "Hello, World"

// 		if got != want {
// 			t.Errorf("got %q want %q", got, want)
// 		}

// })
// }

// 3. Using assertions by manually creating an assertion function which takes the context and actual and expected resposnes
func Test3Hello(t *testing.T) {
	// t.Run("saying hello to people", func(t *testing.T) {
	// 	got := Hello("Chris")
	// 	want := "Hello , Chris"
	// 	assertCorrectMessage(t, got, want)
	// })

	// t.Run("empty string defaults to 'world'", func(t *testing.T) {
	// 	got := Hello("")
	// 	want := "Hello , World"
	// 	assertCorrectMessage(t, got, want)
	// })

	t.Run("in Spanish", func(t *testing.T) {
		got := Hello("Elodie", "Spanish")
		want := "Hola , Elodie"
		assertCorrectMessage(t, got, want)
	})

}

func assertCorrectMessage(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

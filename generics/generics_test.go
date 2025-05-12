package generics

import "testing"

func TestAssertFunctions(t *testing.T) {
	t.Run("Asserting on integers", func(t *testing.T) {
		AssertEqual(t, 1, 1)
		AssertNotEqual(t, 1, 2)
	})
	t.Run("asserting on strings", func(t *testing.T) {
		AssertEqual(t, "hello", "hello")
		AssertNotEqual(t, "hello", "Grace")
	})
}
//basic assert
// func AssertEqual(t *testing.T, got, want int) {
// 	t.Helper()
// 	if got != want {
// 		t.Errorf("got %d , want %d", got, want)
// 	}
// }

// func AssertNotEqual(t *testing.T, got, want int) {
// 	t.Helper()
// 	if got == want {
// 		t.Errorf("didnt want %d", got)
// 	}
// }

//assert with generics
func AssertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %T , want %T", got, want)
	}
}

func AssertNotEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got == want {
		t.Errorf("didnt want %T", got)
	}
}


package shapes

import "testing"

func TestPerimeter(t *testing.T) {
	rectangle := Rectangle{10.0, 10.0}
	got := Perimeter(rectangle)
	want := 40.0
	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}

// func TestArea(t *testing.T) {
// 	rectangle := Rectangle{3.0, 5.0}
// 	got := Area(rectangle)
// 	want := 15.0

// 	if got != want {
// 		t.Errorf("Got %.2f want %.2f", got, want)
// 	}
// }

func TestArea(t *testing.T) {

	t.Run("rectangles", func(t *testing.T) {

		rectangle := Rectangle{3.0, 5.0}
		got := rectangle.Area()
		want := 15.0

		if got != want {
			t.Errorf("Got %.2f want %.2f", got, want)
		}
	})

	t.Run("circles", func(t *testing.T) {
		circle := Circle{10}
		got := circle.Area()
		want := 314.1592653589793

		if got != want {
			t.Errorf("Got %g want %g ", got, want)
		}

	})
}

func TestAre(t *testing.T) {
	checkArea := func(t testing.TB, shape Shape, want float64) {
		t.Helper()
		got := shape.Area()
		if got != want {
			t.Errorf("got %g want %g", got, want)
		}
	}

	t.Run("rectangles", func(t *testing.T) {
		rectangle := Rectangle{4, 6}
		checkArea(t, rectangle, 24.0)

	})
	t.Run("circles", func(t *testing.T) {
		circle := Circle{10}
		checkArea(t, circle, 314.1592653589793)

	})
}

// In Go interface resolution is implicit. If the type you pass in matches what the interface is asking for, it will compile.
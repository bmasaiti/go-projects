package shapes

import "math"

type Rectangle struct {
	Width  float64
	Height float64
}

type Circle struct {
	Radius float64
}
type Shape interface {
	Area() float64
}

func Perimeter(rect Rectangle) float64 {
	return 2 * (rect.Width + rect.Height)

}

// func Area(rect Rectangle) float64 {

//		return rect.Width * rect.Height
//	}
//
// func (receiverName ReceiverType) MethodName(args)
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Wait, what?
// This is quite different to interfaces in most other programming languages. Normally you have to write code to say My type Foo implements interface Bar.

// But in our case

// Rectangle has a method called Area that returns a float64 so it satisfies the Shape interface

// Circle has a method called Area that returns a float64 so it satisfies the Shape interface

// string does not have such a method, so it doesn't satisfy the interface

// etc.
package main

//array
// func Sum(numbers [5]int) int {
// 	sum := 0
// 	for _, number:= range numbers {
// 		sum += number
// 	}
// 	return sum
// }

// slice
func Sum(numbers []int) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}

func SumAll(numbersToSum ...[]int) []int {
	// lengthOfNumbers := len(numbersToSum)
	// sums := make([]int, lengthOfNumbers) //using make to create slice with init capacity
	var sums []int
	for _, numbers := range numbersToSum {
		sums = append(sums, Sum(numbers))
	}

	return sums
}

func SumAllTails(numbersToAdd ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToAdd {
		tail := numbers[1:]
		sums = append(sums, Sum(tail))
	}
	return sums
}

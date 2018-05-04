package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println("hello world")

	myValues := [][]int{{1, 2, 3, 7}, {4, 5, 6, 9}, {7, 8, 9}, {}}
	mySlices := make([][]int, len(myValues))
	for i := 0; i < len(myValues); i++ {
		mySlices[i] = myValues[i][0:]
	}
	fmt.Println(mySlices)
	finalOut := mergeSort(mySlices)
	fmt.Println(finalOut)
	return
}

func mergeSort(input [][]int) []int {
	finalLength := 0
	for i := 0; i < len(input); i++ {
		finalLength += len(input[i])
	}
	output := make([]int, finalLength)

	lowestVal := int(math.Inf(1)) - 1
	bookmark := -1
	for i := 0; i < finalLength; i++ {
		for j := 0; j < len(input); j++ {
			if len(input[j]) == 0 {
				continue
			}
			if lowestVal > input[j][0] {
				lowestVal = input[j][0]
				bookmark = j
			}
		}
		output[i] = lowestVal
		fmt.Printf("%d is the lowVal, %d is the bookmark", lowestVal, bookmark)
		fmt.Println(input[bookmark])
		input[bookmark] = input[bookmark][1:]
		lowestVal = int(math.Inf(1))
	}

	return output
}

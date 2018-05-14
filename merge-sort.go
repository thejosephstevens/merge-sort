package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"unsafe"
)

const (
	multiWorker = iota
	singleLayer
	singleWorker
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Insufficient arguments, please enter `merge-sort.exe <number of arrays> <length of arrays> <sort-mode> (optional)`\n<sort-mode> can be single-layer, single-worker, or multi-worker")
		return
	}

	minVal := 0
	maxVal := 1
	testVal := 0
	intSize := unsafe.Sizeof(testVal)
	if intSize == 8 {
		minVal = math.MinInt64
		maxVal = math.MaxInt64
	} else if intSize == 4 {
		minVal = math.MinInt32
		maxVal = math.MaxInt32
	} else {
		fmt.Printf("panic, unexpected OS instruction size %d\n", intSize*8)
		return
	}

	fmt.Println(minVal)
	fmt.Println(maxVal)

	numSlices, err := strconv.Atoi(os.Args[1])
	if err != nil || numSlices < 0 {
		fmt.Printf("Invalid number of arrays, you inputted %s\n", os.Args[2])
	}

	sliceLength, err := strconv.Atoi(os.Args[2])
	if err != nil || sliceLength < 0 {
		fmt.Printf("Invalid length of arrays, you inputted %s\n", os.Args[3])
	}

	sortMode := 0
	if len(os.Args) > 3 {
		switch sortModeString := os.Args[3]; sortModeString {
		case "multi-worker":
			sortMode = multiWorker
		case "single-layer":
			sortMode = singleLayer
		case "single-worker":
			sortMode = singleWorker
		default:
			fmt.Printf("Unknown sort option selected %s, defaulting to multi-worker\n", sortModeString)
			sortMode = multiWorker
		}
	} else {
		fmt.Println("No sort method chosen, multi-worker will be used as the default.")
		sortMode = multiWorker
	}

	avgInterval := 2 * (maxVal / sliceLength)

	myValues := make([][]int, numSlices)
	for i := 0; i < numSlices; i++ {
		myValues[i] = make([]int, sliceLength)
		curVal := minVal + rand.Intn(avgInterval)/2 // just to help be sure we're not always going over
		for j := 0; j < sliceLength; j++ {
			if curVal > maxVal {
				curVal = maxVal
			}
			myValues[i][j] = curVal
			curVal += rand.Intn(avgInterval)
		}
	}
	fmt.Println(myValues)
	finalOut := make([]int, 0)
	switch sortMode {
	case multiWorker:

	case singleLayer:
		finalOut = mergeSort(myValues, maxVal)
	case singleWorker:

	}
	fmt.Println(finalOut)
	return
}

func mergeSort(input [][]int, maxVal int) []int {
	finalLength := 0
	for i := 0; i < len(input); i++ {
		finalLength += len(input[i])
	}
	output := make([]int, finalLength)

	bookmark := -1
	for i := 0; i < finalLength; i++ {
		for j := 0; j < len(input); j++ {
			if len(input[j]) == 0 {
				continue
			}
			if maxVal > input[j][0] {
				maxVal = input[j][0]
				bookmark = j
			}
		}
		output[i] = maxVal
		fmt.Printf("%d is the maxVal, %d is the bookmark", maxVal, bookmark)
		fmt.Println(input[bookmark])
		input[bookmark] = input[bookmark][1:]
		maxVal = int(math.Inf(1))
	}

	return output
}

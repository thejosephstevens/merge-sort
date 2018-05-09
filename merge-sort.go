package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	MaxInt8   = 1<<7 - 1
	MinInt8   = -1 << 7
	MaxInt16  = 1<<15 - 1
	MinInt16  = -1 << 15
	MaxInt32  = 1<<31 - 1
	MinInt32  = -1 << 31
	MaxInt64  = 1<<63 - 1
	MinInt64  = -1 << 63
	MaxUint8  = 1<<8 - 1
	MaxUint16 = 1<<16 - 1
	MaxUint32 = 1<<32 - 1
	MaxUint64 = 1<<64 - 1
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(errors.New("Input the number of arrays you want to merge"))
		return
	}
	numSlices, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(errors.New("Please input a valid integer in the command line"))
		return
	}
	startPrep := time.Now()
	mySlices := make([][]int, numSlices) //{{1, 2, 3, 7}, {4, 5, 6, 9}, {7, 8, 9}, {}}
	// this might be horribly slow
	totalVals := 0
	for n := 0; n < numSlices; n++ {
		newSlice := make([]int, 0)
		for i := rand.Intn(MaxInt16); i < MaxInt16; i += (rand.Intn(MaxInt16) / 4096) {
			newSlice = append(newSlice, i)
			totalVals += 1
		}
		mySlices[n] = newSlice
	}
	endPrep := time.Now()
	prepTime := int(endPrep.Sub(startPrep) / time.Millisecond)
	fmt.Printf("Generated input arrays in %v milliseconds\n", prepTime)
	//fmt.Println(mySlices)
	routines := 0
	slices := make(chan []int, 16)
	start := time.Now()
	for true {
		if routines == 0 && len(mySlices) == 1 {
			break
		}
		for len(mySlices) > 1 {
			go mergeSort(mySlices[0:2], slices)
			routines += 1
			mySlices = mySlices[2:]
		}
		nextSlice := <-slices
		routines -= 1
		mySlices = append(mySlices, nextSlice)
	}
	end := time.Now()
	timeElapsed := int(end.Sub(start) / time.Millisecond)
	avgLength := totalVals / numSlices
	fmt.Printf("Sorted %d arrays with an average length of %d in %v milliseconds.\n", numSlices, avgLength, timeElapsed)
	return
}

func mergeSort(input [][]int, out chan []int) {
	finalLength := 0
	for i := 0; i < len(input); i++ {
		finalLength += len(input[i])
	}
	output := make([]int, finalLength)

	lowestVal := MaxInt32
	bookmark := 0
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
		input[bookmark] = input[bookmark][1:]
		lowestVal = MaxInt32
	}

	//fmt.Printf("the output is %v\n", output)
	out <- output
}

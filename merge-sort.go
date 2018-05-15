package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	multiWorker = iota
	singleLayer
	singleWorker
)
const maxUint = ^uint(0)
const maxInt = int(maxUint >> 1)
const minInt = -maxInt - 1

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "help") {
		fmt.Printf("This is a sample go program that generates a random set of integer slices and then merge-sorts them.\nEnter `merge-sort.exe <number of arrays> <length of arrays> <sort-mode> (optional)`\n<sort-mode> can be single-layer, single-worker, multi-worker, or you can leave it blank (it will default to multi-worker)\n")
		return
	}
	var stratLookup = map[int]string{
		multiWorker:  "multi-worker",
		singleLayer:  "single-layer",
		singleWorker: "single-worker",
	}
	if len(os.Args) < 3 {
		fmt.Println("Insufficient arguments, please enter `merge-sort.exe <number of arrays> <length of arrays> <sort-mode> (optional)`\n<sort-mode> can be single-layer, single-worker, or multi-worker")
		return
	}

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

	myValues := make([][]int, numSlices)
	newSlices := make(chan []int, 64)
	writeData := make(chan []int, 64)
	unsortedWriteComplete := make(chan bool, 0)
	go generateIntSliceHelper(numSlices, sliceLength, newSlices)
	go writeIntSlices("unsorted_data.txt", numSlices, sliceLength, writeData, unsortedWriteComplete)
	for i := 0; i < numSlices; i++ {
		myValues[i] = <-newSlices
		writeData <- myValues[i]
	}

	routines := 0
	slices := make(chan []int, 16) // a small buffer
	start := time.Now()
	switch sortMode {
	case multiWorker:
		for true {
			if routines == 0 && len(myValues) == 1 {
				break
			}
			for len(myValues) > 1 {
				go mergeSort(myValues[0:2], slices)
				routines++
				myValues = myValues[2:]
			}
			nextSlice := <-slices
			routines--
			myValues = append(myValues, nextSlice)
		}
	case singleLayer:
		go mergeSort(myValues, slices)
		myValues = append(myValues, make([]int, 1))
		myValues = myValues[len(myValues)-1 : len(myValues)]
		myValues[0] = <-slices
	case singleWorker:
		for true {
			if len(myValues) == 1 {
				break
			}
			go mergeSort(myValues[0:2], slices)
			myValues = myValues[2:]
			nextSlice := <-slices
			myValues = append(myValues, nextSlice)
		}
	}
	end := time.Now()
	timeElapsed := int(end.Sub(start) / time.Millisecond)

	filename := "sorted_data.txt"
	os.Create(filename)
	outfile, _ := os.OpenFile(filename, os.O_RDWR, 0666)
	defer outfile.Close()
	outfile.WriteString(fmt.Sprintln(myValues[0]))
	<-unsortedWriteComplete
	fmt.Printf("Sorted %d arrays of length %d using %s strategy in %v milliseconds.\n", numSlices, sliceLength, stratLookup[sortMode], timeElapsed)
	return
}

func mergeSort(input [][]int, out chan []int) {
	finalLength := 0
	for i := 0; i < len(input); i++ {
		finalLength += len(input[i])
	}
	output := make([]int, finalLength)

	highVal := maxInt
	bookmark := -1
	for i := 0; i < finalLength; i++ {
		for j := 0; j < len(input); j++ {
			if len(input[j]) == 0 {
				continue
			}
			if highVal > input[j][0] {
				highVal = input[j][0]
				bookmark = j
			}
		}
		output[i] = highVal
		input[bookmark] = input[bookmark][1:]
		highVal = maxInt
	}

	out <- output
}

func generateIntSliceHelper(numSlices int, length int, out chan []int) {
	for i := 0; i < numSlices; i++ {
		go generateIntSlice(length, out)
	}
}

func generateIntSlice(length int, out chan []int) {
	avgInterval := 2 * (maxInt / length)
	newSlice := make([]int, length)
	curVal := minInt + rand.Intn(avgInterval)/2
	for j := 0; j < length; j++ {
		if curVal > maxInt {
			curVal = maxInt
		}
		newSlice[j] = curVal
		curVal += rand.Intn(avgInterval)
	}
	out <- newSlice
}

func writeIntSlices(filename string, numSlices int, sliceLength int, in chan []int, done chan bool) {
	os.Create(filename)
	outfile, _ := os.OpenFile(filename, os.O_RDWR, 0666)
	defer outfile.Close()
	for i := 0; i < numSlices; i++ {
		nextSlice := <-in
		outfile.WriteString(fmt.Sprintln(nextSlice))
	}
	done <- true
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

func split(start int, end int, shouldTakeUpper bool) (newStart int, newEnd int) {
	if shouldTakeUpper {
		newStart = start + (end-start)/2 + (end-start)%2
		newEnd = end
	} else {
		newStart = start
		newEnd = start + (end-start)/2
	}
	return
}

func getRow(code string) int {
	start := 0
	end := 127
	for i := 0; i < 7; i++ {
		if code[i] == 'F' {
			start, end = split(start, end, false)
		} else if code[i] == 'B' {
			start, end = split(start, end, true)
		} else {
			fmt.Println("Something went wrong, unexpected row code", code, i)
		}
	}
	if start != end {
		fmt.Println("Something went wrong, start doesn't equal end.", start, end)
	}
	return start
}

func getCol(code string) int {
	start := 0
	end := 7
	for i := 7; i < 10; i++ {
		if code[i] == 'L' {
			start, end = split(start, end, false)
		} else if code[i] == 'R' {
			start, end = split(start, end, true)
		} else {
			fmt.Println("Something went wrong, unexpected col code", code, i)
		}
	}
	if start != end {
		fmt.Println("Something went wrong, start doesn't equal end.", start, end)
	}
	return start
}

func getSeatID(code string) int {
	row := getRow(code)
	col := getCol(code)
	return row*8 + col
}

func main() {
	file, err := os.Open("./five-input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	highestID := 0
	allIDs := make([]int, 0)
	for scanner.Scan() {
		// Example: BBFFFBFRLR
		line := scanner.Text()
		if err != nil {
			log.Fatal(err)
		}
		seatID := getSeatID(line)
		allIDs = append(allIDs, seatID)
		if seatID > highestID {
			highestID = seatID
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Highest seat id:", highestID)
	sort.Ints(allIDs)
	for i := 1; i < len(allIDs)-1; i++ {
		// Find the missing seatID
		if allIDs[i]+1 != allIDs[i+1] {
			fmt.Println("Found a gap:", allIDs[i-1], allIDs[i], allIDs[i+1])
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func loadNumbers() []int {
	file, err := os.Open("./ten-input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	list := make([]int, 1) // first value of 0
	for scanner.Scan() {
		// Example: 23
		line := scanner.Text()
		if err != nil {
			log.Fatal(err)
		}
		number, _ := strconv.Atoi(line)
		list = append(list, number)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return list
}

func roundOne() {
	numbers := loadNumbers()
	sort.Ints(numbers)
	// The final device joltage
	numbers = append(numbers, numbers[len(numbers)-1]+3)
	differences := map[int]int{
		1: 0, 2: 0, 3: 0,
	}
	for i := 0; i < len(numbers); i++ {
		previous := 0
		if i > 0 {
			previous = numbers[i-1]
		}
		differences[numbers[i]-previous]++
	}
	fmt.Println("the differences", differences)
	fmt.Println("1 jolts times 3 jolts is", differences[1]*differences[3])
}

// The number of chains to the end from start index in numbers,
// caching previously computed values in knownWays
func numChains(numbers []int, start int, knownWays map[int]int) int {
	if start == len(numbers)-1 {
		knownWays[start] = 1
		return 1
	}
	numWays := 0
	// [..., 4, 5, 6, 9]
	for i := start + 1; i < len(numbers); i++ {
		val := numbers[i]
		if val-numbers[start] > 3 {
			break
		}
		if cached, ok := knownWays[i]; ok {
			numWays += cached
		} else {
			numWays += numChains(numbers, i, knownWays)
		}
	}
	knownWays[start] = numWays
	return numWays
}

func main() {
	// Round 2
	numbers := loadNumbers()
	sort.Ints(numbers)
	// the final device joltage
	numbers = append(numbers, numbers[len(numbers)-1]+3)
	fmt.Println(numbers)
	knownWays := make(map[int]int)
	numChains := numChains(numbers, 0, knownWays)
	fmt.Println(numChains, "ways")
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func loadNumbers() []int {
	file, err := os.Open("./nine-input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	list := make([]int, 0)
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

func isValid(numbers []int, targetIndex int, window int) bool {
	for xIndex := targetIndex - window; xIndex < targetIndex; xIndex++ {
		for yIndex := xIndex + 1; yIndex < targetIndex; yIndex++ {
			if numbers[xIndex]+numbers[yIndex] == numbers[targetIndex] {
				return true
			}
		}
	}
	return false
}

func decode(numbers []int, window int) int {
	for i := window; i < len(numbers)-window; i++ {
		if !isValid(numbers, i, window) {
			return numbers[i]
		}
	}
	return -1
}

func findSum(numbers []int, target int) int {
	for i := 0; i < len(numbers); i++ {
		// Look for a window starting at i
		sum := numbers[i]
		for j := i + 1; j < len(numbers); j++ {
			sum += numbers[j]
			if sum == target {
				// Find smallest and largest in this range
				smallest := numbers[i]
				largest := numbers[i]
				for k := i; k <= j; k++ {
					if numbers[k] < smallest {
						smallest = numbers[k]
					}
					if numbers[k] > largest {
						largest = numbers[k]
					}
				}
				return smallest + largest
			}
			if sum > target {
				break
			}
		}
	}
	return -1
}

func main() {
	// Round 1
	numbers := loadNumbers()
	invalid := decode(numbers, 25)
	fmt.Println("the first number that is bad is", invalid)

	// Round 2
	encryptionWeakness := findSum(numbers, invalid)
	fmt.Println("the encryption weakness is", encryptionWeakness)
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func roundOne() {
	file, err := os.Open("./two-input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	numValid := 0

	for scanner.Scan() {
		// num, err := strconv.Atoi(scanner.Text())
		// Example: 3-6 v: vvdvvbv
		line := scanner.Text()
		if err != nil {
			log.Fatal(err)
		}
		s := strings.Fields(line)
		requiredRange := strings.Split(s[0], "-")
		minCount, _ := strconv.Atoi(requiredRange[0])
		maxCount, _ := strconv.Atoi(requiredRange[1])
		requiredCharacter := strings.TrimRight(s[1], ":")
		password := s[2]
		occurrences := 0
		for _, char := range password {
			if string(char) == requiredCharacter {
				occurrences++
			}
		}
		if occurrences >= minCount && occurrences <= maxCount {
			numValid++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(numValid)
}

func roundTwo() {
	file, err := os.Open("./two-input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	numValid := 0

	for scanner.Scan() {
		// num, err := strconv.Atoi(scanner.Text())
		// Example: 3-6 v: vvdvvbv
		line := scanner.Text()
		if err != nil {
			log.Fatal(err)
		}
		s := strings.Fields(line)
		// These are 1-indexed, not zero indexed
		requiredIndexes := strings.Split(s[0], "-")
		requiredIndex1, _ := strconv.Atoi(requiredIndexes[0])
		requiredIndex2, _ := strconv.Atoi(requiredIndexes[1])
		requiredCharacter := strings.TrimRight(s[1], ":")
		password := s[2]
		numMatches := 0
		if string(password[requiredIndex1-1]) == requiredCharacter {
			numMatches++
		}
		if string(password[requiredIndex2-1]) == requiredCharacter {
			numMatches++
		}
		if numMatches == 1 {
			numValid++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(numValid)
}

func main() {
	roundTwo()
}

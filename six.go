package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func roundOne() {
	file, err := os.Open("./six-input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	total := 0
	questions := make(map[rune]bool)
	for scanner.Scan() {
		// Example: mxfhdeyikljnz
		line := scanner.Text()
		if err != nil {
			log.Fatal(err)
		}
		if line == "" {
			// End of this group
			total += len(questions)
			questions = make(map[rune]bool)
		} else {
			for _, question := range line {
				questions[question] = true
			}
		}
	}
	// Process the final group too
	total += len(questions)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Total sum:", total)
}

func main() {
	// Round 2
	file, err := os.Open("./six-input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	total := 0
	questions := map[rune]bool(nil)
	for scanner.Scan() {
		// Example: mxfhdeyikljnz
		line := scanner.Text()
		if err != nil {
			log.Fatal(err)
		}
		if line == "" {
			// End of this group
			total += len(questions)
			questions = nil
		} else if questions == nil {
			// The first person in the group
			questions = make(map[rune]bool)
			for _, question := range line {
				questions[question] = true
			}
		} else {
			// Remove any questions that this person hasn't answered yes to
			for key := range questions {
				if !strings.ContainsRune(line, key) {
					delete(questions, key)
				}
			}
		}
	}
	// Process the final group too
	total += len(questions)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Total sum:", total)
}

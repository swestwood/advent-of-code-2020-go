package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func load() {
	file, err := os.Open("./eighteen-input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		// Example:
		line := scanner.Text()
		if err != nil {
			log.Fatal(err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	load()
	fmt.Println("Template")
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open("./one-input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	numbers := make([]int, 0)

	for scanner.Scan() {
		num, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		numbers = append(numbers, num)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for _, a := range numbers {
		for _, b := range numbers {
			for _, c := range numbers {

				if a+b+c == 2020 {
					fmt.Println(a * b * c)
					return
				}
			}
		}
	}
}
